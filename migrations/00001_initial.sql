-- +goose Up
CREATE TYPE arcgis_license_type AS ENUM (
	'advancedUT',
	'basicUT',
	'creatorUT',
	'editorUT',
	'fieldWorkerUT',
	'GISProfessionalAdvUT',
	'GISProfessionalBasicUT',
	'GISProfessionalStdUT',
	'IndoorsUserUT',
	'insightsAnalystUT',
	'liteUT',
	'standardUT',
	'storytellerUT',
	'viewerUT');

CREATE TABLE organization (
	id SERIAL PRIMARY KEY,
	name TEXT
);

CREATE TABLE user_ (
	id SERIAL PRIMARY KEY,
	arcgis_access_token TEXT,
	arcgis_license arcgis_license_type,
	arcgis_refresh_token TEXT,
	arcgis_refresh_token_expires TIMESTAMP,
	arcgis_role TEXT,
        display_name VARCHAR(200),
        email TEXT,
	organization_id INTEGER REFERENCES organization (id),
	username TEXT NOT NULL
);

-- +goose Down
DROP TABLE user_;
DROP TABLE organization;
DROP TYPE arcgis_license_type;
