-- TrapLocationBySourceID
SELECT 
	tl.globalid AS trap_location_globalid,
	ST_Distance(pl.geom, tl.geom) AS distance
FROM 
	fs_pointlocation pl
CROSS JOIN
	fs_traplocation tl
WHERE 
	tl.organization_id = $1 AND
	pl.globalid = $2
ORDER BY 
	ST_Distance(pl.geom, tl.geom)
LIMIT 4
