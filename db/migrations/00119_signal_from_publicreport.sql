-- +goose Up
ALTER TABLE signal ADD COLUMN site_id INTEGER REFERENCES site(id);
-- +goose Down
ALTER TABLE signal DROP COLUMN site_id;
