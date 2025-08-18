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

	initialTransactionCode = "TR-000"
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

	transactionCode, err := uc.repository.FindLastTransactionCode(ctx, accountID)
	if err != nil {
		if err.Error() != "no rows in result set" {
			return fmt.Errorf("%w: %s", ErrFailedToCreateTransaction, err.Error())
		}

		transactionCode = initialTransactionCode
	}

	newTransactionCode := helpers.IncrementCode(transactionCode)
	transactionID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateTransaction, err.Error())
	}
	transactionDate, err := time.Parse("2006-01-02", request.TransactionDate)
	if err != nil {
		return fmt.Errorf("invalid transaction date format: %w", err)
	}

	err = uc.repository.CreateTransaction(ctx, repository.CreateTransactionParams{
		TransactionID:   transactionID,
		TransactionCode: newTransactionCode,
		AccountID:       accountID,
		CategoryID:      request.CategoryID,
		Amount:          int32(request.Amount),
		Description:     &request.Description,
		TransactionDate: transactionDate,
	})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateTransaction, err.Error())
	}

	return nil
}
