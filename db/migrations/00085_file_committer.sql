-- +goose Up
ALTER TABLE fileupload.file ADD COLUMN committer INTEGER REFERENCES user_(id);
-- +goose Down
ALTER TABLE fileupload.file DROP COLUMN committer;
