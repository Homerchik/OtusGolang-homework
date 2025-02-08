package models

import "github.com/google/uuid"

type Storage interface {
	AddEvent(event Event) error
	DeleteEvent(uuid uuid.UUID) error
	UpdateEvent(event Event) error
	GetEvents(fromDate, toDate int64) (Schedule, error)
	GetEventByID(eventUUID uuid.UUID) (int, Event, error)
}
