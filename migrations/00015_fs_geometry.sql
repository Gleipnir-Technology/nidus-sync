-- +goose Up
ALTER TABLE fs_pointlocation ADD COLUMN geom geometry(Point, 3857);  -- as specified by the ArcGIS API
UPDATE fs_pointlocation SET geom = ST_SetSRID(ST_MakePoint(geometry_x, geometry_y), 3857);

ALTER TABLE fs_treatment ADD COLUMN geom geometry(Point, 3857);  -- as specified by the ArcGIS API
UPDATE fs_treatment SET geom = ST_SetSRID(ST_MakePoint(geometry_x, geometry_y), 3857);

ALTER TABLE fs_mosquitoinspection ADD COLUMN geom geometry(Point, 3857);  -- as specified by the ArcGIS API
UPDATE fs_mosquitoinspection SET geom = ST_SetSRID(ST_MakePoint(geometry_x, geometry_y), 3857);

-- +goose Down
ALTER TABLE fs_pointlocation DROP COLUMN geom;
ALTER TABLE fs_treatment DROP COLUMN geom;
ALTER TABLE fs_mosquitoinspection DROP COLUMN geom;
