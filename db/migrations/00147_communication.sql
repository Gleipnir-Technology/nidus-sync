-- +goose Up
CREATE TYPE CommunicationStatus AS ENUM (
	'closed',
	'invalid',
	'new',
	'opened',
	'pending',
	'possible-issue',
	'possible-resolved',
	'resolved'
);
CREATE TABLE communication (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL,
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	response_email_log_id INTEGER REFERENCES comms.email_log(id),
	response_text_log_id INTEGER REFERENCES comms.text_log(id),
	source_email_log_id INTEGER REFERENCES comms.email_log(id),
	source_report_id INTEGER REFERENCES publicreport.report(id),
	source_text_log_id INTEGER REFERENCES comms.text_log(id),
	status CommunicationStatus NOT NULL,
	PRIMARY KEY(id)
);
CREATE TYPE CommunicationLogEntry AS ENUM (
	'created',
	'status.closed',
	'status.invalidated',
	'status.opened',
	'status.pending',
	'status.possible-issue',
	'status.possible-resolved'
);
CREATE TABLE communication_log_entry (
	communication_id INTEGER NOT NULL REFERENCES communication(id),
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL,
	type_ CommunicationLogEntry NOT NULL,
	user_ INTEGER REFERENCES user_(id),
	PRIMARY KEY(id)
);
-- +goose Down
DROP TABLE communication;
DROP TABLE communication_log_entry;
DROP TYPE CommunicationLogEntry;
