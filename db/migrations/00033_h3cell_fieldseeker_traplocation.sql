-- +goose Up
ALTER TABLE fieldseeker.traplocation ADD COLUMN h3cell h3index GENERATED ALWAYS AS (h3_latlng_to_cell(geospatial, 15)) STORED;

-- +goose Down
ALTER TABLE fieldseeker.traplocation DROP COLUMN h3cell;
