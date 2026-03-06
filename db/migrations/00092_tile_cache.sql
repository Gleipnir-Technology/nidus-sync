-- +goose Up
CREATE SCHEMA tile;
CREATE TABLE tile.cached_image (
	arcgis_id TEXT NOT NULL REFERENCES arcgis.service_map(arcgis_id),
	x INTEGER NOT NULL,
	y INTEGER NOT NULL,
	z INTEGER NOT NULL,
	UNIQUE(arcgis_id, x, y, z)
);
ALTER TABLE organization ADD COLUMN arcgis_map_service_id TEXT REFERENCES arcgis.service_map(arcgis_id);
-- +goose Down
ALTER TABLE organization DROP COLUMN arcgis_map_service_id;
DROP TABLE tile.cached_image;
DROP SCHEMA tile;
