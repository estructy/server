// Package bycategoryresponse provides a structure for representing a categorized spending report.
package bycategoryresponse

import "github.com/nahtann/controlriver.com/internal/infra/database/repository"

type Transactions struct {
	Amount      uint64  `json:"amount"`
	Description *string `json:"description"`
	Date        string  `json:"date"`
}

type Category struct {
	Name         string         `json:"name"`
	TotalSpent   uint64         `json:"total_spent"`
	Transactions []Transactions `json:"transactions,omitempty"`
	SubCategory  []Category     `json:"sub_category,omitempty"`
}

type Response struct {
	From       string     `json:"from"`
	To         string     `json:"to"`
	Type       string     `json:"type"`
	Categories []Category `json:"categories"`
}

func NewResponse(
	from, to, reportType string,
	rawCategoriesRows []repository.GetReportByCategoriesRow,
	transactions []repository.FindTransactionsByTypeRow,
) *Response {
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

	if len(transactions) > 0 {
		for i := range response.Categories {
			for _, transaction := range transactions {
				if response.Categories[i].Name == *transaction.Name {
					response.Categories[i].Transactions = append(response.Categories[i].Transactions, Transactions{
						Amount:      uint64(transaction.Amount),
						Description: transaction.Description,
						Date:        transaction.TransactionDate.Format("2006-01-02"),
					})
				}
			}
		}
	}

	return response
}

func NewResponseWithSubCategories(
	from, to, reportType string,
	rawCategoriesRows []repository.GetReportByCategoriesWithParentCategoryRow,
	transactions []repository.FindTransactionsByTypeRow,
) *Response {
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

	if len(transactions) > 0 {
		for i := range response.Categories {
			for j := range response.Categories[i].SubCategory {
				for _, transaction := range transactions {
					if response.Categories[i].SubCategory[j].Name == *transaction.Name {
						response.Categories[i].SubCategory[j].Transactions = append(response.Categories[i].SubCategory[j].Transactions, Transactions{
							Amount:      uint64(transaction.Amount),
							Description: transaction.Description,
							Date:        transaction.TransactionDate.Format("2006-01-02"),
						})
					}
					if response.Categories[i].Name == *transaction.Name {
						response.Categories[i].Transactions = append(response.Categories[i].Transactions, Transactions{
							Amount:      uint64(transaction.Amount),
							Description: transaction.Description,
							Date:        transaction.TransactionDate.Format("2006-01-02"),
						})
					}
				}
			}
		}
	}

	return response
}
