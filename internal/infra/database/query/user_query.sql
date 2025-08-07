-- name: CreateUser :one
INSERT INTO users (user_id, name, email) VALUES ($1, $2, $3) RETURNING user_id, name, email;
