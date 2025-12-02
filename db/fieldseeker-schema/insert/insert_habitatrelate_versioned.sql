-- Prepared statement for conditional insert with versioning for fieldseeker.habitatrelate
-- Only inserts a new version if data has changed

PREPARE insert_habitatrelate_versioned(bigint, uuid, uuid, varchar, timestamp, varchar, timestamp, fieldseeker.habitatrelate_habitatrelate_habitattype_2e81cf2f550e400783cf284f3cec3953_enum, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.habitatrelate
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.habitatrelate
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.habitatrelate (
  objectid, foreign_id, globalid, created_user, created_date, last_edited_user, last_edited_date, habitattype, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.foreign_id IS NOT DISTINCT FROM $2 AND
    lv.globalid IS NOT DISTINCT FROM $3 AND
    lv.created_user IS NOT DISTINCT FROM $4 AND
    lv.created_date IS NOT DISTINCT FROM $5 AND
    lv.last_edited_user IS NOT DISTINCT FROM $6 AND
    lv.last_edited_date IS NOT DISTINCT FROM $7 AND
    lv.habitattype IS NOT DISTINCT FROM $8 AND
    lv.creationdate IS NOT DISTINCT FROM $9 AND
    lv.creator IS NOT DISTINCT FROM $10 AND
    lv.editdate IS NOT DISTINCT FROM $11 AND
    lv.editor IS NOT DISTINCT FROM $12
  )
RETURNING *;

-- Example usage: EXECUTE insert_habitatrelate_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: FOREIGN_ID (uuid)
-- $3: GlobalID (uuid)
-- $4: created_user (varchar)
-- $5: created_date (timestamp)
-- $6: last_edited_user (varchar)
-- $7: last_edited_date (timestamp)
-- $8: HABITATTYPE (fieldseeker.habitatrelate_habitatrelate_habitattype_2e81cf2f550e400783cf284f3cec3953_enum)
-- $9: CreationDate (timestamp)
-- $10: Creator (varchar)
-- $11: EditDate (timestamp)
-- $12: Editor (varchar)
