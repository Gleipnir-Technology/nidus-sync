-- +goose Up
DROP TABLE tile.cached_image;
CREATE TABLE tile.cached_image (
	arcgis_id TEXT NOT NULL REFERENCES arcgis.service_map(arcgis_id),
	x INTEGER NOT NULL,
	y INTEGER NOT NULL,
	z INTEGER NOT NULL,
	is_empty BOOLEAN NOT NULL,
	PRIMARY KEY (arcgis_id, x, y, z)
);
-- +goose Down
DROP TABLE tile.cached_image;
CREATE TABLE tile.cached_image (
	arcgis_id TEXT NOT NULL REFERENCES arcgis.service_map(arcgis_id),
	x INTEGER NOT NULL,
	y INTEGER NOT NULL,
	z INTEGER NOT NULL,
	UNIQUE(arcgis_id, x, y, z)
);
