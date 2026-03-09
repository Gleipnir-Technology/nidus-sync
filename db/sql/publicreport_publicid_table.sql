-- PublicreportIDTable
WITH found_tables AS (
	SELECT 
		'nuisance' as table_name,
		id
	FROM publicreport.nuisance
	WHERE public_id = $1
    
	UNION ALL
    
	SELECT
		'water'  as table_name,
		id
	FROM publicreport.water
	WHERE public_id = $1
)
SELECT 
	EXISTS (SELECT 1 FROM found_tables) as exists_somewhere,
	array_agg(table_name) as found_in_tables,
	array_agg(id) as report_ids
FROM found_tables;
