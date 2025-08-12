// Package createaccount provides the use case for creating an account.
package createaccount

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	accountrepository "github.com/nahtann/controlriver.com/internal/domain/accounts/repository"
	createaccountrequest "github.com/nahtann/controlriver.com/internal/domain/accounts/use_cases/create_account/request"
	createaccountresponse "github.com/nahtann/controlriver.com/internal/domain/accounts/use_cases/create_account/response"
	accountroles "github.com/nahtann/controlriver.com/internal/helpers/accounts/roles"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

var (
	ErrFailedToCreateAccount    = fmt.Errorf("failed to create account")
	ErrFailedToAddAccountMember = fmt.Errorf("failed to add account member")
)

type CreateAccountUseCase struct {
	DB                *pgxpool.Pool
	AccountRepository accountrepository.AccountRepository
}

func NewCreateAccountUseCase(db *pgxpool.Pool, accountRepository accountrepository.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		DB:                db,
		AccountRepository: accountRepository,
	}
}

func (uc *CreateAccountUseCase) Execute(userID uuid.UUID, request createaccountrequest.CreateAccountRequest) (*createaccountresponse.Response, error) {
	accountID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	ctx := context.Background()

	tx, err := uc.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}
	defer tx.Rollback(ctx)

	qtx := uc.AccountRepository.WithTx(tx)

	if err := qtx.CreateAccount(ctx, repository.CreateAccountParams{
		AccountID:       accountID,
		CreatedByUserID: userID,
		Name:            request.Name,
		Description:     &request.Description,
		CurrencyCode:    &request.CurrencyCode,
	}); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	if err := qtx.AddAccountMember(ctx, repository.AddAccountMemberParams{
		AccountID: accountID,
		UserID:    userID,
		Role:      accountroles.Owner,
	}); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToAddAccountMember, err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	return &createaccountresponse.Response{
		AccountID: accountID.String(),
	}, nil
}
