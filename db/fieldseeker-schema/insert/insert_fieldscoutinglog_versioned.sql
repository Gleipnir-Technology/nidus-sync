-- Prepared statement for conditional insert with versioning for fieldseeker.fieldscoutinglog
-- Only inserts a new version if data has changed

PREPARE insert_fieldscoutinglog_versioned(bigint, fieldseeker.fieldscoutinglog_fieldscoutingsymbology_enum, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.fieldscoutinglog
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.fieldscoutinglog
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.fieldscoutinglog (
  objectid, status, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.status IS NOT DISTINCT FROM $2 AND
    lv.globalid IS NOT DISTINCT FROM $3 AND
    lv.created_user IS NOT DISTINCT FROM $4 AND
    lv.created_date IS NOT DISTINCT FROM $5 AND
    lv.last_edited_user IS NOT DISTINCT FROM $6 AND
    lv.last_edited_date IS NOT DISTINCT FROM $7 AND
    lv.creationdate IS NOT DISTINCT FROM $8 AND
    lv.creator IS NOT DISTINCT FROM $9 AND
    lv.editdate IS NOT DISTINCT FROM $10 AND
    lv.editor IS NOT DISTINCT FROM $11
  )
RETURNING *;

-- Example usage: EXECUTE insert_fieldscoutinglog_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: STATUS (fieldseeker.fieldscoutinglog_fieldscoutingsymbology_enum)
-- $3: GlobalID (uuid)
-- $4: created_user (varchar)
-- $5: created_date (timestamp)
-- $6: last_edited_user (varchar)
-- $7: last_edited_date (timestamp)
-- $8: CreationDate (timestamp)
-- $9: Creator (varchar)
-- $10: EditDate (timestamp)
-- $11: Editor (varchar)
