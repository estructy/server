// Package createaccountrequest provides the structure for creating a new account request.
package createaccountrequest

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	currencycodes "github.com/nahtann/controlriver.com/internal/helpers/currency/codes"
)

type CreateAccountRequest struct {
	Name         string `json:"name" validate:"required,max=100"`
	Description  string `json:"description" validate:"max=500"`
	CurrencyCode string `json:"currency_code" validate:"required,len=3,uppercase"`
}

func Validate(request *CreateAccountRequest) []string {
	validate := validator.New()

	errorMessages := []string{}
	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		// @todo: improve error messages response
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, e.Error())
		}
	}

	if !currencycodes.IsValid(request.CurrencyCode) {
		errorMessages = append(errorMessages, fmt.Sprintf("Invalid currency code: %s", request.CurrencyCode))
	}

	return errorMessages
}
