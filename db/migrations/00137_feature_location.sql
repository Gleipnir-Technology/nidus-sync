-- +goose Up
ALTER TABLE feature ADD COLUMN location_latitude DOUBLE PRECISION GENERATED ALWAYS AS (ST_Y(location)) STORED;
ALTER TABLE feature ADD COLUMN location_longitude DOUBLE PRECISION GENERATED ALWAYS AS (ST_X(location)) STORED;
-- +goose Down
ALTER TABLE feature DROP COLUMN location_longitude;
ALTER TABLE feature DROP COLUMN location_latitude;
