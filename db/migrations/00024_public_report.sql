-- +goose Up
CREATE SCHEMA IF NOT EXISTS publicreport;
CREATE TABLE publicreport.quick (
	id SERIAL PRIMARY KEY,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	comments TEXT NOT NULL,
	location GEOGRAPHY,
	h3cell h3index,
	uuid UUID NOT NULL
);

CREATE TABLE publicreport.quick_photo (
	id SERIAL PRIMARY KEY,
	size BIGINT NOT NULL,
	filename TEXT NOT NULL,
	quick_id INT NOT NULL REFERENCES publicreport.quick(id),
	uuid UUID NOT NULL
);

-- +goose Down
DROP TABLE publicreport.quick_photo;
DROP TABLE publicreport.quick;
DROP SCHEMA publicreport;

