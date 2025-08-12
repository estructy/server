// Package accountrepository provides the interface for account-related database operations.
package accountrepository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

type AccountRepository interface {
	WithTx(tx pgx.Tx) *repository.Queries
	CreateAccount(ctx context.Context, arg repository.CreateAccountParams) error
}
