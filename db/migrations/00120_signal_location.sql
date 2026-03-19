-- +goose Up
DELETE FROM signal;
ALTER TABLE signal ADD COLUMN location Geometry(Geometry, 4326) NOT NULL;
ALTER TABLE signal ADD COLUMN location_type TEXT GENERATED ALWAYS AS (GeometryType(location)) STORED;
ALTER TABLE signal ADD CONSTRAINT valid_location_types
		CHECK (location_type IN ('POINT', 'POLYGON', 'MULTIPOLYGON'));
CREATE INDEX idx_signal_location ON signal USING GIST(location);
CREATE INDEX idx_signal_location_type ON signal(location_type);
-- +goose Down
ALTER TABLE signal DROP INDEX idx_signal_location_type;
ALTER TABLE signal DROP INDEX idx_signal_location;
ALTER TABLE signal DROP COLUMN location_type;
ALTER TABLE signal DROP COLUMN location;
