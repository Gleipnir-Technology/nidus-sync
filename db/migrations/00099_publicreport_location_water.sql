-- +goose Up
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
		'pool' AS table_name,
		address_raw,
		created,
		location,
		public_id,
		status
	FROM publicreport.water
) AS combined_data;
