-- name: AddAccountTransactionCategories :copyfrom
INSERT INTO account_transaction_categories (category_code, account_id, transaction_category_id, color) VALUES ($1, $2, $3, $4);
