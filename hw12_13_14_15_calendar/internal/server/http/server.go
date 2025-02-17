package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

type Server struct {
	addr   string
	ctx    context.Context
	logger Logger
	mux    *http.ServeMux
	server *http.Server
}

type Logger interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Error(format string, a ...any)
}

func NewServer(addr string, logger Logger, storage models.Storage) *Server {
	mux := http.NewServeMux()
	handler := &CalendarHandler{logger, storage}
	mux.HandleFunc("/", loggingMiddleware(logger, handler.Hello))
	mux.HandleFunc("POST /api/event", loggingMiddleware(logger, handler.CreateEvent))
	mux.HandleFunc("GET /api/events/{id}", loggingMiddleware(logger, handler.GetEvent))
	mux.HandleFunc("GET /api/events", loggingMiddleware(logger, handler.GetEventsForRange))
	mux.HandleFunc("PUT /api/events/{id}", loggingMiddleware(logger, handler.UpdateEvent))
	mux.HandleFunc("DELETE /api/events/{id}", loggingMiddleware(logger, handler.DeleteEvent))
	return &Server{addr, nil, logger, mux, nil}
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:              s.addr,
		Handler:           s.mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	s.server = server
	s.logger.Info("%v", s.addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("server shutdown error: %v", err)
		return err
	}
	s.logger.Info("server stopped.")
	return nil
}
