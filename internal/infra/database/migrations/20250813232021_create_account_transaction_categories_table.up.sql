CREATE TABLE IF NOT EXISTS account_transaction_categories (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	category_code VARCHAR(20),
	account_id UUID REFERENCES accounts(account_id) ON DELETE CASCADE,
	transaction_category_id UUID REFERENCES transaction_categories(transaction_category_id) ON DELETE CASCADE,
	color VARCHAR(7) DEFAULT '#FFFFFF', 
	CHECK (category_code ~ '^TC-[0-9]+$' OR category_code IS NULL),	
	UNIQUE (account_id, transaction_category_id, category_code)
);
