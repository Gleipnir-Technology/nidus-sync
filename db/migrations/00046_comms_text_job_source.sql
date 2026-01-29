-- +goose Up
CREATE TYPE comms.TextJobSource AS ENUM (
	'rmo',
	'nidus'
);

ALTER TABLE comms.text_job ADD COLUMN source comms.TextJobSource;
UPDATE comms.text_job SET source = 'rmo';
ALTER TABLE comms.text_job ALTER COLUMN source SET NOT NULL;
ALTER TABLE comms.text_job ADD COLUMN completed TIMESTAMP WITHOUT TIME ZONE;

