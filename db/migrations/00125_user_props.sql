-- +goose Up
ALTER TABLE user_
	ADD COLUMN avatar UUID,
	ADD COLUMN is_active BOOLEAN,
	ADD COLUMN is_drone_pilot BOOLEAN,
	ADD COLUMN is_warrant BOOLEAN;
UPDATE user_ SET avatar=NULL, is_active = TRUE, is_drone_pilot=FALSE, is_warrant=FALSE;
ALTER TABLE user_ 
	ALTER COLUMN is_active SET NOT NULL,
	ALTER COLUMN is_drone_pilot SET NOT NULL,
	ALTER COLUMN is_warrant SET NOT NULL;
-- +goose Down
ALTER TABLE user_
	DROP COLUMN avatar,
	DROP COLUMN is_active,
	DROP COLUMN is_drone_pilot,
	DROP COLUMN is_warrant;
