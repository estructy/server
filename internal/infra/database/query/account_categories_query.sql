-- name: AddAccountCategories :copyfrom
INSERT INTO account_categories 
	(account_category_id, category_code, parent_id, account_id, category_id, color) 
VALUES 
	(@account_category_id, 
	@category_code, 
	sqlc.narg('parent_id'),
	@account_id, 
	@category_id, 
	@color);

-- name: FindLastAccountCategoryCode :one
SELECT category_code FROM account_categories WHERE account_id = $1 ORDER BY category_code DESC LIMIT 1;

-- name: FindAccountCategoryByCode :one 
SELECT * FROM account_categories WHERE account_id = $1 AND category_code = $2;

-- name: FindAccountCategoriesByAccountID :many 
SELECT 
	ac.category_code, 
	c.name, 
	c.type,
	ac.color
FROM account_categories ac
LEFT JOIN categories c ON ac.category_id = c.category_id
WHERE 
	account_id = $1 
	AND c.type = COALESCE(NULLIF(sqlc.narg('type')::text, ''), c.type)
	AND (
		NOT sqlc.narg('without_parent')::boolean
		OR ac.parent_id IS NULL
	)
ORDER BY category_code;
