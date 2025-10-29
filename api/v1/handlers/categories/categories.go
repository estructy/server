// Package categorieshandler provides a handler for category-related operations.
package categorieshandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	createcategory "github.com/estructy/server/internal/domain/categories/use_cases/create_category"
	createcategoryrequest "github.com/estructy/server/internal/domain/categories/use_cases/create_category/request"
	listcategories "github.com/estructy/server/internal/domain/categories/use_cases/list_categories"
	listcategoriesrequest "github.com/estructy/server/internal/domain/categories/use_cases/list_categories/request"
	contexthelper "github.com/estructy/server/internal/helpers/context"
	jsonhelper "github.com/estructy/server/internal/helpers/json"
	requesthelper "github.com/estructy/server/internal/helpers/request"
	"github.com/estructy/server/internal/infra/database/repository"
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

	var request createcategoryrequest.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	errorMessages := requesthelper.ValidateRequest(request)
	if errorMessages != "" {
		jsonhelper.HTTPError(w, http.StatusBadRequest, errorMessages)
		return
	}

	createCategoryUseCase := createcategory.NewCreateCategoryUseCase(h.db, h.repository)
	err := createCategoryUseCase.Execute(accountID, &request)
	if err != nil {
		errMappings := map[error]jsonhelper.ErrorMappings{
			createcategory.ErrFailedToCreateCategory: {Code: http.StatusInternalServerError, Message: "Failed to create category"},
			createcategory.ErrCategoryAlreadyExists:  {Code: http.StatusConflict, Message: "Category already exists"},
			createcategory.ErrParentCategoryNotFound: {Code: http.StatusNotFound, Message: "Parent category not found"},
		}
		jsonhelper.HandleError(w, err, errMappings)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Category created successfully"))
}

func (h *CategoriesHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	accountID, ok := contexthelper.AccountIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Account ID not found in context", http.StatusInternalServerError)
		return
	}

	categoriesType := r.URL.Query().Get("type")
	subCategories := r.URL.Query().Get("without-parent")

	request := &listcategoriesrequest.Request{}

	if categoriesType != "" {
		request.Type = categoriesType
	}
	if subCategories != "" && subCategories == "true" {
		request.WitoutParent = true
	} else {
		request.WitoutParent = false
	}

	// errorMessages := requesthelper.ValidateRequest(request)
	// if errorMessages != "" {
	// 	jsonhelper.HTTPError(w, http.StatusBadRequest, errorMessages)
	// 	return
	// }

	listCategoriesUseCase := listcategories.NewListCategoriesUseCase(h.repository)
	categories, err := listCategoriesUseCase.Execute(&accountID, request)
	if err != nil {
		errMappings := map[error]jsonhelper.ErrorMappings{
			listcategories.ErrFailedToListCategories: {Code: http.StatusInternalServerError, Message: "Failed to list categories"},
		}
		jsonhelper.HandleError(w, err, errMappings)
		return
	}

	jsonhelper.HTTPResponse(w, http.StatusOK, categories)
}
