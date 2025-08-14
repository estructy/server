-- name: AddAccountTransactionCategories :copyfrom
INSERT INTO account_transaction_categories (category_code, account_id, transaction_category_id, color) VALUES ($1, $2, $3, $4);

-- name: FindLastAccountTransactionCategoryCode :one
SELECT category_code FROM account_transaction_categories WHERE account_id = $1 ORDER BY category_code DESC LIMIT 1;
