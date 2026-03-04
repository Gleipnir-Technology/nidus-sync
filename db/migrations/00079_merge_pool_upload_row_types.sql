-- +goose Up
ALTER TABLE fileupload.pool ADD COLUMN address_number TEXT;
UPDATE fileupload.pool SET address_number = '';
ALTER TABLE fileupload.pool ALTER COLUMN address_number SET NOT NULL;

ALTER TABLE fileupload.pool ADD COLUMN address_locality TEXT;
UPDATE fileupload.pool SET address_locality = address_city;
ALTER TABLE fileupload.pool ALTER COLUMN address_locality SET NOT NULL;

ALTER TABLE fileupload.pool ADD COLUMN address_region TEXT;
UPDATE fileupload.pool SET address_region = '';
ALTER TABLE fileupload.pool ALTER COLUMN address_region SET NOT NULL;

ALTER TABLE fileupload.pool DROP COLUMN address_city;

DROP TABLE fileupload.flyover_aerial_service;
-- +goose Down
CREATE TABLE fileupload.flyover_aerial_service (
	committed BOOLEAN NOT NULL, -- Whether or not its just proposed before a CSV file is committed
	condition fileupload.PoolConditionType NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER REFERENCES user_(id) NOT NULL,
	csv_file INTEGER REFERENCES fileupload.csv(file_id) NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	geom geometry(Point, 4326),
	h3cell h3index,
	id SERIAL,
	organization_id INTEGER REFERENCES organization(id) NOT NULL,
	PRIMARY KEY (id)
);

ALTER TABLE fileupload.pool ADD COLUMN address_city TEXT;
UPDATE fileupload.pool SET address_city = '';
ALTER TABLE fileupload.pool ALTER COLUMN address_city SET NOT NULL;

ALTER TABLE fileupload.pool DROP COLUMN address_region;
ALTER TABLE fileupload.pool DROP COLUMN address_locality;
ALTER TABLE fileupload.pool DROP COLUMN address_number;
