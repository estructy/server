// Package createaccount provides the use case for creating an account.
package createaccount

import (
	"context"
	"fmt"

	createaccountrequest "github.com/estructy/server/internal/domain/accounts/use_cases/create_account/request"
	createaccountresponse "github.com/estructy/server/internal/domain/accounts/use_cases/create_account/response"
	accountroles "github.com/estructy/server/internal/helpers/accounts/roles"
	"github.com/estructy/server/internal/infra/database/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFailedToCreateAccount    = fmt.Errorf("failed to create account")
	ErrFailedToAddAccountMember = fmt.Errorf("failed to add account member")

	defaultCategories = []string{
		"alimentação",
		"transporte",
		"entreterimento",
		"educação",
		"saúde",
		"moradia",
		"investimentos",
		"salário",
	}
	colors = []string{
		"#FF5733", // Red
		"#33FF57", // Green
		"#3357FF", // Blue
		"#F1C40F", // Yellow
		"#9B59B6", // Purple
		"#E67E22", // Orange
		"#1ABC9C", // Teal
		"#34495E", // Dark Blue
	}

	categoryCodePrefix = "AC"
)

type CreateAccountUseCase struct {
	DB   *pgxpool.Pool
	repo *repository.Queries
}

func NewCreateAccountUseCase(db *pgxpool.Pool, repo *repository.Queries) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		DB:   db,
		repo: repo,
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

	qtx := uc.repo.WithTx(tx)

	if err := qtx.CreateAccount(ctx, repository.CreateAccountParams{
		AccountID:       accountID,
		CreatedByUserID: userID,
		Name:            request.AccountName,
		Description:     &request.Description,
		CurrencyCode:    &request.CurrencyCode,
	}); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	if err := qtx.UpdateUserLastAccessedAccount(ctx, repository.UpdateUserLastAccessedAccountParams{
		UserID:              userID,
		LastAccessedAccount: &accountID,
	}); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	if request.UserName != "" {
		if err := qtx.UpdateUserName(ctx, repository.UpdateUserNameParams{
			UserID: userID,
			Name:   request.UserName,
		}); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
		}
	}

	if err := qtx.AddAccountMember(ctx, repository.AddAccountMemberParams{
		AccountID: accountID,
		UserID:    userID,
		Role:      accountroles.Owner,
	}); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToAddAccountMember, err.Error())
	}

	// @todo: Implmement cache for transaction categories.
	categories, err := qtx.FindCategoriesByNames(ctx, defaultCategories)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	accountCategories := []repository.AddAccountCategoriesParams{}
	for index, category := range categories {
		categoryCode := fmt.Sprintf("%s-%02d", categoryCodePrefix, index+1)
		color := colors[index%len(colors)]
		newUUID, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
		}

		accountCategories = append(accountCategories, repository.AddAccountCategoriesParams{
			AccountCategoryID: newUUID,
			CategoryCode:      &categoryCode,
			AccountID:         &accountID,
			CategoryID:        &category.CategoryID,
			Color:             &color,
			ParentID:          nil,
		})
	}

	_, err = qtx.AddAccountCategories(ctx, accountCategories)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateAccount, err.Error())
	}

	// @todo: Implement cache update for account ID.

	return &createaccountresponse.Response{
		AccountID: accountID.String(),
	}, nil
}
