CREATE TABLE IF NOT EXISTS users (
	user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	password_hash TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC')
);
