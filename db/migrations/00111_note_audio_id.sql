-- +goose Up
ALTER TABLE note_audio
ADD COLUMN id INTEGER GENERATED ALWAYS AS IDENTITY;
ALTER TABLE note_audio
ADD CONSTRAINT note_audio_id_unique UNIQUE (id);
-- +goose Down
ALTER TABLE note_audio DROP COLUMN id;
