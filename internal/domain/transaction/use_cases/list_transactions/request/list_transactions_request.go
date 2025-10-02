// Package listtransactionsrequest implements the request structure for listing transactions.
package listtransactionsrequest

import "github.com/google/uuid"

type Request struct {
	AccountID  uuid.UUID `json:"account_id" validate:"required"`
	Type       string    `json:"type"`
	AddedBy    uuid.UUID `json:"added_by"`
	From       string    `json:"from" validate:"omitempty,datetime=2006-01-02"`
	To         string    `json:"to" validate:"omitempty,datetime=2006-01-02"`
	Categories []string  `json:"categories"`
}
