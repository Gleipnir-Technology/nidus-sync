-- +goose Up
ALTER TABLE organization ALTER COLUMN name SET NOT NULL;

-- +goose Down
ALTER TABLE organization ALTER COLUMN name DROP NOT NULL;
