package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func ZapLoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next_handler http.Handler) http.Handler {
		return http.HandlerFunc(func(response_writer http.ResponseWriter, request *http.Request) {
			wrapped_writer := middleware.NewWrapResponseWriter(response_writer, request.ProtoMajor)

			defer func() {
				status := wrapped_writer.Status()

				if status >= 500 {
					logger.Warn("Request failed",
						zap.String("method", request.Method),
						zap.String("path", request.URL.Path),
						zap.Int("status", status))
				} else {
					logger.Info("Request complete",
						zap.String("method", request.Method),
						zap.String("path", request.URL.Path),
						zap.Int("status", status))
				}
			}()
			
			next_handler.ServeHTTP(wrapped_writer, request)
		})
	}
}