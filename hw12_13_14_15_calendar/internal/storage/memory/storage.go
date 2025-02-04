package memorystorage

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logic"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

type Storage struct {
	mu     *sync.RWMutex
	Events map[int64]models.Schedule
}

func New() *Storage {
	mu := &sync.RWMutex{}
	return &Storage{mu, make(map[int64]models.Schedule)}
}

func (s *Storage) AddEvent(event models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := logic.CheckEvent(event, s.GetUserEvents(event.UserID)); err != nil {
		return errors.Join(err, models.ErrEventCantBeAdded)
	}
	s.Events[event.StartDate] = append(s.Events[event.StartDate], event)
	return nil
}

func (s *Storage) DeleteEvent(eventUUID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx, event, err := s.GetEventByID(eventUUID)
	if err != nil {
		return errors.Join(err, models.ErrEventCantBeDeleted)
	}
	s.Events[event.StartDate] = append(s.Events[event.StartDate][:idx], s.Events[event.StartDate][idx+1:]...)
	if len(s.Events[event.StartDate]) == 0 {
		delete(s.Events, event.StartDate)
	}
	return nil
}

func (s *Storage) UpdateEvent(event models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx, e, err := s.GetEventByID(event.ID)
	if err != nil {
		return err
	}
	if e.ID != event.ID || e.UserID != event.UserID {
		return models.ErrEventCantBeUpdated
	}
	if e.HasDifferentDate(event) {
		s.mu.Unlock()
		if err := s.DeleteEvent(e.ID); err != nil {
			s.mu.Lock()
			return errors.Join(err, models.ErrEventCantBeUpdated)
		}
		if err := s.AddEvent(event); err != nil {
			s.AddEvent(e) // If we failed to add new event, return old one
			s.mu.Lock()
			return err
		}
		s.mu.Lock()
	} else {
		s.Events[event.StartDate][idx] = event
	}
	return nil
}

func (s *Storage) GetEvents(fromDate, toDate int64) (models.Schedule, error) {
	result := make(models.Schedule, 0)
	for date, events := range s.Events {
		if date > fromDate && date < toDate {
			result = append(result, events...)
		}
	}
	return result, nil
}

func (s *Storage) GetUserEvents(userID uuid.UUID) models.Schedule {
	schedule := make(models.Schedule, 0)
	for _, events := range s.Events {
		for _, event := range events {
			if event.UserID == userID {
				schedule = append(schedule, event)
			}
		}
	}
	return schedule
}

func (s *Storage) GetEventByID(eventUUID uuid.UUID) (int, models.Event, error) {
	for _, events := range s.Events {
		for i, event := range events {
			if event.ID == eventUUID {
				return i, event, nil
			}
		}
	}
	return -1, models.Event{}, models.ErrNoEventFound
}
