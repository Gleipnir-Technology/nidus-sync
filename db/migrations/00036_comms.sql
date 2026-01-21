-- +goose Up
CREATE SCHEMA comms;
CREATE TYPE comms.MessageTypeText AS ENUM (
	'initial-contact',
	'report-subscription-confirmation',
	'report-status-scheduled',
	'report-status-complete'
);
CREATE TYPE comms.MessageTypeEmail AS ENUM (
	'initial-contact',
	'report-subscription-confirmation',
	'report-status-scheduled',
	'report-status-complete'
);
CREATE TABLE comms.phone (
	e164 TEXT NOT NULL,
	is_subscribed BOOLEAN NOT NULL,
	PRIMARY KEY (e164)
);
CREATE TABLE comms.text_log (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	destination TEXT NOT NULL REFERENCES comms.phone(e164),
	source TEXT NOT NULL REFERENCES comms.phone(e164),
	type comms.MessageTypeText NOT NULL,
	PRIMARY KEY (destination, source, type)
);
CREATE TABLE comms.email (
	address TEXT NOT NULL,
	confirmed BOOLEAN NOT NULL,
	is_subscribed BOOLEAN NOT NULL,
	PRIMARY KEY(address)
);
CREATE TABLE comms.email_log (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	destination TEXT NOT NULL REFERENCES comms.email(address),
	source TEXT NOT NULL REFERENCES comms.phone(e164),
	type comms.MessageTypeEmail NOT NULL,
	PRIMARY KEY(destination, source, type)
);
-- +goose Down
DROP TABLE comms.email_log;
DROP TABLE comms.email;
DROP TABLE comms.text_log;
DROP TABLE comms.phone;
DROP TYPE comms.MessageTypeEmail;
DROP TYPE comms.MessageTypeText;
DROP SCHEMA comms;
