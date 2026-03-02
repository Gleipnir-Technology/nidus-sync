-- +goose Up
DROP TABLE comms.mailer;
DROP TYPE comms.MailerType;
CREATE TABLE comms.mailer (
	address_id INTEGER NOT NULL REFERENCES address(id),
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL NOT NULL,
	recipient TEXT NOT NULL,
	uuid UUID NOT NULL,
	PRIMARY KEY(id)
);
CREATE TABLE compliance_report_request_mailer (
	compliance_report_request_id INTEGER NOT NULL REFERENCES compliance_report_request(id),
	mailer_id INTEGER NOT NULL REFERENCES comms.mailer(id),
	UNIQUE(compliance_report_request_id, mailer_id)
);
-- +goose Down
DROP TABLE compliance_report_request_mailer;
DROP TABLE comms.mailer;
CREATE TYPE comms.MailerType AS ENUM (
	'green-pool'
);
CREATE TABLE comms.mailer (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL NOT NULL,
	type_ comms.MailerType NOT NULL,
	PRIMARY KEY(id)
);

