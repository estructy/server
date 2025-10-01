// Package listtransactionsresponse implements the response structure for listing transactions.
package listtransactionsresponse

import (
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

type User struct {
	Name string `json:"name"`
}

type Category struct {
	CategoryCode string `json:"category_code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Color        string `json:"color"`
}

type Transaction struct {
	TransactionCode string   `json:"transaction_code"`
	Category        Category `json:"category"`
	Date            string   `json:"date"`
	Amount          uint     `json:"amount"`
	Description     string   `json:"description"`
	AddedBy         User     `json:"added_by"`
}

type Response struct {
	From         string        `json:"from"`
	To           string        `json:"to"`
	Type         string        `json:"type"`
	Transactions []Transaction `json:"transactions"`
}

func NewResponse(from, to, categoryType string, transactions []repository.FindTransactionsRow) *Response {
	response := &Response{
		Transactions: make([]Transaction, 0, len(transactions)),
	}

	response.From = from
	response.To = to
	response.Type = categoryType
	if response.Type == "" {
		response.Type = "all"
	}

	for _, t := range transactions {
		formatedDate := t.TransactionDate.Format("2006-01-02")
		response.Transactions = append(response.Transactions, Transaction{
			TransactionCode: t.TransactionCode,
			Category: Category{
				CategoryCode: *t.CategoryCode,
				Name:         *t.CategoryName,
				Type:         *t.CategoryType,
				Color:        *t.CategoryColor,
			},
			Date:        formatedDate,
			Amount:      uint(t.Amount),
			Description: *t.Description,
			AddedBy: User{
				Name: *t.AddedBy,
			},
		})
	}

	return response
}
