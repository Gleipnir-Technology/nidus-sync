-- +goose Up
ALTER TABLE organization ADD COLUMN lob_address_id TEXT;
-- +goose Down
ALTER TABLE organization DROP COLUMN lob_address_id;
