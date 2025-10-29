// Package createuser implements the use case for creating a user.
package createuser

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	userrepository "github.com/estructy/server/internal/domain/user/repository"
	createuserrequest "github.com/estructy/server/internal/domain/user/use_cases/create_user/request"
	"github.com/estructy/server/internal/infra/database/repository"
)

var (
	ErrUserAlreadyExists  = fmt.Errorf("user already exists")
	ErrFailedToCreateUser = fmt.Errorf("failed to create user")
)

type CreateUserUseCase struct {
	UserRepository userrepository.UserRepository
}

func NewCreateUserUseCase(userRepository userrepository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: userRepository,
	}
}

func (uc *CreateUserUseCase) Execute(request createuserrequest.CreateUserRequest) error {
	userExists, err := uc.UserRepository.UserExistsByEmail(context.Background(), request.Email)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateUser, err.Error())
	}
	if userExists {
		return fmt.Errorf("%w: %s", ErrUserAlreadyExists, request.Email)
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateUser, err.Error())
	}

	_, err = uc.UserRepository.CreateUser(context.Background(), repository.CreateUserParams{
		UserID: userID,
		Name:   request.Name,
		Email:  request.Email,
	})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateUser, err.Error())
	}

	return nil
}
