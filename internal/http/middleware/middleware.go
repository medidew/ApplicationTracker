package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func ZapLoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next_handler http.Handler) http.Handler {
		return http.HandlerFunc(func(response_writer http.ResponseWriter, request *http.Request) {
			nab_writer := &NabWriter{response_writer: response_writer, status: 0}

			next_handler.ServeHTTP(nab_writer, request)

			if nab_writer.status >= 500 {
				logger.Warn("Request failed",
					zap.String("method", request.Method),
					zap.String("path", request.URL.Path),
					zap.Int("status", nab_writer.status))
			} else {
				logger.Info("Request complete",
					zap.String("method", request.Method),
					zap.String("path", request.URL.Path),
					zap.Int("status", nab_writer.status))
			}
		})
	}
}

// wrapper for ResponseWriter to nab the status code of the response
type NabWriter struct {
	response_writer http.ResponseWriter
	status          int
}

func (nab_writer *NabWriter) WriteHeader(code int) {
	nab_writer.status = code
	nab_writer.response_writer.WriteHeader(code)
}

func (nab_writer *NabWriter) Header() http.Header {
	return nab_writer.response_writer.Header()
}

func (nab_writer *NabWriter) Write(b []byte) (int, error) {
	return nab_writer.response_writer.Write(b)
}
