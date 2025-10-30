-- name: CreateUser :one
INSERT INTO users (user_id, auth_id, name, email) VALUES ($1, $2, $3, $4) RETURNING user_id;

-- name: UserExistsByEmail :one
SELECT EXISTS (
		SELECT 1 FROM users WHERE email = $1
) AS exists;

-- name: GetUserByAuthID :one 
SELECT user_id FROM users WHERE auth_id = $1;

-- name: FindUserLastAccessedAccount :one 
SELECT last_accessed_account FROM users WHERE user_id = $1;

-- name: UpdateUserLastAccessedAccount :exec 
UPDATE users SET last_accessed_account = $2 WHERE user_id = $1; 

-- name: UpdateUserName :exec
UPDATE users SET name = $2 WHERE user_id = $1; 
