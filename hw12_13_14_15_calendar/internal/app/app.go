package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/homerchik/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Storage Storage
	logger  Logger
	addr    string
}

type Logger interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Error(format string, a ...any)
}

type Storage interface {
	AddEvent(event storage.Event) error
	DeleteEvent(uuid uuid.UUID) error
	UpdateEvent(event storage.Event) error
	GetEvents(fromDate, toDate time.Time) (storage.Schedule, error)
}

func New(logger Logger, storage Storage, addr string) *App {
	return &App{
		storage, logger, addr,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error { //nolint:revive
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) GetAddr() string {
	return a.addr
}

func (a *App) GetLogger() Logger {
	return a.logger
}
