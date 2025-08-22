// Package bycategoryresponse provides a structure for representing a categorized spending report.
package bycategoryresponse

import (
	"fmt"

	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

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
	withSubCategories bool,
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
		if withSubCategories {
			handleSubCategories(row, &categories)
			continue
		}

		category := Category{
			Name:        *row.Name,
			TotalSpent:  uint64(row.TotalSpent),
			SubCategory: []Category{},
		}

		categories = append(categories, category)
	}
	response.Categories = categories

	if len(transactions) > 0 {
		for i := range response.Categories {
			if withSubCategories {
				handleSubcategoriesTransactions()
			}

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

func handleSubCategories(row repository.GetReportByCategoriesRow, categories *[]Category) {
	// check if parent category already exists
	var parentCategory *Category
	for i := range *categories {
		if (*categories)[i].Name == row.Parent {
			parentCategory = &(*categories)[i]
			break
		}
	}

	// if parent category does not exist, create it
	if parentCategory == nil {
		fmt.Println("Creating new parent category:", row.Parent)
		parentCategory = &Category{
			Name:       row.Parent,
			TotalSpent: uint64(row.TotalSpent),
			SubCategory: []Category{
				{
					Name:       *row.Name,
					TotalSpent: uint64(row.TotalSpent),
				},
			},
		}
		*categories = append(*categories, *parentCategory)
		return
	} else {
		// if parent category exists, update its total spent
		parentCategory.TotalSpent += uint64(row.TotalSpent)
	}

	// check if subcategory already exists
	var subCategory *Category
	for j := range parentCategory.SubCategory {
		if parentCategory.SubCategory[j].Name == *row.Name {
			subCategory = &parentCategory.SubCategory[j]
			break
		}
	}

	// if subcategory does not exist, create it
	if subCategory == nil && row.Name != &row.Parent {
		subCategory = &Category{
			Name:       *row.Name,
			TotalSpent: uint64(row.TotalSpent),
		}
		parentCategory.SubCategory = append(parentCategory.SubCategory, *subCategory)
	} else {
		// if subcategory exists, update its total spent
		subCategory.TotalSpent += uint64(row.TotalSpent)
	}
}

func handleSubcategoriesTransactions() {
	// for i := range response.Categories {
	// 	for j := range response.Categories[i].SubCategory {
	// 		for _, transaction := range transactions {
	// 			if response.Categories[i].SubCategory[j].Name == *transaction.Name {
	// 				response.Categories[i].SubCategory[j].Transactions = append(response.Categories[i].SubCategory[j].Transactions, Transactions{
	// 					Amount:      uint64(transaction.Amount),
	// 					Description: transaction.Description,
	// 					Date:        transaction.TransactionDate.Format("2006-01-02"),
	// 				})
	// 			}
	// 			if response.Categories[i].Name == *transaction.Name {
	// 				response.Categories[i].Transactions = append(response.Categories[i].Transactions, Transactions{
	// 					Amount:      uint64(transaction.Amount),
	// 					Description: transaction.Description,
	// 					Date:        transaction.TransactionDate.Format("2006-01-02"),
	// 				})
	// 			}
	// 		}
	// 	}
	// }
}
