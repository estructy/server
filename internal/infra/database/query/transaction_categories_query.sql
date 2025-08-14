-- name: CreateTransactionCategory :one
INSERT INTO transaction_categories (name, type) 
VALUES ($1, $2) RETURNING transaction_category_id;

-- name: FindTransactionCategoriesByNames :many
SELECT * FROM transaction_categories 
WHERE name = ANY($1::varchar[]);
