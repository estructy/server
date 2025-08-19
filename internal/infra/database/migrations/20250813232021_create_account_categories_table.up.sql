CREATE TABLE IF NOT EXISTS account_categories (
	account_category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	category_code VARCHAR(20),
	parent_id UUID REFERENCES account_categories(account_category_id) ON DELETE SET NULL,
	account_id UUID REFERENCES accounts(account_id) ON DELETE CASCADE,
	category_id UUID REFERENCES categories(category_id) ON DELETE CASCADE,
	color VARCHAR(7) DEFAULT '#FFFFFF', 
	CHECK (category_code ~ '^AC-[0-9]+$' OR category_code IS NULL),	
	UNIQUE (account_id, category_id, category_code)
);
