-- +goose Up
CREATE TABLE fieldseeker_sync (
	id SERIAL PRIMARY KEY,
	created TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	records_created INTEGER NOT NULL,
	records_updated INTEGER NOT NULL,
	records_unchanged INTEGER NOT NULL,
	organization_id INTEGER REFERENCES organization(id) NOT NULL
);

-- +goose Down
DROP TABLE fieldseeker_sync;
