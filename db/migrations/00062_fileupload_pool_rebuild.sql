-- +goose Up
DROP TABLE fileupload.pool;

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
	geom geometry(Point, 4326),
	h3cell h3index,
	id SERIAL,
	is_in_district BOOLEAN NOT NULL, -- Whether or not the pool is within the district
	is_new BOOLEAN NOT NULL, -- Whether or not we already have a pool in the system for this row
	notes TEXT NOT NULL,
	organization_id INTEGER REFERENCES organization(id) NOT NULL,
	property_owner_name TEXT NOT NULL,
	property_owner_phone_e164 TEXT REFERENCES comms.phone(e164),
	resident_owned BOOLEAN,
	resident_phone_e164 TEXT REFERENCES comms.phone(e164),
	line_number INTEGER NOT NULL,
	tags HSTORE NOT NULL,
	PRIMARY KEY (id)
);

-- migration 62
-- ALTER TABLE fileupload.pool ADD COLUMN property_owner_phone_e164 TEXT REFERENCES comms.phone(e164);
-- ALTER TABLE fileupload.pool ADD COLUMN resident_phone_e164 TEXT REFERENCES comms.phone(e164);

-- migration 64
-- ALTER TABLE fileupload.pool DROP COLUMN geom;
-- ALTER TABLE fileupload.pool ADD COLUMN geom geometry(Point, 4326);

-- migration 65
-- ALTER TABLE fileupload.pool ADD COLUMN tags HSTORE NOT NULL;

-- migration 66
-- ALTER TABLE fileupload.pool ADD COLUMN row_number INTEGER NOT NULL;

-- +goose Down
