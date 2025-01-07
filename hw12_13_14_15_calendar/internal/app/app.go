package app

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/homerchik/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface {
	AddEvent(event storage.Event) error
	DeleteEvent(uuid uuid.UUID) error
	UpdateEvent(event storage.Event) error
	GetEvents(fromDate, toDate time.Time) (storage.Schedule, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
