-- +goose Up
-- CREATE EXTENSION IF NOT EXISTS hstore;
ALTER TABLE comms.email RENAME TO email_contact;
ALTER TABLE comms.email_contact ADD COLUMN public_id TEXT;
UPDATE comms.email_contact SET public_id = '';
ALTER TABLE comms.email_contact ALTER COLUMN public_id SET NOT NULL;
CREATE TABLE comms.email_template (
	content_html TEXT NOT NULL,
	content_txt TEXT NOT NULL,
	content_hash_html VARCHAR(64) NOT NULL,
	content_hash_txt VARCHAR(64) NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL NOT NULL,
	superceded TIMESTAMP WITHOUT TIME ZONE,
	message_type comms.MessageTypeEmail NOT NULL,
	PRIMARY KEY (id)
);
DROP TABLE comms.email_log;
CREATE TABLE comms.email_log (
	id SERIAL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	delivery_status VARCHAR(16) NOT NULL,
	destination TEXT NOT NULL REFERENCES comms.email_contact(address),
	public_id VARCHAR(64) NOT NULL,
	sent_at TIMESTAMP WITHOUT TIME ZONE,
	source TEXT NOT NULL,
	subject VARCHAR(255) NOT NULL,
	template_id INTEGER REFERENCES comms.email_template(id),
	template_data HSTORE NOT NULL,
	type comms.MessageTypeEmail NOT NULL,
	PRIMARY KEY (id)
);
-- +goose Down
DROP TABLE comms.email_log;
CREATE TABLE comms.email_log (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	destination TEXT NOT NULL REFERENCES comms.email_contact(address),
	source TEXT NOT NULL REFERENCES comms.phone(e164),
	type comms.MessageTypeEmail NOT NULL,
	PRIMARY KEY(destination, source, type)
);
DROP TABLE comms.email_template;
ALTER TABLE comms.email_contact DROP COLUMN public_id;
ALTER TABLE comms.email_contact RENAME TO email;
