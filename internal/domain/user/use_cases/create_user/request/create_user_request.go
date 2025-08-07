// Package createuserrequest provides the request structure for creating a user.
package createuserrequest

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email,max=100"`
}
