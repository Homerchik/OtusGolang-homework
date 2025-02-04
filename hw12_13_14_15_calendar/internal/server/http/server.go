package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/app"
)

type Server struct { // TODO
	addr   string
	ctx    context.Context
	logger Logger
	mux    *http.ServeMux
	server *http.Server
}

type Logger interface { // TODO
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Error(format string, a ...any)
}

type Application interface {
	GetAddr() string
	GetLogger() app.Logger
}

type CusHandler struct{}

func (h *CusHandler) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello world!"))
	time.Sleep(time.Second)
}

func NewServer(app Application) *Server {
	mux := http.NewServeMux()
	handler := &CusHandler{}
	mux.HandleFunc("/", loggingMiddleware(app.GetLogger(), handler.Hello))
	return &Server{app.GetAddr(), nil, app.GetLogger(), mux, nil}
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
