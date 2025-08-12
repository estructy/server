CREATE TABLE IF NOT EXISTS categories (
	category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	account_id UUID REFERENCES accounts(account_id) ON DELETE CASCADE, -- NULL if global
	name VARCHAR(50) NOT NULL,
	type VARCHAR(10) CHECK (type IN ('income', 'expense')) NOT NULL,
	color VARCHAR(7) DEFAULT '#FFFFFF', 
	deleted_at TIMESTAMPTZ DEFAULT NULL
);
