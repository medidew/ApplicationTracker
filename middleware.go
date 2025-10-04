package main

import (
	"net/http"

	"go.uber.org/zap"
)

func ZapLoggerMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			logger.Info("Request complete",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("status", r.Method))
		})
	}
}