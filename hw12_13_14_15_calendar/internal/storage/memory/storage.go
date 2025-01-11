package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/homerchik/hw12_13_14_15_calendar/internal/storage"
	"github.com/homerchik/hw12_13_14_15_calendar/internal/logic"
)



type Storage struct {
	mu *sync.RWMutex
	Events map[time.Time]storage.Schedule
}

func New() *Storage {
	mu := &sync.RWMutex{}
	return &Storage{mu, make(map[time.Time]storage.Schedule)}
}

func (s *Storage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := logic.CheckEvent(event, s.GetUserEvents(event.UserId)); err != nil {
		return errors.Join(err, storage.ErrEventCantBeAdded)
	}
	s.Events[event.StartDate] = append(s.Events[event.StartDate], event)
	return nil
}

func (s *Storage) DeleteEvent(eventUuid uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if idx, event, err := s.GetEvent(eventUuid); err == nil {
		newEvents := append(s.Events[event.StartDate][:idx], s.Events[event.StartDate][idx+1:]...)
		if len(newEvents) == 0 {
			delete(s.Events, event.StartDate)
		} else {
			s.Events[event.StartDate] = newEvents
		}
		return nil
	} else {
		return errors.Join(err, storage.ErrEventCantBeDeleted)
	}
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if idx, e, err := s.GetEvent(event.ID); err == nil {
		if e.ID != event.ID || e.UserId != event.UserId {
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
	} else {
		return err
	}
}

func (s *Storage) GetEvents(fromDate, toDate time.Time) (storage.Schedule, error) {
	result := make(storage.Schedule, 0)
	for date, events := range s.Events {
		if date.After(fromDate) && date.Before(toDate) {
			result = append(result, events...)
		}
	}
	return result, nil
}

func (s *Storage) GetUserEvents(userID uuid.UUID) storage.Schedule {
	schedule := make(storage.Schedule, 0)
	for _, events := range s.Events {
		for _, event := range events {
			if event.UserId == userID {
				schedule = append(schedule, event)
			}
		}
	}
	return schedule
}

func (s *Storage) GetEvent(eventUuid uuid.UUID) (int, *storage.Event, error) {
	for _, events := range s.Events {
		for i, event := range events {
			if event.ID == eventUuid {
				return i, &event, nil
			}
		}
	}
	return -1, nil, storage.ErrNoEventFound
}