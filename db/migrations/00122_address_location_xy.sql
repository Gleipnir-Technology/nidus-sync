-- +goose Up
ALTER TABLE address ADD COLUMN location_x DOUBLE PRECISION GENERATED ALWAYS AS (ST_X(location)) STORED;
ALTER TABLE address ADD COLUMN location_y DOUBLE PRECISION GENERATED ALWAYS AS (ST_Y(location)) STORED;
-- +goose Down
ALTER TABLE address DROP COLUMN location_y;
ALTER TABLE address DROP COLUMN location_x;
