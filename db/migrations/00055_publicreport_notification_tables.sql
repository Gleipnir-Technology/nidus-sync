-- +goose Up
CREATE TABLE publicreport.notify_phone_nuisance (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	nuisance_id INTEGER NOT NULL REFERENCES publicreport.nuisance(id),
	phone_e164 TEXT NOT NULL REFERENCES comms.phone(e164),
	PRIMARY KEY(nuisance_id, phone_e164)
);
CREATE TABLE publicreport.notify_phone_pool (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	phone_e164 TEXT NOT NULL REFERENCES comms.phone(e164),
	pool_id INTEGER NOT NULL REFERENCES publicreport.pool(id),
	PRIMARY KEY(pool_id, phone_e164)
);
CREATE TABLE publicreport.notify_email_nuisance (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	nuisance_id INTEGER NOT NULL REFERENCES publicreport.nuisance(id),
	email_address TEXT NOT NULL REFERENCES comms.email_contact(address),
	PRIMARY KEY(nuisance_id, email_address)
);
CREATE TABLE publicreport.notify_email_pool (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	pool_id INTEGER NOT NULL REFERENCES publicreport.pool(id),
	email_address TEXT NOT NULL REFERENCES comms.email_contact(address),
	PRIMARY KEY(pool_id, email_address)
);
-- +goose Down
DROP TABLE publicreport.notify_email_pool;
DROP TABLE publicreport.notify_email_nuisance;
DROP TABLE publicreport.notify_phone_pool;
DROP TABLE publicreport.notify_phone_nuisance;
