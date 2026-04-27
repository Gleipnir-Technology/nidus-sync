-- +goose Up
ALTER TABLE publicreport.compliance ADD COLUMN submitted TIMESTAMP WITHOUT TIME ZONE;
-- +goose Down
ALTER TABLE publicreport.compliance DROP COLUMN submitted;
