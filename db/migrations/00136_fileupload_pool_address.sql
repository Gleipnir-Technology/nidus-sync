-- +goose Up
ALTER TABLE fileupload.pool ADD COLUMN address_id INTEGER REFERENCES address(id);
-- +goose Down
ALTER TABLE fileupload.pool DROP COLUMN address_id;
