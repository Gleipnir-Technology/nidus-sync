-- +goose Up
DROP TABLE comms.text_log;
DROP TYPE comms.MessageTypeText;
CREATE TYPE comms.TextOrigin AS ENUM (
	'district',
	'llm',
	'website-action'
);
CREATE TYPE comms.TextJobType AS ENUM (
	'report-confirmation'
);
CREATE TABLE comms.text_job (
	content TEXT NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	destination TEXT NOT NULL REFERENCES comms.phone(e164),
	id SERIAL NOT NULL,
	type_ comms.TextJobType NOT NULL,
	PRIMARY KEY(id)
);
COMMENT ON TABLE comms.text_job IS 'Used to track text messages that should be sent later';
CREATE TABLE comms.text_log (
	content TEXT NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	destination TEXT NOT NULL REFERENCES comms.phone(e164),
	id SERIAL,
	is_welcome BOOLEAN NOT NULL,
	origin comms.TextOrigin NOT NULL,
	source TEXT NOT NULL REFERENCES comms.phone(e164),
	PRIMARY KEY(id)
);
COMMENT ON TABLE comms.text_log IS 'Used to track text messages that were sent.';

-- +goose Down
DROP TABLE comms.text_log;
DROP TABLE comms.text_job;
DROP TYPE comms.TextJobType;
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
