-- +goose Up
ALTER TABLE publicreport.nuisance RENAME COLUMN location TO source_location;
CREATE TYPE publicreport.ReportStatusType AS ENUM (
	'reported',
	'reviewed',
	'scheduled',
	'treated'
);

ALTER TABLE publicreport.nuisance ADD COLUMN address TEXT;
UPDATE publicreport.nuisance SET address = '';
ALTER TABLE publicreport.nuisance ALTER COLUMN address SET NOT NULL;
ALTER TABLE publicreport.nuisance ADD COLUMN location GEOGRAPHY;
ALTER TABLE publicreport.nuisance ADD COLUMN status publicreport.ReportStatusType;
UPDATE publicreport.nuisance SET status = 'reported';
ALTER TABLE publicreport.nuisance ALTER COLUMN status SET NOT NULL;

ALTER TABLE publicreport.pool ADD COLUMN status publicreport.ReportStatusType;
UPDATE publicreport.pool SET status = 'reported';
ALTER TABLE publicreport.pool ALTER COLUMN status SET NOT NULL;

ALTER TABLE publicreport.quick ADD COLUMN address TEXT;
UPDATE publicreport.quick SET address = '';
ALTER TABLE publicreport.quick ALTER COLUMN address SET NOT NULL;
ALTER TABLE publicreport.quick ADD COLUMN status publicreport.ReportStatusType;
UPDATE publicreport.quick SET status = 'reported';
ALTER TABLE publicreport.quick ALTER COLUMN status SET NOT NULL;

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
	UNION 
	SELECT
		'quick' AS table_name,
		address,
		created,
		location,
		public_id,
		status
	FROM publicreport.quick
) AS combined_data;
