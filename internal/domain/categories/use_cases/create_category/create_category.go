// Package createcategory provides functionality to create a new category.
package createcategory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	createcategoryrequest "github.com/nahtann/controlriver.com/internal/domain/categories/use_cases/create_category/request"
	"github.com/nahtann/controlriver.com/internal/helpers"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

var (
	ErrFailedToCreateCategory = fmt.Errorf("failed to create category")
	ErrCategoryAlreadyExists  = fmt.Errorf("category already exists")
)

type CreateCategoryUseCase struct {
	db         *pgxpool.Pool
	repository *repository.Queries
}

func NewCreateCategoryUseCase(db *pgxpool.Pool, repository *repository.Queries) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		db:         db,
		repository: repository,
	}
}

func (uc *CreateCategoryUseCase) Execute(accountID uuid.UUID, request *createcategoryrequest.Request) error {
	ctx := context.Background()

	tx, err := uc.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateCategory, err.Error())
	}
	defer tx.Rollback(ctx)
	qtx := uc.repository.WithTx(tx)

	categoryID, err := qtx.CreateCategory(ctx, repository.CreateCategoryParams{
		Name:     request.Name,
		Type:     request.Type,
		ParentID: request.ParentID,
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			return fmt.Errorf("%w: %s", ErrCategoryAlreadyExists, err.Error())
		}
		return fmt.Errorf("%w: %s", ErrFailedToCreateCategory, err.Error())
	}

	// @todo: due to race conditions, could implement a retry mechanism here
	lastCategoryCode, err := qtx.FindLastAccountCategoryCode(ctx, accountID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateCategory, err.Error())
	}

	newCategoryCode := helpers.IncrementCode(*lastCategoryCode)
	newUUID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateCategory, err.Error())
	}

	_, err = qtx.AddAccountCategories(ctx, []repository.AddAccountCategoriesParams{
		{
			AccountCategoryID: newUUID,
			Color:             &request.Color,
			AccountID:         accountID,
			CategoryID:        categoryID,
			CategoryCode:      &newCategoryCode,
		},
	})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateCategory, err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateCategory, err.Error())
	}

	return nil
}
