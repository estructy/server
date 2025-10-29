// Package userrepository provides an interface for user-related database operations.
package userrepository

import (
	"context"

	"github.com/estructy/server/internal/infra/database/repository"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(context.Context, repository.CreateUserParams) (uuid.UUID, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}
