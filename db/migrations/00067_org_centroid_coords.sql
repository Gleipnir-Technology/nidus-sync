-- +goose Up
ALTER TABLE organization ADD COLUMN service_area_centroid_x DOUBLE PRECISION GENERATED ALWAYS AS (ST_X(ST_Centroid(service_area_geometry))) STORED;
ALTER TABLE organization ADD COLUMN service_area_centroid_y DOUBLE PRECISION GENERATED ALWAYS AS (ST_Y(ST_Centroid(service_area_geometry))) STORED;
