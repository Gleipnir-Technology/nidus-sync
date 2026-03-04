-- +goose Up
ALTER TABLE address DROP COLUMN number_;
ALTER TABLE address ADD COLUMN number_ TEXT;
UPDATE address SET number_ = '';
ALTER TABLE address ALTER COLUMN number_ SET NOT NULL;
-- +goose Down
ALTER TABLE address DROP COLUMN number_;
ALTER TABLE address ADD COLUMN number_ INTEGER;
UPDATE address SET number_ = 0;
ALTER TABLE address ALTER COLUMN number_ SET NOT NULL;
