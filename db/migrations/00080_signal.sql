-- +goose Up
CREATE TYPE SignalType AS ENUM (
	'flyover pool',
	'plan followup',
	'publicreport water',
	'publicreport nuisance',
	'residual exiring',
	'surveillance observation',
	'trap spike'
);
CREATE TYPE MosquitoSpecies AS ENUM (
	'none',
	'aedes aegypti',
	'aedes albopictus',
	'culex pipiens',
	'culex tarsalis'
);
CREATE TABLE signal (
	addressed TIMESTAMP WITHOUT TIME ZONE,
	addressor INTEGER REFERENCES user_(id),
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator INTEGER NOT NULL REFERENCES user_(id),
	id SERIAL NOT NULL,
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	species MosquitoSpecies,
	title TEXT NOT NULL,
	type_ SignalType NOT NULL,
	PRIMARY KEY(id)
);
CREATE TABLE signal_pool (
	pool_id INTEGER NOT NULL REFERENCES pool(id),
	signal_id INTEGER NOT NULL REFERENCES signal(id)
);
-- +goose Down
DROP TABLE signal_pool;
DROP TABLE signal;
DROP TYPE MosquitoSpecies;
DROP TYPE SignalType;
