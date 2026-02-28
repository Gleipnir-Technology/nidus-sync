-- +goose Up
CREATE TABLE arcgis.account (
	id TEXT NOT NULL,
	name TEXT NOT NULL,
	organization_id INTEGER NOT NULL REFERENCES organization(id),

	url_features TEXT,
	url_insights TEXT,
	url_geometry TEXT,
	url_notebooks TEXT,
	url_tiles TEXT,
	PRIMARY KEY(id)
);
CREATE TABLE arcgis.service_map (
	account_id TEXT NOT NULL REFERENCES arcgis.account(id),
	arcgis_id TEXT NOT NULL,
	name TEXT NOT NULL,
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	PRIMARY KEY(arcgis_id)
);
ALTER TABLE arcgis.feature_service RENAME TO service_feature;
ALTER TABLE arcgis.service_feature ADD COLUMN account_id TEXT REFERENCES arcgis.account(id);
CREATE TABLE arcgis.oauth_token (
	access_token TEXT NOT NULL,
	access_token_expires TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	arcgis_account_id TEXT REFERENCES arcgis.account(id),
	arcgis_id TEXT,
	arcgis_license_type_id TEXT,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL NOT NULL,
	invalidated_at TIMESTAMP WITHOUT TIME ZONE,
	refresh_token TEXT NOT NULL,
	refresh_token_expires TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	user_id INTEGER NOT NULL REFERENCES user_(id),
	username TEXT NOT NULL,
	PRIMARY KEY (id)
);


ALTER TABLE organization ADD COLUMN arcgis_account_id TEXT REFERENCES arcgis.account(id);
ALTER TABLE organization ADD COLUMN fieldseeker_service_feature_item_id TEXT REFERENCES arcgis.service_feature(item_id);
ALTER TABLE organization DROP COLUMN arcgis_id;
ALTER TABLE organization DROP COLUMN arcgis_name;
ALTER TABLE organization DROP COLUMN fieldseeker_url;
DROP TABLE oauth_token;
-- +goose Down
CREATE TABLE oauth_token (
	access_token TEXT NOT NULL,
	access_token_expires TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	-- arcgis_account_id INTEGER REFERENCES arcgis.account(id),
	arcgis_id TEXT,
	arcgis_license_type_id TEXT,
	created TIMESTAMP WITHOUT TIME ZONE,
	id SERIAL NOT NULL,
	invalidated_at TIMESTAMP WITHOUT TIME ZONE,
	refresh_token TEXT NOT NULL,
	refresh_token_expires TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	user_id INTEGER NOT NULL REFERENCES user_(id),
	username TEXT NOT NULL,
	PRIMARY KEY (id)
);
ALTER TABLE organization ADD COLUMN fieldseeker_url TEXT;
ALTER TABLE organization ADD COLUMN arcgis_name TEXT;
ALTER TABLE organization ADD COLUMN arcgis_id TEXT;
ALTER TABLE organization DROP COLUMN fieldseeker_service_feature_item_id;
ALTER TABLE organization DROP COLUMN arcgis_account_id;
DROP TABLE arcgis.oauth_token;
ALTER TABLE arcgis.service_feature DROP COLUMN account_id;
ALTER TABLE arcgis.service_feature RENAME TO feature_service;
DROP TABLE arcgis.service_map;
DROP TABLE arcgis.account;
