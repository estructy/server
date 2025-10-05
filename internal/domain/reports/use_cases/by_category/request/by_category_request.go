// Package bycategoryrequest provides the request structure for generating a report by category.
package bycategoryrequest

type Request struct {
	Type              string `json:"type"`
	From              string `json:"from" validate:"required,datetime=2006-01-02"`
	To                string `json:"to" validate:"required,datetime=2006-01-02,gtefield=From"`
	WithSubCategories bool   `json:"with_sub_categories" validate:"omitempty,boolean"`
	WithTransactions  bool   `json:"with_transactions" validate:"omitempty,boolean"`
}
