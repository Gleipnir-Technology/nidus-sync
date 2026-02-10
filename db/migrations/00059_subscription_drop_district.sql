-- +goose Up
DROP TABLE publicreport.subscribe_email;
DROP TABLE publicreport.subscribe_phone;

CREATE TABLE publicreport.subscribe_email (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	email_address TEXT NOT NULL REFERENCES comms.email_contact(address),
	id SERIAL,
	PRIMARY KEY(id)
);
CREATE TABLE publicreport.subscribe_phone (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	id SERIAL,
	phone_e164 TEXT NOT NULL REFERENCES comms.phone(e164),
	PRIMARY KEY(id)
);
-- +goose Down
DROP TABLE publicreport.subscribe_email;
DROP TABLE publicreport.subscribe_phone;

CREATE TABLE publicreport.subscribe_email (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	district_id INTEGER REFERENCES organization(id),
	email_address TEXT NOT NULL REFERENCES comms.email_contact(address),
	PRIMARY KEY(district_id, email_address)
);
CREATE TABLE publicreport.subscribe_phone (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	district_id INTEGER REFERENCES organization(id),
	phone_e164 TEXT NOT NULL REFERENCES comms.phone(e164),
	PRIMARY KEY(district_id, phone_e164)
);
