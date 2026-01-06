-- +goose Up
ALTER TABLE user_ ALTER COLUMN organization_id SET NOT NULL;

-- +goose Down
ALTER TABLE user_ ALTER COLUMN organization_id DROP NOT NULL;
