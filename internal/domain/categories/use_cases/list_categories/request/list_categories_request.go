// Package listcategoriesrequest implements the request structure for listing categories.
package listcategoriesrequest

type Request struct {
	Type         string `json:"type" validate:"oneof=income expense"`
	WitoutParent bool   `json:"without_parent" validate:"omitempty"`
}
