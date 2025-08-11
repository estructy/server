// Package usershandler provides handlers for user-related operations.
package usershandler

import (
	"encoding/json"
	"net/http"

	createuser "github.com/nahtann/controlriver.com/internal/domain/user/use_cases/create_user"
	createuserrequest "github.com/nahtann/controlriver.com/internal/domain/user/use_cases/create_user/request"
	jsonhelper "github.com/nahtann/controlriver.com/internal/helpers/json"
	requesthelper "github.com/nahtann/controlriver.com/internal/helpers/request"
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

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody createuserrequest.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		jsonhelper.HTTPError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errorMessages := requesthelper.ValidateRequest(requestBody)
	if errorMessages != "" {
		jsonhelper.HTTPError(w, http.StatusBadRequest, errorMessages)
		return
	}

	userUseCase := createuser.NewCreateUserUseCase(h.repository)

	if err := userUseCase.Execute(requestBody); err != nil {
		errorMappings := map[error]jsonhelper.ErrorMappings{
			createuser.ErrFailedToCreateUser: {Code: http.StatusInternalServerError, Message: "Failed to create user"},
			createuser.ErrUserAlreadyExists:  {Code: http.StatusConflict, Message: "User already exists"},
		}

		jsonhelper.HandleError(w, err, errorMappings)
	}

	w.WriteHeader(http.StatusCreated)
}
