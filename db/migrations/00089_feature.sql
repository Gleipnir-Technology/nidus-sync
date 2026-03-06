-- +goose Up
CREATE TABLE feature (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER NOT NULL REFERENCES user_(id),
	id SERIAL NOT NULL,
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	site_id INTEGER NOT NULL,
	site_version INTEGER NOT NULL,
	geometry geometry(Point,4326),
	FOREIGN KEY (site_id, site_version) REFERENCES site(id, version),
	PRIMARY KEY(id)
);

CREATE TABLE feature_pool (
	feature_id INTEGER REFERENCES feature(id) NOT NULL,
	condition PoolConditionType NOT NULL,
	depth_meters FLOAT,
	geometry geometry(Polygon, 4326),
	PRIMARY KEY(feature_id)
);

DROP TABLE signal_pool;
DROP TABLE pool;

-- +goose Down

CREATE TABLE pool(
	id SERIAL,
	PRIMARY KEY(id)
);
CREATE TABLE signal_pool(
	id SERIAL,
	PRIMARY KEY(id)
);
DROP TABLE feature_pool;
DROP TABLE feature;
