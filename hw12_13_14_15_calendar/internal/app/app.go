package app

import (
	"context"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

type App struct {
	Storage models.Storage
	logger  Logger
	addr    string
}

type Logger interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Error(format string, a ...any)
}

func New(logger Logger, storage models.Storage, addr string) *App {
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
