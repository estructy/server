// Package accountshandler provides handlers for account-related operations.
package accountshandler

import (
	"encoding/json"
	"net/http"

	createaccount "github.com/nahtann/controlriver.com/internal/domain/accounts/use_cases/create_account"
	createaccountrequest "github.com/nahtann/controlriver.com/internal/domain/accounts/use_cases/create_account/request"
	contexthelper "github.com/nahtann/controlriver.com/internal/helpers/context"
	jsonhelper "github.com/nahtann/controlriver.com/internal/helpers/json"
	requesthelper "github.com/nahtann/controlriver.com/internal/helpers/request"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

type AccountsHandler struct {
	repository *repository.Queries
}

func NewAccountsHandler(repository *repository.Queries) *AccountsHandler {
	return &AccountsHandler{
		repository: repository,
	}
}

func (accounts *AccountsHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
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

	errorMessages := requesthelper.ValidateRequest(requestBody)
	if errorMessages != "" {
		jsonhelper.HTTPError(w, http.StatusBadRequest, errorMessages)
		return
	}

	createAccountUseCase := createaccount.NewCreateAccountUseCase(accounts.repository)
	if err := createAccountUseCase.Execute(userID, requestBody); err != nil {
		errorMappings := map[error]jsonhelper.ErrorMappings{
			createaccount.ErrFailedToCreateAccount: {Code: http.StatusInternalServerError, Message: "Failed to create account"},
		}

		jsonhelper.HandleError(w, err, errorMappings)
		return
	}

	jsonhelper.HTTPResponse(w, http.StatusCreated, map[string]string{
		"message": "Account created successfully",
	})
}
