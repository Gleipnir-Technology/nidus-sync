-- +goose Up
ALTER TABLE organization ADD COLUMN arcgis_id TEXT;
ALTER TABLE organization ADD COLUMN arcgis_name TEXT;
ALTER TABLE oauth_token ADD COLUMN arcgis_id TEXT;
ALTER TABLE oauth_token ADD COLUMN arcgis_license_type_id TEXT;
-- +goose Down
ALTER TABLE organization DROP COLUMN arcgis_id;
ALTER TABLE organization DROP COLUMN arcgis_name;
ALTER TABLE oauth_token DROP COLUMN arcgis_id;
ALTER TABLE oauth_token DROP COLUMN arcgis_license_type_id;
