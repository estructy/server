-- name: CreateTransactionCategory :one
INSERT INTO transaction_categories (parent_id, name, type) 
VALUES (NULLIF(sqlc.narg(parent_id), '00000000-0000-0000-0000-000000000000'::uuid), @name, @type) 
ON CONFLICT (name, type) DO NOTHING 
RETURNING transaction_category_id;

-- name: FindTransactionCategoriesByNames :many
SELECT * FROM transaction_categories 
WHERE name = ANY(sqlc.arg(name)::varchar[]);
