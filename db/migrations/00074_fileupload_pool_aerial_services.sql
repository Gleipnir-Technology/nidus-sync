-- +goose Up
ALTER TYPE fileupload.CSVType ADD VALUE 'Flyover' AFTER 'PoolList';
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
-- +goose Down
DROP TABLE fileupload.flyover_aerial_services;
ALTER TYPE fileupload.CSVType DROP VALUE 'Flyover';
