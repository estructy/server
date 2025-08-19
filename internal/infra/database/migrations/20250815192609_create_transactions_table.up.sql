CREATE TABLE IF NOT EXISTS transactions (
	transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	transaction_code VARCHAR(20) NOT NULL,
  account_id UUID REFERENCES accounts(account_id) ON DELETE CASCADE,
  category_id UUID REFERENCES account_categories(account_category_id),
  amount INT NOT NULL CHECK (amount >= 0),
  description TEXT,
  transaction_date DATE NOT NULL,
  version INT DEFAULT 1,
	created_at TIMESTAMPTZ NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  deleted_at TIMESTAMP DEFAULT NULL,
	CHECK (transaction_code ~ '^TR-[0-9]+$')
);
