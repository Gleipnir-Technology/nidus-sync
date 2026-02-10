-- +goose Up
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
-- +goose Down
DROP TABLE publicreport.subscribe_phone;
DROP TABLE publicreport.subscribe_email;
