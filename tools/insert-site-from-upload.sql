-- This query is just a one-off to quickly get uploaded data turned into sites.
INSERT INTO site (
    address_id,
    created,
    creator_id,
    notes,
    organization_id,
    owner_name,
    parcel_id,
    tags,
    version
)
SELECT DISTINCT ON (closest_address.id)
    closest_address.id AS address_id,
    fas.created,
    fas.creator_id,
    '' AS notes,
    fas.organization_id,
    '' AS owner_name,
    containing_parcel.id AS parcel_id,
    ''::hstore AS tags,
    1 AS version
FROM fileupload.flyover_aerial_service fas
CROSS JOIN LATERAL (
    SELECT a.id
    FROM address a
    ORDER BY a.geom <-> fas.geom
    LIMIT 1
) closest_address
CROSS JOIN LATERAL (
    SELECT p.id
    FROM parcel p
    WHERE ST_Contains(p.geometry, fas.geom)
    LIMIT 1
) containing_parcel
WHERE fas.geom IS NOT NULL
  AND fas.deleted IS NULL
ORDER BY closest_address.id, fas.created ASC  -- Keep the earliest created per address
ON CONFLICT (address_id) DO NOTHING;  -- Skip if address already has a site

