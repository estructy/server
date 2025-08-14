// Package createcategory provides functionality to create a new category.
package createcategory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	createcategoryrequest "github.com/nahtann/controlriver.com/internal/domain/categories/use_cases/create_category/request"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

var ErrFailedToCreateCategory = fmt.Errorf("failed to create category")

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

	categoryID, err := qtx.CreateTransactionCategory(ctx, repository.CreateTransactionCategoryParams{
		Name: request.Name,
		Type: request.Type,
	})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateCategory, err.Error())
	}

	// @todo: due to race conditions, could implement a retry mechanism here
	lastCategoryCode, err := qtx.FindLastAccountTransactionCategoryCode(ctx, accountID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateCategory, err.Error())
	}

	newCategoryCode := uc.incrementCategoryCode(*lastCategoryCode)

	_, err = qtx.AddAccountTransactionCategories(ctx, []repository.AddAccountTransactionCategoriesParams{
		{
			Color:                 &request.Color,
			AccountID:             accountID,
			TransactionCategoryID: categoryID,
			CategoryCode:          &newCategoryCode,
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

func (uc *CreateCategoryUseCase) incrementCategoryCode(lastCategoryCode string) string {
	var prefix string
	var number int
	fmt.Sscanf(lastCategoryCode, "%2s-%d", &prefix, &number)
	number++ // Increment the numeric part

	return fmt.Sprintf("%s-%02d", prefix, number) // Format it back to the same structure
}
