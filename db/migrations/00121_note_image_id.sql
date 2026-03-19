-- +goose Up
ALTER TABLE note_image
ADD COLUMN id INTEGER GENERATED ALWAYS AS IDENTITY;
ALTER TABLE note_image
ADD CONSTRAINT note_image_id_unique UNIQUE (id);
-- +goose Down
ALTER TABLE note_image DROP COLUMN id;
