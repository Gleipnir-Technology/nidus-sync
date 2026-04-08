-- +goose Up
ALTER TABLE address ADD COLUMN gid TEXT;
UPDATE address SET gid = '';
ALTER TABLE address ALTER COLUMN gid SET NOT NULL;
ALTER TABLE publicreport.report ADD COLUMN address_gid TEXT;
UPDATE publicreport.report SET address_gid = '';
ALTER TABLE publicreport.report ALTER COLUMN address_gid SET NOT NULL;
-- +goose Down
ALTER TABLE publicreport.report DROP COLUMN address_gid;
ALTER TABLE address DROP COLUMN gid;
