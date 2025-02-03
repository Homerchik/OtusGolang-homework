package memorystorage

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logic"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     *sync.RWMutex
	Events map[int64]storage.Schedule
}

func New() *Storage {
	mu := &sync.RWMutex{}
	return &Storage{mu, make(map[int64]storage.Schedule)}
}

func (s *Storage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := logic.CheckEvent(event, s.GetUserEvents(event.UserID)); err != nil {
		return errors.Join(err, storage.ErrEventCantBeAdded)
	}
	s.Events[event.StartDate] = append(s.Events[event.StartDate], event)
	return nil
}

func (s *Storage) DeleteEvent(eventUUID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx, event, err := s.GetEvent(eventUUID)
	if err != nil {
		return errors.Join(err, storage.ErrEventCantBeDeleted)
	}
	s.Events[event.StartDate] = append(s.Events[event.StartDate][:idx], s.Events[event.StartDate][idx+1:]...)
	if len(s.Events[event.StartDate]) == 0 {
		delete(s.Events, event.StartDate)
	}
	return nil
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx, e, err := s.GetEvent(event.ID)
	if err != nil {
		return err
	}
	if e.ID != event.ID || e.UserID != event.UserID {
		return storage.ErrEventCantBeUpdated
	}
	if e.HasDifferentDate(event) {
		s.mu.Unlock()
		if err := s.DeleteEvent(e.ID); err != nil {
			s.mu.Lock()
			return errors.Join(err, storage.ErrEventCantBeUpdated)
		}
		if err := s.AddEvent(event); err != nil {
			s.AddEvent(*e) // If we failed to add new event, return old one
			s.mu.Lock()
			return err
		}
		s.mu.Lock()
	} else {
		s.Events[event.StartDate][idx] = event
	}
	return nil
}

func (s *Storage) GetEvents(fromDate, toDate int64) (storage.Schedule, error) {
	result := make(storage.Schedule, 0)
	for date, events := range s.Events {
		if date > fromDate && date < toDate {
			result = append(result, events...)
		}
	}
	return result, nil
}

func (s *Storage) GetUserEvents(userID uuid.UUID) storage.Schedule {
	schedule := make(storage.Schedule, 0)
	for _, events := range s.Events {
		for _, event := range events {
			if event.UserID == userID {
				schedule = append(schedule, event)
			}
		}
	}
	return schedule
}

func (s *Storage) GetEvent(eventUUID uuid.UUID) (int, *storage.Event, error) {
	for _, events := range s.Events {
		for i, event := range events {
			if event.ID == eventUUID {
				return i, &event, nil
			}
		}
	}
	return -1, nil, storage.ErrNoEventFound
}
