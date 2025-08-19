-- name: CreateCategory :one
INSERT INTO categories (parent_id, name, type) 
VALUES (NULLIF(sqlc.narg(parent_id), '00000000-0000-0000-0000-000000000000'::uuid), @name, @type) 
ON CONFLICT (name, type) DO NOTHING 
RETURNING category_id;

-- name: FindCategoriesByNames :many
SELECT * FROM categories 
WHERE name = ANY(sqlc.arg(name)::varchar[]);
