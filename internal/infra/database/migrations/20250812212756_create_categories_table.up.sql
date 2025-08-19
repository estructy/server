CREATE TABLE IF NOT EXISTS categories (
	category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	parent_id UUID REFERENCES categories(transaction_category_id) ON DELETE SET NULL,
	name VARCHAR(50) NOT NULL,
	type VARCHAR(10) CHECK (type IN ('income', 'expense')) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	UNIQUE (name, type)
);
