-- PublicreportPublicIDSuggestion
SELECT 
  'nuisance' AS table_name,
  public_id,
  location
FROM 
  publicreport.nuisance
WHERE 
  public_id LIKE $1

UNION ALL

SELECT 
  'water' AS table_name,
  public_id,
  location
FROM 
  publicreport.water
WHERE 
  public_id LIKE $1
ORDER BY
  public_id;
