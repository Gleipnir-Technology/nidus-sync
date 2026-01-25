-- +goose Up
DROP TABLE comms.text_log;
DROP TYPE comms.MessageTypeText;
CREATE TYPE comms.TextOrigin AS ENUM (
	'district',
	'llm',
	'website-action'
);
CREATE TABLE comms.text_log (
	content TEXT NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	destination TEXT NOT NULL REFERENCES comms.phone(e164),
	id SERIAL,
	origin comms.TextOrigin NOT NULL,
	source TEXT NOT NULL REFERENCES comms.phone(e164),
	PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE comms.text_log;
DROP TYPE comms.TextOrigin;
CREATE TYPE comms.MessageTypeText AS ENUM (
	'initial-contact',
	'report-subscription-confirmation',
	'report-status-scheduled',
	'report-status-complete'
);
CREATE TABLE comms.text_log (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	destination TEXT NOT NULL REFERENCES comms.phone(e164),
	source TEXT NOT NULL REFERENCES comms.phone(e164),
	type comms.MessageTypeText NOT NULL,
	PRIMARY KEY (destination, source, type)
);
