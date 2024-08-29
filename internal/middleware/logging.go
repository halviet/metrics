package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		rd *responseData
	}
)

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.rd.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.rd.status = statusCode
}

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := loggingResponseWriter{
				ResponseWriter: w,
				rd: &responseData{
					status: 0,
					size:   0,
				},
			}
			start := time.Now()

			next.ServeHTTP(&lw, r)

			duration := time.Since(start)

			l.Info(
				"Incoming request",
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Duration("duration", duration),
			)
			l.Info(
				"Server response",
				zap.Int("status", lw.rd.status),
				zap.Int("size", lw.rd.size),
			)
		})
	}
}
