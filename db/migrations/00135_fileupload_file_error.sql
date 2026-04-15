-- +goose Up
ALTER TABLE fileupload.file ADD COLUMN error TEXT;
UPDATE fileupload.file SET error = '';
ALTER TABLE fileupload.file ALTER COLUMN error SET NOT NULL;
-- +goose Down
ALTER TABLE fileupload.file DROP COLUMN error;
