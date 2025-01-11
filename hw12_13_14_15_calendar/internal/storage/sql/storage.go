package sqlstorage

import (
	"context"
	"time"
	"errors"

	"database/sql"

	"github.com/google/uuid"
	"github.com/homerchik/hw12_13_14_15_calendar/internal/logic"
	"github.com/homerchik/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
)


type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, connStr string, driver string) error {
	var err error
	if s.db, err = sql.Open(driver, connStr); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddEvent(event storage.Event) error {
	if err := logic.CheckEvent(event, s.GetUserEvents(event.UserId)); err != nil {
		return errors.Join(err, storage.ErrEventCantBeAdded)
	}
	notifyBefore := event.NotifyBefore.Seconds()
	if _, err := s.db.Exec(
		"INSERT INTO events (id, user_id, title, description, start_date, end_date, notify_before) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		event.ID, event.UserId, event.Title, event.Description, event.StartDate, event.EndDate, int(notifyBefore),
	); err != nil {
		return errors.Join(err, storage.ErrEventCantBeAdded)
	}
	return nil
}

func (s *Storage) DeleteEvent(eventUuid uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM events WHERE id = $1", eventUuid)
	if err != nil {
		return errors.Join(err, storage.ErrEventCantBeDeleted)
	}
	return nil
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	if event.HasDifferentDate(event) {
		events := s.GetUserEvents(event.UserId)
		idx := -1
		for i, e := range events {
			if e.ID == event.ID {
				idx = i
				break
			}
		}
		if idx != -1 {
			events = append(events[:idx], events[idx+1:]...)
		}
		if err := logic.CheckEvent(event, events); err != nil {
			return errors.Join(err, storage.ErrEventCantBeUpdated)
		}
	}
	_, err := s.db.Exec(
		"UPDATE events SET title = $1, description = $2, start_date = $3, end_date = $4 WHERE id = $5",
		event.Title, event.Description, event.StartDate, event.EndDate, event.ID,
	)
	if err != nil {
		return errors.Join(err, storage.ErrEventCantBeUpdated)
	}
	return nil
}

func (s *Storage) GetEvents(fromDate, toDate time.Time) (storage.Schedule, error) {
	rows, err := s.db.Query("SELECT * FROM events WHERE start_date >= $1 AND end_date <= $2", fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events storage.Schedule
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event.ID, &event.UserId, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.NotifyBefore); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) GetUserEvents(userId uuid.UUID) storage.Schedule {
	rows, err := s.db.Query("SELECT * FROM events WHERE user_id = $1", userId)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var events storage.Schedule
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event.ID, &event.UserId, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.NotifyBefore); err != nil {
			return nil
		}
		events = append(events, event)
	}
	return events
}