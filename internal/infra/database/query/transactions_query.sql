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
	ac.category_code AS category_code,
	c.name AS category_name,
	c.type AS category_type,
	t.transaction_date,
	t.amount, 
	t.description, 
	t.created_at
FROM transactions t
LEFT JOIN account_categories ac ON 
	t.category_id = ac.category_id
LEFT JOIN categories c ON 
	ac.category_id = c.category_id
WHERE transaction_id = $1;
