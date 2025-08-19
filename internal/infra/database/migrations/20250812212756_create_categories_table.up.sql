CREATE TABLE IF NOT EXISTS categories (
	category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(50) NOT NULL,
	type VARCHAR(10) CHECK (type IN ('income', 'expense')) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	UNIQUE (name, type)
);
