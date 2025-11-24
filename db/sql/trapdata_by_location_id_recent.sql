-- TrapDataByLocationIDRecent
SELECT enddatetime, globalid, loc_id
FROM (
	SELECT enddatetime, globalid, loc_id, ROW_NUMBER()
	OVER (PARTITION BY loc_id ORDER BY enddatetime DESC) as row_num
	FROM fs_trapdata
	WHERE 
		organization_id = $1 AND
		loc_id IN ($2)
) ranked_data
WHERE row_num <= 10
ORDER BY enddatetime DESC;
