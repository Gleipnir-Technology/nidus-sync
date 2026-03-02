-- +goose Up
CREATE INDEX idx_address_geom ON address USING GIST (geom);
CREATE INDEX idx_parcel_geometry ON parcel USING GIST (geometry);
