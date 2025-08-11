// Package accountrepository provides the interface for account-related database operations.
package accountrepository

import (
	"context"

	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, arg repository.CreateAccountParams) error
}
