-- +goose Up
ALTER TABLE publicreport.report ADD COLUMN location_latitude DOUBLE PRECISION GENERATED ALWAYS AS (ST_Y(location)) STORED;
ALTER TABLE publicreport.report ADD COLUMN location_longitude DOUBLE PRECISION GENERATED ALWAYS AS (ST_X(location)) STORED;
-- +goose Down
ALTER TABLE publicreport.report DROP COLUMN location_longitude;
ALTER TABLE publicreport.report DROP COLUMN location_latitude;
