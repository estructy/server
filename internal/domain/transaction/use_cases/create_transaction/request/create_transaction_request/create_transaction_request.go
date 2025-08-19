// Package createtransactionrequest provides the request structure for creating a transaction.
package createtransactionrequest

type Request struct {
	CategoryCode    string `json:"category_code" validate:"required,startswith=AC"`
	Amount          uint   `json:"amount" validate:"required,gte=0"`
	Description     string `json:"description,omitempty" `
	TransactionDate string `json:"transaction_date" validate:"required,datetime=2006-01-02"`
}
