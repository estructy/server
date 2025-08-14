package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Seed the transaction categories
	transactionCategoriesSeed(dbpool)
}

func transactionCategoriesSeed(db *pgxpool.Pool) {
	query := `INSERT INTO transaction_categories (name, type)
VALUES ('alimentação', 'expense'), ('transporte', 'expense'), ('entreterimento', 'expense'),
	('educação', 'expense'), ('saúde', 'expense'), ('moradia', 'expense'), ('investimentos', 'income'), ('salário', 'income');`

	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to seed transaction categories: %v\n", err)
	} else {
		log.Println("Transaction categories seeded successfully.")
	}
}
