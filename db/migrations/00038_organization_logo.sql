-- +goose UP
ALTER TABLE public.organization ADD COLUMN logo_uuid UUID;
