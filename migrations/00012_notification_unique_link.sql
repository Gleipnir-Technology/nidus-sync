-- +goose Up
ALTER TABLE notification ADD COLUMN resolved_at TIMESTAMP WITHOUT TIME ZONE;
CREATE UNIQUE INDEX unique_user_link_not_resolved
ON notification (user_id, link)
WHERE resolved_at IS NULL;

-- +goose Down
DROP INDEX unique_user_link_not_resolved;
ALTER TABLE notification DROP COLUMN resolved_at;
