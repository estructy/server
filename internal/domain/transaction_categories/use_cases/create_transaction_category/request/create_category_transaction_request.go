// Package createtransactioncategoryrequest defines the request structure for creating a new transaction category in the application.
package createtransactioncategoryrequest

import "github.com/google/uuid"

type Request struct {
	Name     string    `json:"name" validate:"required"`
	Type     string    `json:"type" validate:"required,oneof=income expense"`
	Color    string    `json:"color" validate:"required,hexcolor"`
	ParentID uuid.UUID `json:"parent_id,omitempty" validate:"omitempty,uuid"`
}
