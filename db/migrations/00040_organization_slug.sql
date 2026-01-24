-- +goose Up
ALTER TABLE organization ADD COLUMN slug VARCHAR(24) UNIQUE;
-- +goose Down
ALTER TABLE organization DROP COLUMN slug;
