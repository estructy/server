// Package requesthelper provides utility functions for request validation.
package requesthelper

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateRequest(s any) string {
	validate := validator.New()

	err := validate.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errorMessages := []string{}
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, e.Error())
		}

		return strings.Join(errorMessages, " ")
	}

	return ""
}
