-- +goose Up
ALTER TABLE pool ADD COLUMN geometry Geometry(Point, 4326);
-- +goose Down
ALTER TABLE pool DROP COLUMN geometry;
