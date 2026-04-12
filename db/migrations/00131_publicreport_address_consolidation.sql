-- +goose Up
WITH new_addresses AS (
  INSERT INTO address (
    country,
    locality,
    postal_code,
    street,
    number_,
    region,
    gid,
    location,
    h3cell,
    created,
    unit
  )
  SELECT DISTINCT ON (r.address_gid)
    r.address_country,
    r.address_locality,
    r.address_postal_code,
    r.address_street,
    r.address_number,
    r.address_region,
    r.address_gid,
    r.location,
    r.h3cell,
    r.created,
    '' -- default empty string for unit since there's no corresponding column
  FROM publicreport.report r
  WHERE r.address_id IS NULL
    AND r.location IS NOT NULL
    AND r.h3cell IS NOT NULL
    AND NOT EXISTS (
      SELECT 1 FROM address a WHERE a.gid = r.address_gid
    )
  RETURNING id, gid
)
UPDATE publicreport.report r
SET address_id = a.id
FROM address a
WHERE r.address_gid = a.gid
  AND r.address_id IS NULL;
ALTER TABLE publicreport.report
	DROP COLUMN address_number,
	DROP COLUMN address_street,
	DROP COLUMN address_locality,
	DROP COLUMN address_region,
	DROP COLUMN address_postal_code,
	DROP COLUMN address_country;

