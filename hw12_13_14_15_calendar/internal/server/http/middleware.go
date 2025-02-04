package internalhttp

import (
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (rec *StatusRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(log Logger, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UnixMilli()
		statusRec := NewStatusRecorder(w)
		next(statusRec, r)
		end := time.Now().UnixMilli()
		log.Info(
			"%v [%v] %v %v %v %v %v %v",
			r.RemoteAddr, time.Now().UTC().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.RequestURI(),
			r.Proto, statusRec.StatusCode, end-start, r.UserAgent(),
		)
	})
}
