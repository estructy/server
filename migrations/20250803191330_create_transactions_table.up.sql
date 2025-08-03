CREATE TABLE IF NOT EXISTS transactions (
	id UUID PRIMARY KEY NOT NULL,
	user_id UUID NOT NULL,
	category_id INT,
	amount NUMERIC(12,2) NOT NULL,
	type VARCHAR(10) CHECK (type IN ('income', 'expense')) NOT NULL,
	description TEXT,
	transaction_date DATE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
	deleted_at TIMESTAMPTZ 
);
