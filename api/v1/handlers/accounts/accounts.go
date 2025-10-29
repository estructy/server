// Package accountshandler provides handlers for account-related operations.
package accountshandler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	createaccount "github.com/estructy/server/internal/domain/accounts/use_cases/create_account"
	createaccountrequest "github.com/estructy/server/internal/domain/accounts/use_cases/create_account/request"
	contexthelper "github.com/estructy/server/internal/helpers/context"
	jsonhelper "github.com/estructy/server/internal/helpers/json"
	"github.com/estructy/server/internal/infra/database/repository"
)

type AccountsHandler struct {
	db         *pgxpool.Pool
	repository *repository.Queries
}

func NewAccountsHandler(db *pgxpool.Pool, repository *repository.Queries) *AccountsHandler {
	return &AccountsHandler{
		db:         db,
		repository: repository,
	}
}

func (h *AccountsHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := contexthelper.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}

	var requestBody createaccountrequest.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		jsonhelper.HTTPError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errorMessages := createaccountrequest.Validate(&requestBody)
	if len(errorMessages) > 0 {
		jsonhelper.HTTPResponse(w, http.StatusBadRequest, map[string]any{
			"errors": errorMessages,
		})
		return
	}

	// Repository with transaction
	rtx := repository.New(h.db)

	createAccountUseCase := createaccount.NewCreateAccountUseCase(h.db, rtx)
	response, err := createAccountUseCase.Execute(userID, requestBody)
	if err != nil {
		errorMappings := map[error]jsonhelper.ErrorMappings{
			createaccount.ErrFailedToCreateAccount:    {Code: http.StatusInternalServerError, Message: "Failed to create account"},
			createaccount.ErrFailedToAddAccountMember: {Code: http.StatusInternalServerError, Message: "Failed to add account member"},
		}

		jsonhelper.HandleError(w, err, errorMappings)
		return
	}

	jsonhelper.HTTPResponse(w, http.StatusCreated, response)
}
