-- name: CreateCategory :exec
INSERT INTO categories (category_id, account_id, name, type, color) 
VALUES ($1, $2, $3, $4, $5);
