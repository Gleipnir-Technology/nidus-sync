-- +goose Up
CREATE TYPE publicreport.PermissionAccess AS ENUM (
	'denied',
	'granted',
	'unselected',
	'with-owner'
);
ALTER TABLE publicreport.compliance
	ALTER COLUMN permission_type
	TYPE publicreport.PermissionAccess USING permission_type::text::publicreport.PermissionAccess;
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
	TYPE PermissionAccessType USING permission_type::text::PermissionAccessType;
DROP TYPE publicreport.PermissionAccessType;
