-- +goose UP

-- not text as we want exactly 64 characters
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose DOWN

ALTER TABLE feeds DROP COLUMN last_fetched_at;