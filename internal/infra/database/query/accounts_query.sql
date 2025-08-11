-- name: CreateAccount :exec
INSERT INTO accounts (account_id, name, description, currency_code, created_by_user_id) 
VALUES ($1, $2, $3, $4, $5);

