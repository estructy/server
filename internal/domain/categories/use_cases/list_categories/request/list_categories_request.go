// Package listcategoriesrequest implements the request structure for listing categories.
package listcategoriesrequest

type Request struct {
	Type string `json:"type" validate:"required,oneof=income expense"`
}
