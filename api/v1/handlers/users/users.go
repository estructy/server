// Package usershandler provides handlers for user-related operations.
package usershandler

import (
	"net/http"

	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

type UsersHandler struct {
	repository *repository.Queries
}

func NewUsersHandler(repository *repository.Queries) *UsersHandler {
	return &UsersHandler{
		repository: repository,
	}
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {}
