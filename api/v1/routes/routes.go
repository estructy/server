// Package routesv1 provides the routing for version 1 of the API.
package routesv1

import (
	"net/http"

	healthhandler "github.com/nahtann/controlriver.com/api/v1/handlers/health"
	usershandler "github.com/nahtann/controlriver.com/api/v1/handlers/users"
	middleswares "github.com/nahtann/controlriver.com/api/v1/middlewares"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

func health(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator) {
	router.HandleFunc("GET /health", middlewares.Chain(healthhandler.GetHealth, middlewares.Logger))
}

func users(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator, repository *repository.Queries) {
	usersHandler := usershandler.NewUsersHandler(repository)

	router.HandleFunc("/users", middlewares.Chain(usersHandler.CreateUser, middlewares.Logger))
}

func NewRouterV1(middlewares *middleswares.MiddlewareOrchestrator, repository *repository.Queries) *http.ServeMux {
	mux := http.NewServeMux()

	health(mux, middlewares)
	users(mux, middlewares, repository)

	return mux
}
