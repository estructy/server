ALTER TABLE users 
ADD COLUMN last_accessed_account UUID REFERENCES accounts(account_id) ON DELETE SET NULL;
