-- +goose Up
ALTER TABLE user_
	ADD COLUMN is_active BOOLEAN,
	ADD COLUMN is_drone_pilot BOOLEAN,
	ADD COLUMN is_warrant BOOLEAN,
	ADD COLUMN avatar UUID;
UPDATE user_ SET is_active = TRUE;
ALTER TABLE user_ ALTER COLUMN is_active SET NOT NULL;
-- +goose Down
ALTER TABLE user_
	DROP COLUMN is_active,
	DROP COLUMN is_drone_pilot,
	DROP COLUMN is_warrant,
	DROP COLUMN avatar;
