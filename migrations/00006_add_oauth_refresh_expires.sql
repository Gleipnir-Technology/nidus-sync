-- +goose Up
ALTER TABLE oauth_token RENAME COLUMN expires TO access_token_expires;
ALTER TABLE oauth_token ADD COLUMN refresh_token_expires TIMESTAMP NOT NULL DEFAULT current_timestamp;

-- +goose Down
ALTER TABLE oauth_token DROP COLUMN refresh_token_expires;
ALTER TABLE oauth_token RENAME COLUMN access_token_expires TO expires;
