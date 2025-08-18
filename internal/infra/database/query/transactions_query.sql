-- name: FindLastTransactionCode :one 
SELECT transaction_code FROM transactions WHERE account_id = $1 ORDER BY transaction_code DESC LIMIT 1;

-- name: CreateTransaction :exec
INSERT INTO transactions (
	transaction_id,
	transaction_code, 
	account_id,
	category_id, 
	amount,
	description, 
	transaction_date
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
); 

-- name: FindTransactionById :one
SELECT 
	t.transaction_code, 
	atc.category_code AS category_code,
	tc.name AS category_name,
	tc.type AS category_type,
	t.transaction_date,
	t.amount, 
	t.description, 
	t.created_at
FROM transactions t
LEFT JOIN account_transaction_categories atc ON 
	t.category_id = atc.transaction_category_id
LEFT JOIN transaction_categories tc ON 
	atc.transaction_category_id = tc.id
WHERE transaction_id = $1;
