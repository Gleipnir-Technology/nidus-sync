-- +goose Up
ALTER TABLE publicreport.pool ADD COLUMN has_backyard_permission BOOLEAN;
ALTER TABLE publicreport.pool ADD COLUMN is_reporter_confidential BOOLEAN;
ALTER TABLE publicreport.pool ADD COLUMN is_reporter_owner BOOLEAN;

UPDATE publicreport.pool SET has_backyard_permission=false;
UPDATE publicreport.pool SET is_reporter_confidential=false;
UPDATE publicreport.pool SET is_reporter_owner=false;

ALTER TABLE publicreport.pool ALTER COLUMN has_backyard_permission SET NOT NULL;
ALTER TABLE publicreport.pool ALTER COLUMN is_reporter_confidential SET NOT NULL;
ALTER TABLE publicreport.pool ALTER COLUMN is_reporter_owner SET NOT NULL;

ALTER TABLE publicreport.pool ALTER COLUMN map_zoom TYPE REAL;

ALTER TABLE publicreport.pool DROP COLUMN subscribe;
-- +goose Down
ALTER TABLE publicreport.pool ADD COLUMN subscribe BOOLEAN;

ALTER TABLE publicreport.pool ALTER COLUMN map_zoom TYPE FLOAT;

ALTER TABLE publicreport.pool DROP COLUMN has_backyard_permission;
ALTER TABLE publicreport.pool DROP COLUMN is_reporter_confidential;
ALTER TABLE publicreport.pool DROP COLUMN is_reporter_owner;
