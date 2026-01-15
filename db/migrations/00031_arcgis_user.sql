-- +goose Up
CREATE SCHEMA arcgis;
CREATE TABLE arcgis.user_ (
	access TEXT NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	email TEXT NOT NULL,
	full_name TEXT NOT NULL,
	id TEXT NOT NULL,
	level TEXT NOT NULL,
	org_id TEXT NOT NULL,
	public_user_id INTEGER NOT NULL REFERENCES public.user_(id),
	region TEXT NOT NULL,
	role TEXT NOT NULL,
	role_id TEXT NOT NULL,
	username TEXT NOT NULL,
	user_license_type_id TEXT NOT NULL,
	user_type TEXT NOT NULL,
	PRIMARY KEY (id)
);
CREATE TABLE arcgis.user_privilege (
	user_id TEXT NOT NULL REFERENCES arcgis.user_(id),
	privilege TEXT NOT NULL,
	PRIMARY KEY(user_id, privilege)
);
-- +goose Down
DROP TABLE arcgis.user_privilege;
DROP TABLE arcgis.user_;
DROP SCHEMA arcgis;
