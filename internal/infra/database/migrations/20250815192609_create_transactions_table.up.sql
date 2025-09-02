CREATE TABLE IF NOT EXISTS transactions (
	transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	code VARCHAR(20) NOT NULL,
  account_id UUID REFERENCES accounts(account_id) ON DELETE CASCADE,
  account_category_id UUID REFERENCES account_categories(account_category_id),
  amount INT NOT NULL CHECK (amount >= 0),
  description TEXT,
  date DATE NOT NULL,
  version INT DEFAULT 1,
	added_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
  deleted_at TIMESTAMP DEFAULT NULL,
	CHECK (code ~ '^TR-[0-9]+$')
);
