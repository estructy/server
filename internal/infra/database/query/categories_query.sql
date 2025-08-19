-- name: CreateCategory :one
INSERT INTO categories (name, type) 
VALUES ($1, $2) 
ON CONFLICT (name, type) DO NOTHING 
RETURNING category_id;

-- name: FindCategoriesByNames :many
SELECT * FROM categories 
WHERE name = ANY(sqlc.arg(name)::varchar[]);
