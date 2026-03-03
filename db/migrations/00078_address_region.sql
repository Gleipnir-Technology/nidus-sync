-- +goose Up
ALTER TABLE address ADD COLUMN region TEXT;
UPDATE address SET region = 'California';
ALTER TABLE address ALTER COLUMN region SET NOT NULL;
-- +goose Down
ALTER TABLE address DROP COLUMN region;
