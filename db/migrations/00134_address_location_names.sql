-- +goose Up
ALTER TABLE address
	DROP COLUMN location_x,
	DROP COLUMN location_y;
ALTER TABLE address
	ADD COLUMN location_latitude DOUBLE PRECISION GENERATED ALWAYS AS (ST_Y(location)) STORED,
	ADD COLUMN location_longitude DOUBLE PRECISION GENERATED ALWAYS AS (ST_X(location)) STORED;
-- +goose Down
ALTER TABLE address
	DROP COLUMN location_latitude,
	DROP COLUMN location_longitude;
ALTER TABLE address
	ADD COLUMN location_y DOUBLE PRECISION GENERATED ALWAYS AS (ST_Y(location)) STORED,
	ADD COLUMN location_x DOUBLE PRECISION GENERATED ALWAYS AS (ST_X(location)) STORED;
