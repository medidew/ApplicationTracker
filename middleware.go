package main

import (
	"net/http"

	"go.uber.org/zap"
)

func ZapLoggerMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &NabWriter{ResponseWriter: w, status: 0}

			next.ServeHTTP(sw, r)

			if sw.status >= 500 {
				logger.Errorw("Request failed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sw.status))
			} else {
				logger.Infow("Request complete",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sw.status))
			}
		})
	}
}

// wrapper for ResponseWriter to nab the status code of the response
type NabWriter struct {
	http.ResponseWriter
	status int
}

func (w *NabWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}