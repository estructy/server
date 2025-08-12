// Package middleswares provides middleware orchestration for the API.
package middleswares

import (
	"net/http"

	accountmiddleware "github.com/nahtann/controlriver.com/api/v1/middlewares/account"
	authmiddleware "github.com/nahtann/controlriver.com/api/v1/middlewares/auth"
	logger_middleware "github.com/nahtann/controlriver.com/api/v1/middlewares/logger"
	"go.uber.org/zap"
)

type Middleware func(http.Handler) http.HandlerFunc

type MiddlewareOrchestrator struct {
	Logger  func(http.Handler) http.HandlerFunc
	Auth    func(http.Handler) http.HandlerFunc
	Account func(http.Handler) http.HandlerFunc
}

func NewMiddlewareOrchestration(logger *zap.Logger) *MiddlewareOrchestrator {
	loggerMiddleware := logger_middleware.NewLoggerMiddleware(logger)
	authMiddleware := authmiddleware.NewAuthMiddleware()
	accountMiddleware := accountmiddleware.NewAccountMiddleware()

	return &MiddlewareOrchestrator{
		Logger:  loggerMiddleware.Handle,
		Auth:    authMiddleware.Handle,
		Account: accountMiddleware.Handle,
	}
}

func (mo *MiddlewareOrchestrator) Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
