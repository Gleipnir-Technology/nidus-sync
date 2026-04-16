-- +goose Up
ALTER TABLE comms.mailer ADD COLUMN external_id TEXT NOT NULL;
-- +goose Down
ALTER TABLE comms.mailer DROP COLUMN external_id;
