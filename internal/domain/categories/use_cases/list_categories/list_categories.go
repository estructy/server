// Package listcategories implements the use case for listing categories associated with a specific account.
package listcategories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	listcategoriesrequest "github.com/estructy/server/internal/domain/categories/use_cases/list_categories/request"
	"github.com/estructy/server/internal/infra/database/repository"
)

var ErrFailedToListCategories = fmt.Errorf("failed to list categories")

type ListCategoriesUseCase struct {
	repository *repository.Queries
}

func NewListCategoriesUseCase(repo *repository.Queries) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		repository: repo,
	}
}

func (uc *ListCategoriesUseCase) Execute(accountID *uuid.UUID, request *listcategoriesrequest.Request) ([]repository.FindAccountCategoriesByAccountIDRow, error) {
	ctx := context.Background()

	accountCategories, err := uc.repository.FindAccountCategoriesByAccountID(ctx, repository.FindAccountCategoriesByAccountIDParams{
		AccountID:     accountID,
		Type:          &request.Type,
		WithoutParent: &request.WitoutParent,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToListCategories, err)
	}
	if accountCategories == nil {
		accountCategories = []repository.FindAccountCategoriesByAccountIDRow{}
	}

	return accountCategories, nil
}
