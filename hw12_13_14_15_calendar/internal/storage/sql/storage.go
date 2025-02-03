package sqlstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logic"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib" // importing driver for pg
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(_ context.Context, connStr string, driver string) error {
	var err error
	if s.db, err = sql.Open(driver, connStr); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db == nil {
		return nil 
	}
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddEvent(event storage.Event) error {
	events, err := s.GetUserEvents(event.UserID)
	if err != nil {
		return err
	}
	if err := logic.CheckEvent(event, events); err != nil {
		return errors.Join(err, storage.ErrEventCantBeAdded)
	}
	if _, err := s.db.Exec(
		"INSERT INTO events (id, user_id, title, description, start_date, end_date, notify_before)"+
			"VALUES ($1, $2, $3, $4, $5, $6, $7)",
		event.ID, event.UserID, event.Title, event.Description, event.StartDate, event.EndDate, event.NotifyBefore,
	); err != nil {
		return errors.Join(err, storage.ErrEventCantBeAdded)
	}
	return nil
}

func (s *Storage) DeleteEvent(eventUUID uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM events WHERE id = $1", eventUUID)
	if err != nil {
		return errors.Join(err, storage.ErrEventCantBeDeleted)
	}
	return nil
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	if event.HasDifferentDate(event) {
		events, err := s.GetUserEvents(event.UserID)
		if err != nil {
			return err
		}
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
		"UPDATE events SET title = $1, description = $2, start_date = $3, end_date = $4, notify_before = $5 WHERE id = $6",
		event.Title, event.Description, event.StartDate, event.EndDate, event.NotifyBefore, event.ID,
	)
	if err != nil {
		return errors.Join(err, storage.ErrEventCantBeUpdated)
	}
	return nil
}

func (s *Storage) GetEvents(fromDate, toDate int64) (storage.Schedule, error) {
	rows, err := s.db.Query("SELECT * FROM events WHERE start_date >= $1 AND end_date <= $2", fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events storage.Schedule
	for rows.Next() {
		var event storage.Event
		err := rows.Scan(
			&event.ID, &event.UserID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.NotifyBefore,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Storage) GetUserEvents(userID uuid.UUID) (storage.Schedule, error) {
	rows, err := s.db.Query("SELECT * FROM events WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events storage.Schedule
	for rows.Next() {
		event, err := parseEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) GetEventByID(id uuid.UUID) (storage.Event, error) {
	row, err := s.db.Query("SELECT * FROM events WHERE id = $1", id)
	if err != nil {
		return storage.Event{}, err
	}
	if row.Next() {
		event, err := parseEvent(row)
		if err != nil {
			return storage.Event{}, err
		}
		return event, nil
	}
	return storage.Event{}, storage.ErrNoEventFound
}

func parseEvent(rows *sql.Rows) (storage.Event, error) {
	var event storage.Event
	err := rows.Scan(
		&event.ID, &event.UserID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.NotifyBefore,
	)
	if err != nil {
		return storage.Event{}, err
	}
	return event, nil
}
