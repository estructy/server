package dbseeds

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CategoriesUp(db *pgxpool.Pool) {
	query := `INSERT INTO categories (name, type)
VALUES ('alimentação', 'expense'), ('transporte', 'expense'), ('entreterimento', 'expense'),
	('educação', 'expense'), ('saúde', 'expense'), ('moradia', 'expense'), ('investimentos', 'income'), 
	('salário', 'income') ON CONFLICT DO NOTHING;`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to seed categories: %v\n", err)
	} else {
		log.Println("Categories seeded successfully.")
	}
	defer rows.Close()
}
