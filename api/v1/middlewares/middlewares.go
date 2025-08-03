package middleswares

import (
	"net/http"

	logger_middleware "github.com/nahtann/controlriver.com/api/v1/middlewares/logger"
	"go.uber.org/zap"
)

type Middleware func(http.Handler) http.HandlerFunc

type MiddlewareOrchestrator struct {
	Logger func(http.Handler) http.HandlerFunc
}

func NewMiddlewareOrchestration(logger *zap.Logger) *MiddlewareOrchestrator {
	loggerMiddleware := logger_middleware.NewLoggerMiddleware(logger)

	return &MiddlewareOrchestrator{
		Logger: loggerMiddleware.Handle,
	}
}

func (mo *MiddlewareOrchestrator) Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
