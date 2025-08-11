// Package createaccountrequest provides the structure for creating a new account request.
package createaccountrequest

type CreateAccountRequest struct {
	Name         string `json:"name" validate:"required,max=100"`
	Description  string `json:"description" validate:"max=500"`
	CurrencyCode string `json:"currency_code" validate:"required,len=3,uppercase"`
}
