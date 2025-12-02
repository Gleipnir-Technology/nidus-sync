-- Prepared statement for conditional insert with versioning for fieldseeker.inspectionsample
-- Only inserts a new version if data has changed

PREPARE insert_inspectionsample_versioned(bigint, uuid, varchar, fieldseeker.inspectionsample_notinuit_f_enum, varchar, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.inspectionsample
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.inspectionsample
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.inspectionsample (
  objectid, insp_id, sampleid, processed, idbytech, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.insp_id IS NOT DISTINCT FROM $2 AND
    lv.sampleid IS NOT DISTINCT FROM $3 AND
    lv.processed IS NOT DISTINCT FROM $4 AND
    lv.idbytech IS NOT DISTINCT FROM $5 AND
    lv.globalid IS NOT DISTINCT FROM $6 AND
    lv.created_user IS NOT DISTINCT FROM $7 AND
    lv.created_date IS NOT DISTINCT FROM $8 AND
    lv.last_edited_user IS NOT DISTINCT FROM $9 AND
    lv.last_edited_date IS NOT DISTINCT FROM $10 AND
    lv.creationdate IS NOT DISTINCT FROM $11 AND
    lv.creator IS NOT DISTINCT FROM $12 AND
    lv.editdate IS NOT DISTINCT FROM $13 AND
    lv.editor IS NOT DISTINCT FROM $14
  )
RETURNING *;

-- Example usage: EXECUTE insert_inspectionsample_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: INSP_ID (uuid)
-- $3: SAMPLEID (varchar)
-- $4: PROCESSED (fieldseeker.inspectionsample_notinuit_f_enum)
-- $5: IDBYTECH (varchar)
-- $6: GlobalID (uuid)
-- $7: created_user (varchar)
-- $8: created_date (timestamp)
-- $9: last_edited_user (varchar)
-- $10: last_edited_date (timestamp)
-- $11: CreationDate (timestamp)
-- $12: Creator (varchar)
-- $13: EditDate (timestamp)
-- $14: Editor (varchar)
