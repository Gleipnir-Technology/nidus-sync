-- +goose Up
CREATE TYPE CountryType AS ENUM (
	'usa'
);
CREATE TABLE address (
	country CountryType NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	geom geometry(Point,4326) NOT NULL,
	h3cell h3index NOT NULL,
	id SERIAL NOT NULL,
	locality TEXT NOT NULL,
	number_ INTEGER NOT NULL,
	postal_code TEXT NOT NULL,
	street TEXT NOT NULL,
	unit TEXT NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(country, locality, number_, street)
);

CREATE TABLE site (
	address_id INTEGER REFERENCES address(id) NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER REFERENCES user_(id) NOT NULL,
	file_id INTEGER REFERENCES fileupload.file(id),
	id SERIAL NOT NULL,
	notes text NOT NULL,
	organization_id INTEGER NOT NULL,
	owner_name TEXT NOT NULL,
	owner_phone_e164 TEXT,
	resident_owned BOOLEAN,
	resident_phone_e164 TEXT,
	tags HSTORE NOT NULL,
	version INTEGER NOT NULL,
	PRIMARY KEY(id, version),
	UNIQUE(address_id)
);

CREATE TYPE PoolConditionType AS ENUM (
	'blue',
	'dry',
	'false pool',
	'green',
	'murky'
);
CREATE TABLE pool (
	condition PoolConditionType NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER REFERENCES user_(id) NOT NULL,
	id SERIAL NOT NULL,
	site_id INTEGER,
	PRIMARY KEY(id)
);
