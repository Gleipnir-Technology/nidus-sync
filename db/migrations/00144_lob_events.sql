-- +goose Up
CREATE SCHEMA lob;
CREATE TABLE lob.event (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	body JSONB NOT NULL,
	id TEXT NOT NULL,
	type_ TEXT NOT NULL,
	PRIMARY KEY(id)
);
-- +goose Down
DROP TABLE lob.event;
DROP SCHEMA lob;
