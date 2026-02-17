-- +goose Up
ALTER TABLE organization ADD COLUMN general_manager_name TEXT;
ALTER TABLE organization ADD COLUMN mailing_address_city TEXT;
ALTER TABLE organization ADD COLUMN mailing_address_postal_code TEXT;
ALTER TABLE organization ADD COLUMN mailing_address_street TEXT;
ALTER TABLE organization ADD COLUMN office_address_city TEXT;
ALTER TABLE organization ADD COLUMN office_address_postal_code TEXT;
ALTER TABLE organization ADD COLUMN office_address_street TEXT;
ALTER TABLE organization ADD COLUMN office_fax TEXT;
ALTER TABLE organization ADD COLUMN office_phone TEXT;

ALTER TABLE organization ADD COLUMN service_area_geometry geometry(MultiPolygon,4326);

ALTER TABLE organization ADD COLUMN service_area_square_meters numeric GENERATED ALWAYS AS (ST_Area(service_area_geometry)) STORED;
ALTER TABLE organization ADD COLUMN service_area_centroid geometry(Point,4326) GENERATED ALWAYS AS (ST_Centroid(service_area_geometry)) STORED;
ALTER TABLE organization ADD COLUMN service_area_centroid_geojson TEXT GENERATED ALWAYS AS (ST_AsGeoJSON(ST_Centroid(service_area_geometry))) STORED;
ALTER TABLE organization ADD COLUMN service_area_extent geometry(Polygon,4326) GENERATED ALWAYS AS (ST_Envelope(service_area_geometry)) STORED;
ALTER TABLE organization ADD COLUMN service_area_xmin DOUBLE PRECISION GENERATED ALWAYS AS (ST_XMin(ST_Envelope(service_area_geometry))) STORED;
ALTER TABLE organization ADD COLUMN service_area_ymin DOUBLE PRECISION GENERATED ALWAYS AS (ST_YMin(ST_Envelope(service_area_geometry))) STORED;
ALTER TABLE organization ADD COLUMN service_area_xmax DOUBLE PRECISION GENERATED ALWAYS AS (ST_XMax(ST_Envelope(service_area_geometry))) STORED;
ALTER TABLE organization ADD COLUMN service_area_ymax DOUBLE PRECISION GENERATED ALWAYS AS (ST_YMax(ST_Envelope(service_area_geometry))) STORED;

ALTER TABLE organization DROP CONSTRAINT organization_website_key;
-- +goose Down
