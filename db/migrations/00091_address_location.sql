-- +goose Up
ALTER TABLE address RENAME COLUMN geom TO location;
-- +goose Down
ALTER TABLE address RENAME COLUMN location TO geom;

