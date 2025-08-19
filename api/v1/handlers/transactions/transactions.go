// Package transactionshandler provides handlers for transaction-related operations.
package transactionshandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	createtransaction "github.com/nahtann/controlriver.com/internal/domain/transaction/use_cases/create_transaction"
	createtransactionrequest "github.com/nahtann/controlriver.com/internal/domain/transaction/use_cases/create_transaction/request/create_transaction_request"
	contexthelper "github.com/nahtann/controlriver.com/internal/helpers/context"
	jsonhelper "github.com/nahtann/controlriver.com/internal/helpers/json"
	requesthelper "github.com/nahtann/controlriver.com/internal/helpers/request"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

type TransactionsHandler struct {
	repository *repository.Queries
}

func NewTransactionsHandler(repository *repository.Queries) *TransactionsHandler {
	return &TransactionsHandler{
		repository: repository,
	}
}

func (uc *TransactionsHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	accountID, ok := contexthelper.AccountIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Account ID not found in context", http.StatusInternalServerError)
		return
	}

	var request createtransactionrequest.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	errorMessages := requesthelper.ValidateRequest(request)
	if errorMessages != "" {
		jsonhelper.HTTPError(w, http.StatusBadRequest, errorMessages)
		return
	}

	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(uc.repository)
	err := createTransactionUseCase.Execute(accountID, request)
	if err != nil {
		errMappings := map[error]jsonhelper.ErrorMappings{
			createtransaction.ErrFailedToCreateTransaction: {
				Code:    http.StatusInternalServerError,
				Message: "Failed to create transaction",
			},
			createtransaction.ErrCategoryNotFound: {
				Code:    http.StatusNotFound,
				Message: "Category not found",
			},
		}
		jsonhelper.HandleError(w, err, errMappings)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Transaction created successfully"))
}
