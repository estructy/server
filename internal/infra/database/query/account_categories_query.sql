-- name: AddAccountCategories :copyfrom
INSERT INTO account_categories (account_category_id, category_code, account_id, category_id, color) 
VALUES ($1, $2, $3, $4, $5);

-- name: FindLastAccountCategoryCode :one
SELECT category_code FROM account_categories WHERE account_id = $1 ORDER BY category_code DESC LIMIT 1;
