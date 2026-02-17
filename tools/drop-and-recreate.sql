DROP DATABASE "nidus-sync";
-- ALTER DATABASE "nidus-sync" OWNER TO $1;
CREATE DATABASE "nidus-sync" WITH OWNER $1;
GRANT CONNECT ON DATABASE "nidus-sync" TO $1;
\c nidus-sync;
CREATE EXTENSION h3;
CREATE EXTENSION h3_postgis CASCADE;
CREATE EXTENSION hstore;
CREATE SCHEMA import;
ALTER SCHEMA import OWNER TO $1;
GRANT USAGE ON SCHEMA fileupload TO "tegola";
GRANT USAGE ON SCHEMA import TO "tegola";
GRANT USAGE ON SCHEMA publicreport TO "tegola";
GRANT SELECT ON fileupload.pool TO "tegola";
GRANT SELECT ON h3_aggregation to "tegola";
GRANT SELECT ON organization TO "tegola";
GRANT SELECT ON publicreport.report_location TO "tegola";
GRANT ALL PRIVILEGES ON SCHEMA public TO $1;
-- do import of district data
ALTER TABLE import.district ADD COLUMN geom_4326 geometry(MultiPolygon,4326) GENERATED ALWAYS AS (ST_Transform(geom, 4326)) STORED;
ALTER TABLE import.district ADD COLUMN centroid_4326 geometry(Point,4326) GENERATED ALWAYS AS (ST_Transform(ST_Centroid(geom), 4326)) STORED;
ALTER TABLE import.district ADD COLUMN extent_4326 geometry(Polygon,4326) GENERATED ALWAYS AS (ST_Transform(ST_Envelope(geom), 4326)) STORED;
ALTER TABLE import.district ADD COLUMN area_4326_sqm numeric GENERATED ALWAYS AS (ST_Area(ST_Transform(geom, 4326)::geography)) STORED;


UPDATE organization AS org
SET 
    website = dist.website,
    general_manager_name = dist.general_mg,
    mailing_address_city = dist.city2,
    mailing_address_postal_code = dist.postal_c_1::text,
    mailing_address_street = dist.address2,
    office_address_city = dist.city1,
    office_address_postal_code = dist.postal_cod::text,
    office_address_street = dist.address,
    office_phone = dist.phone1,
    office_fax = dist.fax1,
    service_area_geometry = dist.geom_4326
FROM import.district AS dist
WHERE org.import_district_gid = dist.gid;
