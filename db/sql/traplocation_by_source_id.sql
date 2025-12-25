-- TrapLocationBySourceID
SELECT 
	tl.globalid AS trap_location_globalid,
	ST_Distance(pl.geospatial, tl.geospatial) AS distance
FROM 
	fieldseeker.pointlocation pl
CROSS JOIN
	fieldseeker.traplocation tl
WHERE 
	tl.organization_id = $1 AND
	pl.globalid = $2
ORDER BY 
	ST_Distance(pl.geospatial, tl.geospatial)
LIMIT 4
