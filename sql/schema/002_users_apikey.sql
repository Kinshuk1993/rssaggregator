-- +goose UP

-- not text as we want exactly 64 characters
ALTER TABLE users ADD COLUMN api_key VARCHAR(64) NOT NULL UNIQUE DEFAULT (
	encode(sha256(random()::text::bytea), 'hex')
);

-- +goose DOWN

ALTER TABLE users DROP COLUMN api_key;