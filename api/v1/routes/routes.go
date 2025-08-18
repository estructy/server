// Package routesv1 provides the routing for version 1 of the API.
package routesv1

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	accountshandler "github.com/nahtann/controlriver.com/api/v1/handlers/accounts"
	transactioncategorieshandler "github.com/nahtann/controlriver.com/api/v1/handlers/categories"
	healthhandler "github.com/nahtann/controlriver.com/api/v1/handlers/health"
	transactionshandler "github.com/nahtann/controlriver.com/api/v1/handlers/transactions"
	usershandler "github.com/nahtann/controlriver.com/api/v1/handlers/users"
	middleswares "github.com/nahtann/controlriver.com/api/v1/middlewares"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

func health(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator) {
	router.HandleFunc("GET /health", middlewares.Chain(healthhandler.GetHealth, middlewares.Logger))
}

func users(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator, repository *repository.Queries) {
	usersHandler := usershandler.NewUsersHandler(repository)

	router.HandleFunc("POST /users", middlewares.Chain(usersHandler.CreateUser, middlewares.Logger))
}

func accounts(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator, db *pgxpool.Pool, repository *repository.Queries) {
	accountsHandler := accountshandler.NewAccountsHandler(db, repository)

	router.HandleFunc("POST /accounts", middlewares.Chain(accountsHandler.CreateAccount, middlewares.Logger, middlewares.Auth))
}

func categories(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator, db *pgxpool.Pool, repository *repository.Queries) {
	categoriesHandler := transactioncategorieshandler.NewCategoriesHandler(db, repository)

	router.HandleFunc("POST /categories", middlewares.Chain(categoriesHandler.CreateCategory, middlewares.Logger, middlewares.Account))
}

func transactions(router *http.ServeMux, middlewares *middleswares.MiddlewareOrchestrator, repository *repository.Queries) {
	transactionsHandler := transactionshandler.NewTransactionsHandler(repository)

	router.HandleFunc("POST /transactions", middlewares.Chain(transactionsHandler.CreateTransaction, middlewares.Logger, middlewares.Account))
}

func NewRouterV1(middlewares *middleswares.MiddlewareOrchestrator, db *pgxpool.Pool, repository *repository.Queries) *http.ServeMux {
	mux := http.NewServeMux()

	health(mux, middlewares)
	users(mux, middlewares, repository)
	accounts(mux, middlewares, db, repository)
	categories(mux, middlewares, db, repository)
	transactions(mux, middlewares, repository)

	return mux
}
