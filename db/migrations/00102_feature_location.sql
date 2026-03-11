-- +goose Up
ALTER TABLE feature RENAME COLUMN geometry TO location;
-- +goose Down
ALTER TABLE feature RENAME COLUMN location TO geometry;
