-- +goose Up
INSERT INTO tile.service (name, arcgis_id) VALUES ('stadia', NULL);
ALTER TABLE tile.service
	ADD CONSTRAINT service_name_unique UNIQUE (name);
-- +goose Down
ALTER TABLE tile.service DROP CONSTRAINT service_name_unique;
DELETE FROM tile.service WHERE name = 'stadia';
