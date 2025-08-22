-- name: FindLastTransactionCode :one 
SELECT code FROM transactions WHERE account_id = $1 ORDER BY code DESC LIMIT 1;

-- name: CreateTransaction :exec
INSERT INTO transactions (
	transaction_id,
	code, 
	account_id,
	account_category_id, 
	amount,
	description, 
	date
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
); 

-- name: FindTransactionById :one
SELECT 
	t.code as transaction_code, 
	ac.category_code AS category_code,
	c.name AS category_name,
	c.type AS category_type,
	t.date as transaction_date,
	t.amount, 
	t.description, 
	t.created_at
FROM transactions t
LEFT JOIN account_categories ac ON 
	t.account_category_id = ac.account_category_id
LEFT JOIN categories c ON 
	ac.category_id = c.category_id
WHERE transaction_id = $1;

-- name: FindTransactionsByType :many
SELECT 
	c.name,
  t.amount, 
	t.description, 
	t.date as transaction_date 
FROM transactions t
LEFT JOIN account_categories ac ON ac.account_category_id = t.account_category_id
LEFT JOIN categories c ON c.category_id = ac.category_id
WHERE 
	t.account_id = $1
	AND c.type = $2
	AND t.date BETWEEN sqlc.arg('from') AND sqlc.arg('to')
ORDER BY 
    c.name ASC,
    t.date ASC;
