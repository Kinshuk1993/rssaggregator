-- +goose UP

CREATE TABLE IF NOT EXISTS feed_follows (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
	-- Ensure a user cannot follow the same feed more than once
	UNIQUE(user_id, feed_id)
);

-- +goose DOWN

DROP TABLE IF EXISTS feed_follows;