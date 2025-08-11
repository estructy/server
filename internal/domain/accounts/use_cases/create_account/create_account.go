// Package createaccount provides the use case for creating an account.
package createaccount

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	accountrepository "github.com/nahtann/controlriver.com/internal/domain/accounts/repository"
	createaccountrequest "github.com/nahtann/controlriver.com/internal/domain/accounts/use_cases/create_account/request"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

var ErrFailedToCreateAccount = fmt.Errorf("failed to create account")

type CreateAccountUseCase struct {
	AccountRepository accountrepository.AccountRepository
}

func NewCreateAccountUseCase(accountRepository accountrepository.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountRepository: accountRepository,
	}
}

func (uc *CreateAccountUseCase) Execute(userID uuid.UUID, request createaccountrequest.CreateAccountRequest) error {
	accountID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	if err := uc.AccountRepository.CreateAccount(context.Background(), repository.CreateAccountParams{
		AccountID:       accountID,
		CreatedByUserID: userID,
		Name:            request.Name,
		Description:     &request.Description,
		CurrencyCode:    &request.CurrencyCode,
	}); err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	return nil
}
