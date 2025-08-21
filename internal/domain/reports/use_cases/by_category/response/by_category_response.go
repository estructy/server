// Package bycategoryresponse provides a structure for representing a categorized spending report.
package bycategoryresponse

import "github.com/nahtann/controlriver.com/internal/infra/database/repository"

type Category struct {
	Name        string     `json:"name"`
	TotalSpent  uint64     `json:"total_spent"`
	SubCategory []Category `json:"sub_category,omitempty"`
}

type Response struct {
	From       string     `json:"from"`
	To         string     `json:"to"`
	Type       string     `json:"type"`
	Categories []Category `json:"categories"`
}

func NewResponse(from, to, reportType string, rawCategoriesRows []repository.GetReportByCategoriesRow) *Response {
	response := &Response{
		From: from,
		To:   to,
		Type: reportType,
	}

	var categories []Category
	for _, row := range rawCategoriesRows {
		category := Category{
			Name:       *row.Name,
			TotalSpent: uint64(row.TotalSpent),
		}
		categories = append(categories, category)
	}
	response.Categories = categories

	return response
}

func NewResponseWithSubCategories(from, to, reportType string, rawCategoriesRows []repository.GetReportByCategoriesWithParentCategoryRow) *Response {
	response := &Response{
		From: from,
		To:   to,
		Type: reportType,
	}

	var categories []Category
	for _, row := range rawCategoriesRows {
		category := Category{
			Name:       *row.Name,
			TotalSpent: uint64(row.TotalSpent),
		}

		if row.Parent == *row.Name {
			categories = append(categories, category)
			continue
		}

		for i := range categories {
			if categories[i].Name == row.Parent {
				categories[i].SubCategory = append(categories[i].SubCategory, category)
				break
			}
		}
	}
	response.Categories = categories

	return response
}
