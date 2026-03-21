-- +goose up
ALTER TABLE signal DROP COLUMN title;
-- +goose down
ALTER TABLE signal ADD COLUMN title TEXT NOT NULL;
