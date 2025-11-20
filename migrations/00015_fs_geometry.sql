-- +goose Up
ALTER TABLE fs_mosquitoinspection ADD COLUMN geom geometry(Point, 3857);  -- as specified by the ArcGIS API
UPDATE fs_mosquitoinspection SET geom = ST_SetSRID(ST_MakePoint(geometry_x, geometry_y), 3857);
CREATE INDEX idx_fs_mosquitoinspection_geom ON fs_mosquitoinspection USING GIST(geom);

ALTER TABLE fs_pointlocation ADD COLUMN geom geometry(Point, 3857);  -- as specified by the ArcGIS API
UPDATE fs_pointlocation SET geom = ST_SetSRID(ST_MakePoint(geometry_x, geometry_y), 3857);
CREATE INDEX idx_fs_pointlocation_geom ON fs_pointlocation USING GIST(geom);


--ALTER TABLE fs_trapdata ADD COLUMN geom geometry(Point, 3857);  -- as specified by the ArcGIS API
--UPDATE fs_trapdata SET geom = ST_SetSRID(ST_MakePoint(geometry_x, geometry_y), 3857);

ALTER TABLE fs_traplocation ADD COLUMN geom geometry(Point, 3857);  -- as specified by the ArcGIS API
UPDATE fs_traplocation SET geom = ST_SetSRID(ST_MakePoint(geometry_x, geometry_y), 3857);
CREATE INDEX idx_fs_traplocation_geom ON fs_traplocation USING GIST(geom);

ALTER TABLE fs_treatment ADD COLUMN geom geometry(Point, 3857);  -- as specified by the ArcGIS API
UPDATE fs_treatment SET geom = ST_SetSRID(ST_MakePoint(geometry_x, geometry_y), 3857);
CREATE INDEX idx_fs_treatment_geom ON fs_treatment USING GIST(geom);

-- +goose Down
DROP  INDEX idx_fs_mosquitoinspection_geom;
ALTER TABLE fs_mosquitoinspection DROP COLUMN geom;
DROP  INDEX idx_fs_pointlocation_geom;
ALTER TABLE fs_pointlocation DROP COLUMN geom;
--ALTER TABLE fs_trapdata DROP COLUMN geom;
DROP  INDEX idx_fs_traplocation_geom;
ALTER TABLE fs_traplocation DROP COLUMN geom;
DROP  INDEX idx_fs_treatment_geom;
ALTER TABLE fs_treatment DROP COLUMN geom;
