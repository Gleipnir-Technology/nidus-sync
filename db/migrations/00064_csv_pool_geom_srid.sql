-- +goose Up
ALTER TABLE fileupload.pool DROP COLUMN geom;
ALTER TABLE fileupload.pool ADD COLUMN geom geometry(Point, 4326);
-- +goose Down
ALTER TABLE fileupload.pool DROP COLUMN geom;
ALTER TABLE fileupload.pool ADD COLUMN geom geometry(Point, 3857);
