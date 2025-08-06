package usershanlder

import "github.com/nahtann/controlriver.com/internal/infra/database/repository"

type UsersHandler struct {
	repository *repository.Queries
}

func NewUsersHandler(repository *repository.Queries) *UsersHandler {
	return &UsersHandler{
		repository: repository,
	}
}

func (h *UsersHandler) CreateUser() {}
