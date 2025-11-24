-- +goose Up
CREATE TYPE HashType AS ENUM (
	'bcrypt-14');

ALTER TABLE user_ ADD COLUMN password_hash_type HashType NOT NULL;
ALTER TABLE user_ ADD COLUMN password_hash TEXT NOT NULL;
-- +goose Down
ALTER TABLE user_ DROP COLUMN password_hash;
ALTER TABLE user_ DROP COLUMN password_hash_type;
DROP TYPE HashType;
