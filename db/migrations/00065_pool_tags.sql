-- +goose Up
ALTER TABLE fileupload.pool ADD COLUMN tags HSTORE NOT NULL;
-- +goose Down
ALTER TABLE fileupload.pool DROP COLUMN tags;
