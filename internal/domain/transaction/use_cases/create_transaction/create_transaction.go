// Package createtransaction provides the use case for creating a transaction.
package createtransaction

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	createtransactionrequest "github.com/nahtann/controlriver.com/internal/domain/transaction/use_cases/create_transaction/request/create_transaction_request"
	"github.com/nahtann/controlriver.com/internal/helpers"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

var (
	ErrFailedToCreateTransaction = fmt.Errorf("failed to create transaction")
	ErrCategoryNotFound          = fmt.Errorf("category not found")

	initialCode = "TR-000"
)

type CreateTransactionUseCase struct {
	repository *repository.Queries
}

func NewCreateTransactionUseCase(repository *repository.Queries) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		repository: repository,
	}
}

func (uc *CreateTransactionUseCase) Execute(accountID uuid.UUID, request createtransactionrequest.Request) error {
	ctx := context.Background()

	code, err := uc.repository.FindLastTransactionCode(ctx, &accountID)
	if err != nil {
		if err.Error() != "no rows in result set" {
			return fmt.Errorf("%w: %s", ErrFailedToCreateTransaction, err.Error())
		}

		code = initialCode
	}

	newCode := helpers.IncrementCode(code)
	transactionID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateTransaction, err.Error())
	}
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return fmt.Errorf("invalid transaction date format: %w", err)
	}

	category, err := uc.repository.FindAccountCategoryByCode(ctx, repository.FindAccountCategoryByCodeParams{
		AccountID:    &accountID,
		CategoryCode: &request.CategoryCode,
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			return fmt.Errorf("%w: %s", ErrCategoryNotFound, err.Error())
		}
		return fmt.Errorf("%w: %s", ErrFailedToCreateTransaction, err.Error())
	}

	err = uc.repository.CreateTransaction(ctx, repository.CreateTransactionParams{
		TransactionID:     transactionID,
		Code:              newCode,
		AccountID:         &accountID,
		AccountCategoryID: &category.AccountCategoryID,
		Amount:            int32(request.Amount),
		Description:       &request.Description,
		Date:              date,
	})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateTransaction, err.Error())
	}

	return nil
}
