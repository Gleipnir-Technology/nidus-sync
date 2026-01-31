-- +goose Up
CREATE TYPE comms.PhoneStatusType AS ENUM (
	'unconfirmed',
	'ok-to-send',
	'stopped'
);
ALTER TABLE comms.phone ADD COLUMN status comms.PhoneStatusType;
UPDATE comms.phone SET status = 'unconfirmed';
ALTER TABLE comms.phone ALTER COLUMN status SET NOT NULL;
UPDATE comms.phone SET is_subscribed = FALSE;
ALTER TABLE comms.phone ALTER COLUMN is_subscribed SET NOT NULL;
CREATE TABLE district_subscription_email (
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	email_contact_address TEXT NOT NULL REFERENCES comms.email_contact(address),
	PRIMARY KEY(organization_id, email_contact_address)
);
CREATE TABLE district_subscription_phone (
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	phone_e164 TEXT NOT NULL REFERENCES comms.phone(e164),
	PRIMARY KEY(organization_id, phone_e164)
);
