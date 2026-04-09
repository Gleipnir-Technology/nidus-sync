-- +goose Up
CREATE TYPE PermissionAccessType AS ENUM (
	'denied',
	'granted',
	'unselected',
	'with-owner'
);
CREATE TABLE publicreport.compliance (
	access_instructions TEXT NOT NULL,
	availability_notes TEXT NOT NULL,
	comments TEXT NOT NULL,
	gate_code TEXT NOT NULL,
	has_dog BOOLEAN,
	permission_type PermissionAccessType NOT NULL,
	report_id INTEGER REFERENCES publicreport.report(id),
	report_phone_can_text BOOLEAN,
	wants_scheduled BOOLEAN,
	PRIMARY KEY(report_id)
);
CREATE TABLE publicreport.client (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	user_agent TEXT NOT NULL,
	uuid UUID NOT NULL,
	PRIMARY KEY(uuid)
);
ALTER TABLE publicreport.report ADD COLUMN client_uuid UUID REFERENCES publicreport.client(uuid);
-- +goose Down
ALTER TABLE publicreport.report DROP COLUMN client_uuid;
DROP TABLE publicreport.client;
DROP TABLE publicreport.compliance;
DROP TYPE PermissionAccessType;
