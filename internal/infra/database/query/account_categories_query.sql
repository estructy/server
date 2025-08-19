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
