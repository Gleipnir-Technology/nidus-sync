-- +goose Up
ALTER TABLE public.organization ADD COLUMN import_district_gid INTEGER UNIQUE REFERENCES import.district(gid);
ALTER TABLE public.organization ADD COLUMN website TEXT UNIQUE;
-- +goose Down
ALTER TABLE public.organization DROP COLUMN website;
ALTER TABLE public.organization DROP COLUMN import_district_gid;
