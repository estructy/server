// Package listtransactions implements the use case for listing transactions.
package listtransactions

import (
	"context"
	"fmt"
	"time"

	listtransactionsrequest "github.com/estructy/server/internal/domain/transaction/use_cases/list_transactions/request"
	listtransactionsresponse "github.com/estructy/server/internal/domain/transaction/use_cases/list_transactions/response"
	"github.com/estructy/server/internal/infra/database/repository"
)

var (
	ErrFailedToListTransactions = fmt.Errorf("failed to list transactions")
	ErrInvalidDateFormat        = fmt.Errorf("invalid date format, expected YYYY-MM-DD")
)

type ListTransactionsUseCase struct {
	repository *repository.Queries
}

func NewListTransactionsUseCase(repository *repository.Queries) *ListTransactionsUseCase {
	return &ListTransactionsUseCase{repository: repository}
}

func (uc *ListTransactionsUseCase) Execute(request *listtransactionsrequest.Request) (*listtransactionsresponse.Response, error) {
	fromDate, err := time.Parse("2006-01-02", request.From)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDateFormat, request.From)
	}
	toDate, err := time.Parse("2006-01-02", request.To)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDateFormat, request.To)
	}

	ctx := context.Background()
	transactions, err := uc.repository.FindTransactions(ctx, repository.FindTransactionsParams{
		AccountID:  &request.AccountID,
		From:       fromDate,
		To:         toDate,
		Type:       &request.Type,
		Categories: request.Categories,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrFailedToListTransactions, err)
	}

	return listtransactionsresponse.NewResponse(request.From, request.To, request.Type, transactions), nil
}
