-- +goose Up
ALTER TABLE oauth_token ADD COLUMN invalidated_at timestamp without time zone;

-- +goose Down
ALTER TABLE oauth_token DROP COLUMN invalidated_at;
