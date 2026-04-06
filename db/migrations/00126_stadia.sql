-- +goose Up
CREATE SCHEMA stadia;
CREATE TABLE stadia.api_request (
	id BIGSERIAL,
	request TEXT NOT NULL UNIQUE,  -- hash or identifier
	response JSONB NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
 	PRIMARY KEY(id)
);
-- +goose Down
DROP TABLE stadia.api_request;
DROP SCHEMA stadia;
