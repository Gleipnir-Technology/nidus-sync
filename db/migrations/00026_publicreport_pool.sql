-- +goose Up
CREATE TYPE publicreport.PoolSourceDuration AS ENUM (
	'none',
	'less-than-week',
	'1-2-weeks',
	'2-4-weeks',
	'1-3-months',
	'more-than-3-months'
);
CREATE TABLE publicreport.pool (
	id SERIAL PRIMARY KEY,
	access_comments TEXT NOT NULL,
	access_gate BOOLEAN NOT NULL,
	access_fence BOOLEAN NOT NULL,
	access_locked BOOLEAN NOT NULL,
	access_dog BOOLEAN NOT NULL,
	access_other BOOLEAN NOT NULL,
	address TEXT NOT NULL,
	address_country TEXT NOT NULL,
	address_post_code TEXT NOT NULL,
	address_place TEXT NOT NULL,
	address_street TEXT NOT NULL,
	address_region TEXT NOT NULL,
	comments TEXT NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	h3cell h3index,
	has_adult BOOLEAN NOT NULL,
	has_larvae BOOLEAN NOT NULL,
	has_pupae BOOLEAN NOT NULL,
	location GEOGRAPHY,
	map_zoom FLOAT NOT NULL,
	owner_email TEXT NOT NULL,
	owner_name TEXT NOT NULL,
	owner_phone TEXT NOT NULL,
	public_id TEXT NOT NULL UNIQUE,
	reporter_email TEXT NOT NULL,
	reporter_name TEXT NOT NULL,
	reporter_phone TEXT NOT NULL,
	subscribe BOOLEAN NOT NULL
);

CREATE TABLE publicreport.pool_photo (
	id SERIAL PRIMARY KEY,
	size BIGINT NOT NULL,
	filename TEXT NOT NULL,
	pool_id INT NOT NULL REFERENCES publicreport.pool(id),
	uuid UUID NOT NULL
);

-- +goose Down
DROP TABLE publicreport.pool_photo;
DROP TABLE publicreport.pool;
DROP TYPE publicreport.PoolSourceDuration;
