-- +goose Up
ALTER TABLE fileupload.pool DROP COLUMN organization_id;
-- +goose Down
ALTER TABLE fileupload.pool ADD COLUMN organization_id INTEGER REFERENCES organization(id);
