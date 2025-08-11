CREATE TABLE IF NOT EXISTS account_members (
	account_id UUID REFERENCES accounts(account_id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
  role VARCHAR(20) CHECK (role IN ('owner', 'editor', 'viewer')) NOT NULL DEFAULT 'editor',
  joined_at TIMESTAMPTZ NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),
	removed_at TIMESTAMPTZ DEFAULT NULL, 
  PRIMARY KEY (account_id, user_id)
);
