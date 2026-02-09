-- +goose Up
CREATE SCHEMA fileupload;
CREATE TYPE fileupload.FileStatusType AS ENUM (
	'error',
	'parsed',
	'uploaded'
);
CREATE TYPE fileupload.CSVType AS ENUM (
	'PoolList'
);
CREATE TABLE fileupload.file (
	id SERIAL,
	content_type TEXT NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER REFERENCES user_(id) NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	name TEXT NOT NULL,
	organization_id INTEGER REFERENCES organization(id) NOT NULL,
	status fileupload.FileStatusType NOT NULL,
	size_bytes INTEGER NOT NULL,
	file_uuid uuid NOT NULL,
	PRIMARY KEY(id)
);
CREATE TABLE fileupload.csv (
	committed TIMESTAMP WITHOUT TIME ZONE,
	file_id INTEGER REFERENCES fileupload.file(id) NOT NULL,
	rowcount INTEGER NOT NULL,
	type_ fileupload.CSVType NOT NULL,
	PRIMARY KEY (file_id)
);
CREATE TABLE fileupload.error_file (
	file_id INTEGER REFERENCES fileupload.file(id) NOT NULL,
	id SERIAL,
	message TEXT NOT NULL,
	PRIMARY KEY (id)
);
CREATE TABLE fileupload.error_csv (
	col INTEGER NOT NULL,
	csv_file_id INTEGER REFERENCES fileupload.csv(file_id) NOT NULL,
	id SERIAL,
	line INTEGER NOT NULL,
	message TEXT NOT NULL,
	PRIMARY KEY (id)
);
CREATE TYPE fileupload.PoolConditionType AS ENUM (
	'green',
	'murky',
	'blue',
	'unknown'
);
CREATE TABLE fileupload.pool (
	address_city TEXT NOT NULL,
	address_postal_code TEXT NOT NULL,
	address_street TEXT NOT NULL,
	committed BOOLEAN NOT NULL, -- Whether or not its just proposed before a CSV file is committed
	condition fileupload.PoolConditionType NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER REFERENCES user_(id) NOT NULL,
	csv_file INTEGER REFERENCES fileupload.csv(file_id) NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	geom geometry(Point, 3857),
	h3cell h3index,
	id SERIAL,
	is_in_district BOOLEAN NOT NULL, -- Whether or not the pool is within the district
	is_new BOOLEAN NOT NULL, -- Whether or not we already have a pool in the system for this row
	notes TEXT NOT NULL,
	organization_id INTEGER REFERENCES organization(id) NOT NULL,
	property_owner_name TEXT NOT NULL,
	property_owner_phone comms.phone,
	resident_owned BOOLEAN,
	resident_phone comms.phone,
	version integer,
	PRIMARY KEY (id, version)
);
-- +goose Down
DROP TABLE fileupload.pool;
DROP TYPE fileupload.PoolConditionType;
DROP TABLE fileupload.error_csv;
DROP TABLE fileupload.error_file;
DROP TABLE fileupload.csv;
DROP TABLE fileupload.file;
DROP TYPE fileupload.CSVType;
DROP TYPE fileupload.FileStatusType;
DROP SCHEMA fileupload;
