-- +goose Up
DROP VIEW publicreport.report_location;

ALTER TABLE publicreport.nuisance ADD COLUMN loc2 geometry(Point, 4326);
UPDATE publicreport.nuisance SET loc2 = location::geometry(Point, 4326);
ALTER TABLE publicreport.nuisance DROP COLUMN location;
ALTER TABLE publicreport.nuisance RENAME COLUMN loc2 TO location;

ALTER TABLE publicreport.pool ADD COLUMN loc2 geometry(Point, 4326);
UPDATE publicreport.pool SET loc2 = location::geometry(Point, 4326);
ALTER TABLE publicreport.pool DROP COLUMN location;
ALTER TABLE publicreport.pool RENAME COLUMN loc2 TO location;

DROP TABLE publicreport.quick_image;
DROP TABLE publicreport.quick;

CREATE VIEW publicreport.report_location AS 
SELECT
	ROW_NUMBER() OVER (ORDER BY table_name, public_id) AS id,
	table_name,
	address,
	created,
	location,
	public_id,
	status
FROM (
	SELECT
		'nuisance' AS table_name,
		address,
		created,
		location,
		public_id,
		status
	FROM publicreport.nuisance
	UNION 
	SELECT
		'pool' AS table_name,
		address,
		created,
		location,
		public_id,
		status
	FROM publicreport.pool
) AS combined_data;
-- +goose Down
