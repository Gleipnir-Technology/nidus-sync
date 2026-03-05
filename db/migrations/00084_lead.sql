-- +goose Up
CREATE TYPE LeadType AS ENUM (
	'green-pool'
);
CREATE TABLE lead (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator INTEGER NOT NULL REFERENCES user_(id),
	id SERIAL NOT NULL,
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	site_id INTEGER,
	site_version INTEGER,
	type_ LeadType NOT NULL,
	FOREIGN KEY (site_id, site_version) REFERENCES site(id, version),
	PRIMARY KEY (id)

);
ALTER TABLE compliance_report_request
	DROP CONSTRAINT compliance_report_request_site_id_site_version_fkey,
	DROP COLUMN site_id,
	DROP COLUMN site_version,
	ADD COLUMN lead_id INTEGER REFERENCES lead(id);
-- +goose Down
ALTER TABLE compliance_report_request
	DROP COLUMN lead_id,
	ADD COLUMN site_id INTEGER,
	ADD COLUMN site_version INTEGER,
	ADD FOREIGN KEY (site_id, site_version) REFERENCES site(id, version);
DROP TABLE lead;
DROP TYPE LeadType;
