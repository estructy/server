// Package createcategoryrequest defines the request structure for creating a new category in the application.
package createcategoryrequest

type Request struct {
	Name  string `json:"name" validate:"required"`
	Type  string `json:"type" validate:"required,oneof=income expense"`
	Color string `json:"color" validate:"required,hexcolor"`
}
