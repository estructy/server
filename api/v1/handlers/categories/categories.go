// Package transactioncategorieshandler provides a handler for category-related operations.
package transactioncategorieshandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	createtransactioncategory "github.com/nahtann/controlriver.com/internal/domain/transaction_categories/use_cases/create_transaction_category"
	createtransactioncategoryrequest "github.com/nahtann/controlriver.com/internal/domain/transaction_categories/use_cases/create_transaction_category/request"
	contexthelper "github.com/nahtann/controlriver.com/internal/helpers/context"
	jsonhelper "github.com/nahtann/controlriver.com/internal/helpers/json"
	requesthelper "github.com/nahtann/controlriver.com/internal/helpers/request"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

type CategoriesHandler struct {
	db         *pgxpool.Pool
	repository *repository.Queries
}

func NewCategoriesHandler(db *pgxpool.Pool, repository *repository.Queries) *CategoriesHandler {
	return &CategoriesHandler{
		db:         db,
		repository: repository,
	}
}

func (h *CategoriesHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	accountID, ok := contexthelper.AccountIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Account ID not found in context", http.StatusInternalServerError)
		return
	}

	var request createtransactioncategoryrequest.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	errorMessages := requesthelper.ValidateRequest(request)
	if errorMessages != "" {
		jsonhelper.HTTPError(w, http.StatusBadRequest, errorMessages)
		return
	}

	createCategoryUseCase := createtransactioncategory.NewCreateCategoryUseCase(h.db, h.repository)
	err := createCategoryUseCase.Execute(accountID, &request)
	if err != nil {
		errMappings := map[error]jsonhelper.ErrorMappings{
			createtransactioncategory.ErrFailedToCreateCategory: {Code: http.StatusInternalServerError, Message: "Failed to create category"},
			createtransactioncategory.ErrCategoryAlreadyExists:  {Code: http.StatusConflict, Message: "Category already exists"},
		}
		jsonhelper.HandleError(w, err, errMappings)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Category created successfully"))
}
