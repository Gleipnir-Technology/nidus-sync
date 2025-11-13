-- +goose Up
-- CREATE EXTENSION h3;
-- CREATE EXTENSION h3_postgis CASCADE;
CREATE TYPE H3AggregationType AS ENUM (
	'MosquitoSource',
	'ServiceRequest');

CREATE TABLE h3_aggregation (
	id SERIAL PRIMARY KEY,
	cell h3index NOT NULL,
	resolution INT NOT NULL,
	count_ INTEGER NOT NULL,
	type_ H3AggregationType NOT NULL,
	organization_id INTEGER REFERENCES organization (id) NOT NULL,
	UNIQUE(cell, organization_id, type_));

-- +goose Down
DROP TABLE h3_aggregation;
DROP TYPE H3AggregationType;
-- DROP EXTENSION h3;
