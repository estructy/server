-- name: CreateUser :one
INSERT INTO users (user_id, name, email) VALUES ($1, $2, $3) RETURNING user_id, name, email;

-- name: UserExistsByEmail :one
SELECT EXISTS (
		SELECT 1 FROM users WHERE email = $1
) AS exists;
