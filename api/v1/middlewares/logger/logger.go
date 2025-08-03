package logger_middleware

import (
	"net/http"
	"time"

	response_helper "github.com/nahtann/controlriver.com/internal/helpers/response"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
	logger *zap.Logger
}

func NewLoggerMiddleware(logger *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

func (m *LoggerMiddleware) Handle(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := response_helper.NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)
		duration := time.Since(start).Round(time.Millisecond)

		m.logger.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("url", r.RequestURI),
			zap.Int("status", lrw.StatusCode),
			zap.String("duration", duration.String()),
		)
	})
}
