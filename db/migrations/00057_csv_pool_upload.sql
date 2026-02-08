-- +goose Up
CREATE SCHEMA fileupload;
CREATE TYPE fileupload.FileStatusType AS ENUM (
	'uploaded',
	'parsed'
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
	file_id INTEGER REFERENCES fileupload.file(id) NOT NULL,
	type_ fileupload.CSVType NOT NULL,
	PRIMARY KEY (file_id)
);
CREATE TABLE fileupload.error (
	file_id INTEGER REFERENCES fileupload.file(id) NOT NULL,
	id SERIAL,
	line INTEGER NOT NULL,
	message TEXT NOT NULL,
	PRIMARY KEY (id)
);
-- +goose Down
DROP TABLE fileupload.error;
DROP TABLE fileupload.csv;
DROP TABLE fileupload.file;
DROP TYPE fileupload.CSVType;
DROP TYPE fileupload.FileStatusType;
DROP SCHEMA fileupload;
