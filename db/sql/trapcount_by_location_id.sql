-- TrapCountByLocationID
SELECT 
    td.globalid AS trapdata_globalid,
    td.enddatetime AS trapdata_enddate,
    COALESCE(SUM(sa.females), 0) AS total_females,
    COALESCE(SUM(sa.males), 0) AS total_males,
    COALESCE(SUM(sa.total), 0) AS total
FROM 
    fieldseeker.trapdata td
LEFT JOIN 
    fieldseeker.speciesabundance sa ON td.globalid = sa.trapdata_id
WHERE 
    td.organization_id = $1
    AND td.loc_id IN ($2)
GROUP BY 
    td.globalid, td.loc_id, td.enddatetime;

