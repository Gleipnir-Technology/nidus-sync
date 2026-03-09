-- +goose Up
ALTER TABLE publicreport.image ADD COLUMN loc2 geometry(Point, 4326);
UPDATE publicreport.image SET loc2 = location::geometry(Point, 4326);
ALTER TABLE publicreport.image DROP COLUMN location;
ALTER TABLE publicreport.image RENAME COLUMN loc2 TO location;

