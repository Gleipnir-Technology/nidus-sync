-- Prepared statement for conditional insert with versioning for fieldseeker.stormdrain
-- Only inserts a new version if data has changed

PREPARE insert_stormdrain_versioned(bigint, timestamp, timestamp, varchar, fieldseeker.stormdrain_stormdrainsymbology_enum, uuid, varchar, timestamp, varchar, timestamp, varchar, varchar, varchar, timestamp, varchar, timestamp, varchar, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.stormdrain
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.stormdrain
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.stormdrain (
  objectid, nexttreatmentdate, lasttreatdate, lastaction, symbology, globalid, created_user, created_date, last_edited_user, last_edited_date, laststatus, zone, zone2, creationdate, creator, editdate, editor, type, jurisdiction,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.nexttreatmentdate IS NOT DISTINCT FROM $2 AND
    lv.lasttreatdate IS NOT DISTINCT FROM $3 AND
    lv.lastaction IS NOT DISTINCT FROM $4 AND
    lv.symbology IS NOT DISTINCT FROM $5 AND
    lv.globalid IS NOT DISTINCT FROM $6 AND
    lv.created_user IS NOT DISTINCT FROM $7 AND
    lv.created_date IS NOT DISTINCT FROM $8 AND
    lv.last_edited_user IS NOT DISTINCT FROM $9 AND
    lv.last_edited_date IS NOT DISTINCT FROM $10 AND
    lv.laststatus IS NOT DISTINCT FROM $11 AND
    lv.zone IS NOT DISTINCT FROM $12 AND
    lv.zone2 IS NOT DISTINCT FROM $13 AND
    lv.creationdate IS NOT DISTINCT FROM $14 AND
    lv.creator IS NOT DISTINCT FROM $15 AND
    lv.editdate IS NOT DISTINCT FROM $16 AND
    lv.editor IS NOT DISTINCT FROM $17 AND
    lv.type IS NOT DISTINCT FROM $18 AND
    lv.jurisdiction IS NOT DISTINCT FROM $19
  )
RETURNING *;

-- Example usage: EXECUTE insert_stormdrain_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: NextTreatmentDate (timestamp)
-- $3: LastTreatDate (timestamp)
-- $4: LastAction (varchar)
-- $5: Symbology (fieldseeker.stormdrain_stormdrainsymbology_enum)
-- $6: GlobalID (uuid)
-- $7: created_user (varchar)
-- $8: created_date (timestamp)
-- $9: last_edited_user (varchar)
-- $10: last_edited_date (timestamp)
-- $11: LastStatus (varchar)
-- $12: ZONE (varchar)
-- $13: ZONE2 (varchar)
-- $14: CreationDate (timestamp)
-- $15: Creator (varchar)
-- $16: EditDate (timestamp)
-- $17: Editor (varchar)
-- $18: TYPE (varchar)
-- $19: JURISDICTION (varchar)
