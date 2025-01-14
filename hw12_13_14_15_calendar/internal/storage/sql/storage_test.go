package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"                                       //nolint:depguard
	"github.com/homerchik/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	_ "github.com/jackc/pgx/v5/stdlib"                             //nolint:depguard
	"github.com/pressly/goose/v3"                                  //nolint:depguard
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"                  //nolint:depguard
	"github.com/testcontainers/testcontainers-go/modules/postgres" //nolint:depguard
	"github.com/testcontainers/testcontainers-go/wait"             //nolint:depguard
)

type DBSuite struct {
	suite.Suite
	PgContainer postgres.PostgresContainer
	PgConfig    DBConfig
	SQLStorage  *Storage
}

type DBConfig struct {
	name     string
	password string
	DBName   string
}

func (s *DBSuite) PerformMigration(path string) {
	ctx := context.Background()
	dsn, err := s.PgContainer.ConnectionString(ctx)
	s.NoError(err, "can't acquire dsn")
	db, err := goose.OpenDBWithDriver("pgx", dsn)
	s.NoError(err, "can't connect to db")
	defer func() { s.NoError(db.Close()) }()
	s.NoError(goose.RunContext(ctx, "up", db, path), "can't apply migrations")
}

func (s *DBSuite) SetupSuite() {
	ctx := context.Background()
	s.PgConfig = DBConfig{"calendar", "postgres", "postgres"}

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(s.PgConfig.DBName),
		postgres.WithUsername(s.PgConfig.name),
		postgres.WithPassword(s.PgConfig.password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	s.NoError(err, "failed to start container")
	s.PgContainer = *postgresContainer
	s.PerformMigration("../../../migrations")
	s.SQLStorage = New()
	dsn, err := s.PgContainer.ConnectionString(ctx)
	s.NoError(err, "can't acquire dsn")
	s.NoError(s.SQLStorage.Connect(ctx, dsn, "pgx"))
}

func (s *DBSuite) TeardownSuite() {
	err := testcontainers.TerminateContainer(s.PgContainer)
	s.NoError(err, "failed to terminate container")
}

func (s *DBSuite) TestEventAddedToEmptyStorage() {
	startTime := time.Now()
	event := storage.NewEvent(
		uuid.New(), "Event 1", "Description 1", startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
	)
	err := s.SQLStorage.AddEvent(event)
	s.NoError(err)
	eventFromDB, err := s.SQLStorage.GetEventByID(event.ID)
	s.NoError(err)
	s.Equal(event, eventFromDB, "Fetched and pushed events are different")
}

func (s *DBSuite) TestDeleteExistingEventOneForADate() {
	userID := uuid.New()
	startTime := time.Now()
	events := storage.Schedule{
		storage.NewEvent(
			userID, "Event 1", "Description 1", startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 2", "Description 2", startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 3", "Description 3", startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		s.NoError(s.SQLStorage.AddEvent(event))
	}
	s.NoError(s.SQLStorage.DeleteEvent(events[1].ID))
	_, err := s.SQLStorage.GetEventByID(events[1].ID)
	s.Error(err, storage.ErrNoEventFound)
}

func (s *DBSuite) TestEventUpdateSimpleFields() {
	startTime := time.Now()
	event := storage.NewEvent(
		uuid.New(), "Event 1", "Description 1", startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
	)
	s.NoError(s.SQLStorage.AddEvent(event))
	updatedEvent := storage.NewEvent(
		event.UserID, "Better than event 1", "Simple Des",
		startTime.Add(time.Hour), startTime.Add(2*time.Hour), 10*time.Minute,
	)
	updatedEvent.ID = event.ID
	s.NoError(s.SQLStorage.UpdateEvent(updatedEvent))
	eventFromDB, err := s.SQLStorage.GetEventByID(event.ID)
	s.NoError(err)
	s.Equal(updatedEvent, eventFromDB, "Fetched and pushed events are different")
}

func (s *DBSuite) TestEventUpdateDateFields() {
	userID := uuid.New()
	startTime := time.Now()
	events := storage.Schedule{
		storage.NewEvent(
			userID, "Event 1", "Description 1", startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 2", "Description 2", startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 3", "Description 3", startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		s.NoError(s.SQLStorage.AddEvent(event))
	}
	s.T().Run("Check start date is changed, it's possible", func(_ *testing.T) {
		updatedEvent := storage.NewEvent(
			userID, "Event 1", "Description 1", startTime.Add(30*time.Minute), events[0].EndDate, time.Minute,
		)
		updatedEvent.ID = events[0].ID
		s.NoError(s.SQLStorage.UpdateEvent(updatedEvent))
		eventFromDB, err := s.SQLStorage.GetEventByID(updatedEvent.ID)
		s.NoError(err)
		s.Equal(updatedEvent, eventFromDB, "Fetched and pushed events are different")
	})

	s.T().Run("Check end date is changed, and it's possible", func(_ *testing.T) {
		updatedEvent := storage.NewEvent(
			userID, "Event 1", "Description 1", events[0].StartDate, events[0].EndDate.Add(-30*time.Minute), time.Minute,
		)
		updatedEvent.ID = events[0].ID
		s.NoError(s.SQLStorage.UpdateEvent(updatedEvent))
		eventFromDB, err := s.SQLStorage.GetEventByID(updatedEvent.ID)
		s.NoError(err)
		s.Equal(updatedEvent, eventFromDB, "Fetched and pushed events are different")
	})
}

func TestSQLStorage(t *testing.T) {
	suite.Run(t, new(DBSuite))
}
