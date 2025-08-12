// Package categorieshandler provides a handler for category-related operations.
package categorieshandler

import (
	"fmt"
	"net/http"

	contexthelper "github.com/nahtann/controlriver.com/internal/helpers/context"
)

type CategoriesHandler struct{}

func NewCategoriesHandler() *CategoriesHandler {
	return &CategoriesHandler{}
}

func (uc *CategoriesHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	accountID, ok := contexthelper.AccountIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Account ID not found in context", http.StatusInternalServerError)
		return
	}

	fmt.Println(accountID)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Category created successfully"))
}
