-- +goose Up
CREATE TYPE comms.MailerType AS ENUM (
	'green-pool'
);
CREATE TABLE comms.mailer (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL NOT NULL,
	type_ comms.MailerType NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE compliance_report_request (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator INTEGER REFERENCES user_(id) NOT NULL,
	id SERIAL NOT NULL,
	public_id TEXT NOT NULL UNIQUE,
	site_id INTEGER NOT NULL,
	site_version INTEGER NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY(site_id, site_version) REFERENCES site(id, version)
);
-- +goose Down
DROP TABLE compliance_report_request;
DROP TABLE comms.mailer;
DROP TYPE comms.MailerType;

