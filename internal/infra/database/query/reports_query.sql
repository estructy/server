-- name: GetReportByCategories :many
SELECT
  coalesce(pc_c.name, c.name) AS parent,
  c.name,
	sum(t.amount) AS total_spent,
	ac.color AS category_color,
	c.type AS category_type
FROM transactions t
LEFT JOIN account_categories ac ON ac.account_category_id = t.account_category_id
LEFT JOIN categories c ON c.category_id = ac.category_id
LEFT JOIN account_categories pc
    ON pc.account_category_id = ac.parent_id
LEFT JOIN categories pc_c
    ON pc_c.category_id = pc.category_id
WHERE
	t.account_id = $1
  AND t.date BETWEEN sqlc.arg('from') AND sqlc.arg('to')
	AND c.type = COALESCE(NULLIF(sqlc.narg('type')::text, ''), c.type)
GROUP BY c.name, pc_c.name, ac.color, c.type
ORDER BY total_spent DESC;
