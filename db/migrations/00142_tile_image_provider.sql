-- +goose Up
CREATE TABLE tile.service (
	id SERIAL,
	name TEXT NOT NULL,
	arcgis_id TEXT REFERENCES arcgis.service_map(arcgis_id),
	PRIMARY KEY(id)
);

INSERT INTO tile.service (name, arcgis_id)
	SELECT name, arcgis_id
	FROM arcgis.service_map;

ALTER TABLE tile.cached_image ADD COLUMN service_id INTEGER REFERENCES tile.service(id);

UPDATE tile.cached_image
SET service_id = tile.service.id
FROM tile.service
WHERE tile.service.arcgis_id = tile.cached_image.arcgis_id;

ALTER TABLE tile.cached_image 
	DROP CONSTRAINT cached_image_pkey,
	ALTER COLUMN arcgis_id DROP NOT NULL,
	ALTER COLUMN service_id SET NOT NULL,
	ADD PRIMARY KEY (service_id, x, y, z),
	DROP COLUMN arcgis_id;
