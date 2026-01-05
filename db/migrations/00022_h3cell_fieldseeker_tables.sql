-- +goose Up
ALTER TABLE fieldseeker.pointlocation ADD COLUMN h3cell h3index GENERATED ALWAYS AS (h3_latlng_to_cell(geospatial, 15)) STORED;
ALTER TABLE fieldseeker.servicerequest ADD COLUMN h3cell h3index GENERATED ALWAYS AS (h3_latlng_to_cell(geospatial, 15)) STORED;
ALTER TABLE fieldseeker.trapdata ADD COLUMN h3cell h3index GENERATED ALWAYS AS (h3_latlng_to_cell(geospatial, 15)) STORED;
ALTER TABLE fieldseeker.treatment ADD COLUMN h3cell h3index GENERATED ALWAYS AS (h3_latlng_to_cell(geospatial, 15)) STORED;

-- +goose Down
ALTER TABLE fieldseeker.pointlocation DROP COLUMN h3cell;
ALTER TABLE fieldseeker.servicerequest DROP COLUMN h3cell;
ALTER TABLE fieldseeker.trapdata DROP COLUMN h3cell;
ALTER TABLE fieldseeker.treatment DROP COLUMN h3cell;
