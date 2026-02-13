-- +goose Up
ALTER TABLE oauth_token ADD COLUMN created TIMESTAMP WITHOUT TIME ZONE;
UPDATE oauth_token SET created = now();
ALTER TABLE oauth_token ALTER COLUMN created SET NOT NULL;
-- +goose Down
ALTER TABLE oauth_token DROP COLUMN created;
