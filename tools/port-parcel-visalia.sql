-- Use this to port over data that was imported from Visalia public parcels
-- in create-import-parcel-visalia.sql
INSERT INTO parcel(apn, description, geometry) 
SELECT 
    p.apn_id,
    p.propertysitus,
    -- g.geometrytype
    -- g.properties,
    g.geom
FROM import.csv_parcel p
CROSS JOIN LATERAL public.geojsontogeom(p.geometry::jsonb) g;
