-- +goose Up
CREATE TYPE publicreport.PermissionAccessType AS ENUM (
	'denied',
	'granted',
	'unselected',
	'with-owner'
);
ALTER TABLE publicreport.compliance
	ALTER COLUMN permission_type
	TYPE publicreport.PermissionAccessType USING permission_type::text::publicreport.PermissionAccessType;
DROP TYPE PermissionAccessType;
-- +goose Down
CREATE TYPE PermissionAccessType AS ENUM (
	'denied',
	'granted',
	'unselected',
	'with-owner'
);
ALTER TABLE publicreport.compliance
	ALTER COLUMN permission_type
	TYPE PermissionAccessType;
DROP TYPE publicreport.PermissionAccessType;
