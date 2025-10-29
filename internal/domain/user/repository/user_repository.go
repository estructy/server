// Package userrepository provides an interface for user-related database operations.
package userrepository

import (
	"context"

	"github.com/estructy/server/internal/infra/database/repository"
)

type UserRepository interface {
	CreateUser(context.Context, repository.CreateUserParams) (repository.CreateUserRow, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}
