// Package createtransactionrequest provides the request structure for creating a transaction.
package createtransactionrequest

import (
	"github.com/google/uuid"
)

type Request struct {
	CategoryID      uuid.UUID `json:"category_id" validate:"required,uuid"`
	Amount          uint      `json:"amount" validate:"required,gte=0"`
	Description     string    `json:"description,omitempty" `
	TransactionDate string    `json:"transaction_date" validate:"required,datetime=2006-01-02"`
}
