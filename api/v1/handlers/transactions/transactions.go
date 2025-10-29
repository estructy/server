// Package transactionshandler provides handlers for transaction-related operations.
package transactionshandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	createtransaction "github.com/estructy/server/internal/domain/transaction/use_cases/create_transaction"
	createtransactionrequest "github.com/estructy/server/internal/domain/transaction/use_cases/create_transaction/request/create_transaction_request"
	listtransactions "github.com/estructy/server/internal/domain/transaction/use_cases/list_transactions"
	listtransactionsrequest "github.com/estructy/server/internal/domain/transaction/use_cases/list_transactions/request"
	contexthelper "github.com/estructy/server/internal/helpers/context"
	jsonhelper "github.com/estructy/server/internal/helpers/json"
	requesthelper "github.com/estructy/server/internal/helpers/request"
	"github.com/estructy/server/internal/infra/database/repository"
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
	userID, ok := contexthelper.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusBadRequest)
		return
	}

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
	err := createTransactionUseCase.Execute(userID, accountID, request)
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

func (uc *TransactionsHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	accountID, ok := contexthelper.AccountIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Account ID not found in context", http.StatusInternalServerError)
		return
	}

	transactionType := r.URL.Query().Get("type")
	addedByRaw := r.URL.Query().Get("added_by")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	categoriesRaw := r.URL.Query().Get("categories")

	var categories []string

	if categoriesRaw != "" {
		categories = strings.Split(categoriesRaw, ",")
	}

	addedBy, err := uuid.Parse(addedByRaw)
	if addedByRaw != "" && err != nil {
		http.Error(w, "Invalid added_by UUID", http.StatusBadRequest)
		return
	}

	request := &listtransactionsrequest.Request{
		AccountID:  accountID,
		Type:       transactionType,
		AddedBy:    addedBy,
		From:       from,
		To:         to,
		Categories: categories,
	}

	errorMessages := requesthelper.ValidateRequest(request)
	if errorMessages != "" {
		jsonhelper.HTTPError(w, http.StatusBadRequest, errorMessages)
		return
	}

	if request.Type == "all" {
		request.Type = ""
	}

	listTransactionsUseCase := listtransactions.NewListTransactionsUseCase(uc.repository)
	response, err := listTransactionsUseCase.Execute(request)
	if err != nil {
		errMappings := map[error]jsonhelper.ErrorMappings{
			createtransaction.ErrFailedToCreateTransaction: {
				Code:    http.StatusInternalServerError,
				Message: "Failed to list transactions",
			},
		}
		jsonhelper.HandleError(w, err, errMappings)
		return
	}

	jsonhelper.HTTPResponse(w, http.StatusOK, response)
}
