package routes_v1

import (
	"net/http"

	health_handler "github.com/nahtann/controlriver.com/api/v1/handlers/health"
	middleswares "github.com/nahtann/controlriver.com/api/v1/middlewares"
)

func health(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator) {
	router.HandleFunc("GET /health", middlewares.Chain(health_handler.GetHealth, middlewares.Logger))
}

func NewRouterV1(middlewares *middleswares.MiddlewareOrchestrator) *http.ServeMux {
	mux := http.NewServeMux()

	health(mux, middlewares)

	return mux
}
