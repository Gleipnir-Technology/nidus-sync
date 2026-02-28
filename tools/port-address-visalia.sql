-- Use this to port over data that was imported from Visalia public parcels
-- in create-import-address-visalia.sql
INSERT INTO address(
 country,
 created     ,
 geom       ,
 h3cell ,
 locality,
 number_   ,
 postal_code,
 street,
 unit
) SELECT 
	'usa',
	NOW(),
	g.geom,
	h3_latlng_to_cell(g.geom, 15),
	a.municipality,
	TO_NUMBER(a.addrnum, '999999'),
	a.zipcode,
	a.fullname,
	COALESCE(a.unitid, '')
FROM import.addresses_visalia a
CROSS JOIN LATERAL public.geojsontogeom(a.geometry::jsonb) g;
