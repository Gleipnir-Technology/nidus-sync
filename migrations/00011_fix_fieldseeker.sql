-- +goose Up
ALTER TABLE history_treatment ADD COLUMN created TIMESTAMP;

-- +goose Down
ALTER TABLE history_treatment DROP COLUMN created;
