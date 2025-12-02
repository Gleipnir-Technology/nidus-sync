-- Prepared statement for conditional insert with versioning for fieldseeker.treatmentarea
-- Only inserts a new version if data has changed

PREPARE insert_treatmentarea_versioned(bigint, uuid, uuid, timestamp, varchar, uuid, varchar, timestamp, varchar, timestamp, smallint, varchar, timestamp, varchar, timestamp, varchar, double precision, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.treatmentarea
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.treatmentarea
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.treatmentarea (
  objectid, treat_id, session_id, treatdate, comments, globalid, created_user, created_date, last_edited_user, last_edited_date, notified, type, creationdate, creator, editdate, editor, shape__area, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.treat_id IS NOT DISTINCT FROM $2 AND
    lv.session_id IS NOT DISTINCT FROM $3 AND
    lv.treatdate IS NOT DISTINCT FROM $4 AND
    lv.comments IS NOT DISTINCT FROM $5 AND
    lv.globalid IS NOT DISTINCT FROM $6 AND
    lv.created_user IS NOT DISTINCT FROM $7 AND
    lv.created_date IS NOT DISTINCT FROM $8 AND
    lv.last_edited_user IS NOT DISTINCT FROM $9 AND
    lv.last_edited_date IS NOT DISTINCT FROM $10 AND
    lv.notified IS NOT DISTINCT FROM $11 AND
    lv.type IS NOT DISTINCT FROM $12 AND
    lv.creationdate IS NOT DISTINCT FROM $13 AND
    lv.creator IS NOT DISTINCT FROM $14 AND
    lv.editdate IS NOT DISTINCT FROM $15 AND
    lv.editor IS NOT DISTINCT FROM $16 AND
    lv.shape__area IS NOT DISTINCT FROM $17 AND
    lv.shape__length IS NOT DISTINCT FROM $18
  )
RETURNING *;

-- Example usage: EXECUTE insert_treatmentarea_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: TREAT_ID (uuid)
-- $3: SESSION_ID (uuid)
-- $4: TREATDATE (timestamp)
-- $5: COMMENTS (varchar)
-- $6: GlobalID (uuid)
-- $7: created_user (varchar)
-- $8: created_date (timestamp)
-- $9: last_edited_user (varchar)
-- $10: last_edited_date (timestamp)
-- $11: Notified (smallint)
-- $12: Type (varchar)
-- $13: CreationDate (timestamp)
-- $14: Creator (varchar)
-- $15: EditDate (timestamp)
-- $16: Editor (varchar)
-- $17: Shape__Area (double precision)
-- $18: Shape__Length (double precision)
