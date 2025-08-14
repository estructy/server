-- name: CreateTransactionCategory :exec
INSERT INTO transaction_categories (name, type) 
VALUES ($1, $2);

-- name: FindTransactionCategoriesByNames :many
SELECT * FROM transaction_categories 
WHERE name = ANY($1::varchar[]);
