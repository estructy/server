package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/estructy/server/internal/helpers"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Seed the user
	userID := usersSeed(dbpool)
	// Seed the accounts
	accountID := accountSeed(dbpool, userID)
	// Seed the account members
	accountMembersSeed(dbpool, accountID, userID)
	// Seed the transactions
	transactionsSeed(dbpool, accountID, accountCategoriesSeed(dbpool, accountID, categoriesSeed(dbpool, accountID)))
}

func usersSeed(db *pgxpool.Pool) uuid.UUID {
	query := fmt.Sprintf(`INSERT INTO users (name, email)
	 VALUES ('user_test', '%s') RETURNING user_id;`, gofakeit.Email())

	var userID uuid.UUID
	err := db.QueryRow(context.Background(), query).Scan(&userID)
	if err != nil {
		log.Fatalf("Failed to seed user: %v\n", err)
	} else {
		log.Printf("User seeded successfully with ID: %d\n", userID)
	}

	return userID
}

func accountSeed(db *pgxpool.Pool, userID uuid.UUID) uuid.UUID {
	query := fmt.Sprintf(`INSERT INTO accounts (name, description, currency_code, created_by_user_id) VALUES ('Conta Pessoal', 'Conta para uso pessoal', 'BRL', '%s') RETURNING account_id`, userID)

	var accountID uuid.UUID
	err := db.QueryRow(context.Background(), query).Scan(&accountID)
	if err != nil {
		log.Fatalf("Failed to seed accounts: %v\n", err)
	} else {
		log.Println("Accounts seeded successfully.")
	}

	return accountID
}

func accountMembersSeed(db *pgxpool.Pool, accountID uuid.UUID, userID uuid.UUID) {
	query := fmt.Sprintf(`INSERT INTO account_members (account_id, user_id, role) VALUES ('%s', '%s', 'owner')`, accountID, userID)

	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to seed account members: %v\n", err)
	} else {
		log.Println("Account members seeded successfully.")
	}
}

func categoriesSeed(db *pgxpool.Pool, accountID uuid.UUID) []uuid.UUID {
	accountIDStr := accountID.String()
	accountIDSuffix := accountIDStr[len(accountIDStr)-4:]

	query := fmt.Sprintf(`INSERT INTO categories (name, type)
VALUES ('alimentação_%s', 'expense'), ('transporte_%s', 'expense'), ('entreterimento_%s', 'expense'),
	('educação_%s', 'expense'), ('saúde_%s', 'expense'), ('moradia_%s', 'expense'), ('investimentos_%s', 'income'), 
		('salário_%s', 'income') RETURNING category_id;`,
		accountIDSuffix, accountIDSuffix, accountIDSuffix, accountIDSuffix,
		accountIDSuffix, accountIDSuffix, accountIDSuffix, accountIDSuffix)

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to seed categories: %v\n", err)
	} else {
		log.Println("Categories seeded successfully.")
	}
	defer rows.Close()

	var categoryIDs []uuid.UUID
	for rows.Next() {
		var categoryID uuid.UUID
		if err := rows.Scan(&categoryID); err != nil {
			log.Fatalf("Failed to scan category ID: %v\n", err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v\n", err)
	}

	return categoryIDs
}

func accountCategoriesSeed(db *pgxpool.Pool, accountID uuid.UUID, categoryIDs []uuid.UUID) []uuid.UUID {
	categoryCodeCount := 1
	var parentValues []string

	// range over 4
	for i := 0; i < 4; i++ {
		parentValues = append(parentValues, fmt.Sprintf("('AC-0%d', '%s', '%s', '%s')",
			categoryCodeCount, accountID, categoryIDs[i], gofakeit.HexColor()))
		categoryCodeCount++
	}

	parentQuery := fmt.Sprintf(`INSERT INTO account_categories (category_code, account_id, category_id, color) VALUES %s 
		RETURNING account_category_id;`,
		strings.Join(parentValues, ", "))

	rows, err := db.Query(context.Background(), parentQuery)
	if err != nil {
		log.Fatalf("Failed to seed account categories: %v\n", err)
	} else {
		log.Println("Account categories seeded successfully.")
	}
	defer rows.Close()

	var accountCategoryIDs []uuid.UUID
	for rows.Next() {
		var accountCategoryID uuid.UUID
		if err := rows.Scan(&accountCategoryID); err != nil {
			log.Fatalf("Failed to scan account category ID: %v\n", err)
		}
		accountCategoryIDs = append(accountCategoryIDs, accountCategoryID)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v\n", err)
	}

	var childValues []string

	for i := 0; i < 4; i++ {
		childValues = append(childValues, fmt.Sprintf("('AC-0%d', '%s','%s', '%s', '%s')",
			categoryCodeCount, accountCategoryIDs[i], accountID, categoryIDs[i+4], gofakeit.HexColor()))
		categoryCodeCount++
	}

	childQuery := fmt.Sprintf(`INSERT INTO account_categories (category_code, parent_id, account_id, category_id, color) VALUES %s 
		RETURNING account_category_id;`,
		strings.Join(childValues, ", "))

	rows, err = db.Query(context.Background(), childQuery)
	if err != nil {
		log.Fatalf("Failed to seed account categories: %v\n", err)
	} else {
		log.Println("Account categories seeded successfully.")
	}
	defer rows.Close()

	var accountCategoryIDsChild []uuid.UUID
	for rows.Next() {
		var accountCategoryID uuid.UUID
		if err := rows.Scan(&accountCategoryID); err != nil {
			log.Fatalf("Failed to scan account category ID: %v\n", err)
		}
		accountCategoryIDsChild = append(accountCategoryIDsChild, accountCategoryID)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v\n", err)
	}

	accountCategoryIDs = append(accountCategoryIDs, accountCategoryIDsChild...)

	return accountCategoryIDs
}

func transactionsSeed(db *pgxpool.Pool, accountID uuid.UUID, accountCategoryIDs []uuid.UUID) {
	currentYear := time.Now().Year()
	code := "TR-01" // Starting code for transactions
	rows := make([][]any, 0, 1000)
	for i := 0; i < 1000; i++ {
		rows = append(rows, []any{
			code,
			accountID,
			accountCategoryIDs[gofakeit.Number(0, len(accountCategoryIDs)-1)], // randomly select a category ID
			gofakeit.Price(100, 10000),                                        // random amount between 100 and 10000
			gofakeit.DateRange(
				time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(currentYear, 12, 31, 23, 59, 59, 999999999, time.UTC)).Format("2006-01-02"), // random date between 2024-01-01 and 2025-12-31
		})

		code = helpers.IncrementCode(code)
	}

	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"transactions"},
		[]string{"transaction_code", "account_id", "category_id", "amount", "transaction_date"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Fatalf("Failed to seed transactions: %v\n", err)
	} else {
		log.Printf("Transactions seeded successfully. Total rows copied: %d\n", copyCount)
	}
}
