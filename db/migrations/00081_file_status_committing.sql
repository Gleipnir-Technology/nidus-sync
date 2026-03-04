-- +goose Up
ALTER TYPE fileupload.FileStatusType ADD VALUE 'parsing' AFTER 'uploaded';
ALTER TYPE fileupload.FileStatusType ADD VALUE 'committing' AFTER 'parsing';
-- +goose Down
