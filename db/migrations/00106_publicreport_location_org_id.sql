-- +goose Up
DROP VIEW publicreport.report_location;
CREATE VIEW publicreport.report_location AS 
SELECT
	ROW_NUMBER() OVER (ORDER BY table_name, public_id) AS id,
	table_name,
	address_id,
	address_raw,
	created,
	location,
	location_latitude,
	location_longitude,
	organization_id,
	public_id,
	status
FROM (
	SELECT
		'nuisance' AS table_name,
		address_id,
		address_raw,
		created,
		location,
		ST_X(location) AS location_longitude,
		ST_Y(location) AS location_latitude,
		organization_id,
		public_id,
		status
	FROM publicreport.nuisance
	UNION 
	SELECT
		'water' AS table_name,
		address_id,
		address_raw,
		created,
		location,
		ST_X(location) AS location_longitude,
		ST_Y(location) AS location_latitude,
		organization_id,
		public_id,
		status
	FROM publicreport.water
) AS combined_data;
-- +goose Down
DROP VIEW publicreport.report_location;
CREATE VIEW publicreport.report_location AS 
SELECT
	ROW_NUMBER() OVER (ORDER BY table_name, public_id) AS id,
	table_name,
	address_raw,
	created,
	location,
	public_id,
	status
FROM (
	SELECT
		'nuisance' AS table_name,
		address_raw,
		created,
		location,
		public_id,
		status
	FROM publicreport.nuisance
	UNION 
	SELECT
		'water' AS table_name,
		address_raw,
		created,
		location,
		public_id,
		status
	FROM publicreport.water
) AS combined_data;
