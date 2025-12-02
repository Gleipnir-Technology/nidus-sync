-- Prepared statement for conditional insert with versioning for fieldseeker.zones2
-- Only inserts a new version if data has changed

PREPARE insert_zones2_versioned(bigint, varchar, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar, double precision, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.zones2
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.zones2
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.zones2 (
  objectid, name, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, shape__area, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.name IS NOT DISTINCT FROM $2 AND
    lv.globalid IS NOT DISTINCT FROM $3 AND
    lv.created_user IS NOT DISTINCT FROM $4 AND
    lv.created_date IS NOT DISTINCT FROM $5 AND
    lv.last_edited_user IS NOT DISTINCT FROM $6 AND
    lv.last_edited_date IS NOT DISTINCT FROM $7 AND
    lv.creationdate IS NOT DISTINCT FROM $8 AND
    lv.creator IS NOT DISTINCT FROM $9 AND
    lv.editdate IS NOT DISTINCT FROM $10 AND
    lv.editor IS NOT DISTINCT FROM $11 AND
    lv.shape__area IS NOT DISTINCT FROM $12 AND
    lv.shape__length IS NOT DISTINCT FROM $13
  )
RETURNING *;

-- Example usage: EXECUTE insert_zones2_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: NAME (varchar)
-- $3: GlobalID (uuid)
-- $4: created_user (varchar)
-- $5: created_date (timestamp)
-- $6: last_edited_user (varchar)
-- $7: last_edited_date (timestamp)
-- $8: CreationDate (timestamp)
-- $9: Creator (varchar)
-- $10: EditDate (timestamp)
-- $11: Editor (varchar)
-- $12: Shape__Area (double precision)
-- $13: Shape__Length (double precision)
