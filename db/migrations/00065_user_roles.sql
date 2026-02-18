-- +goose Up
CREATE TYPE UserRole AS ENUM (
	'root',
	'account-owner'
);
ALTER TABLE user_ ADD COLUMN role UserRole;
UPDATE user_ SET role = 'account-owner';
ALTER TABLE user_ ALTER COLUMN role SET NOT NULL;
-- +goose Down
ALTER TABLE user_ DROP COLUMN role;
DROP TYPE UserRole;
