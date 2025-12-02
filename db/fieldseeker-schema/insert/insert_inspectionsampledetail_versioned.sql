-- Prepared statement for conditional insert with versioning for fieldseeker.inspectionsampledetail
-- Only inserts a new version if data has changed

PREPARE insert_inspectionsampledetail_versioned(bigint, uuid, fieldseeker.inspectionsampledetail_mosquitofieldspecies_enum, smallint, smallint, smallint, varchar, fieldseeker.inspectionsampledetail_mosquitodominantstage_enum, fieldseeker.inspectionsampledetail_mosquitoadultactivity_enum, varchar, smallint, smallint, smallint, fieldseeker.inspectionsampledetail_mosquitodominantstage_enum, varchar, uuid, varchar, timestamp, varchar, timestamp, smallint, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.inspectionsampledetail
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.inspectionsampledetail
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.inspectionsampledetail (
  objectid, inspsample_id, fieldspecies, flarvcount, fpupcount, feggcount, flstages, fdomstage, fadultact, labspecies, llarvcount, lpupcount, leggcount, ldomstage, comments, globalid, created_user, created_date, last_edited_user, last_edited_date, processed, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.inspsample_id IS NOT DISTINCT FROM $2 AND
    lv.fieldspecies IS NOT DISTINCT FROM $3 AND
    lv.flarvcount IS NOT DISTINCT FROM $4 AND
    lv.fpupcount IS NOT DISTINCT FROM $5 AND
    lv.feggcount IS NOT DISTINCT FROM $6 AND
    lv.flstages IS NOT DISTINCT FROM $7 AND
    lv.fdomstage IS NOT DISTINCT FROM $8 AND
    lv.fadultact IS NOT DISTINCT FROM $9 AND
    lv.labspecies IS NOT DISTINCT FROM $10 AND
    lv.llarvcount IS NOT DISTINCT FROM $11 AND
    lv.lpupcount IS NOT DISTINCT FROM $12 AND
    lv.leggcount IS NOT DISTINCT FROM $13 AND
    lv.ldomstage IS NOT DISTINCT FROM $14 AND
    lv.comments IS NOT DISTINCT FROM $15 AND
    lv.globalid IS NOT DISTINCT FROM $16 AND
    lv.created_user IS NOT DISTINCT FROM $17 AND
    lv.created_date IS NOT DISTINCT FROM $18 AND
    lv.last_edited_user IS NOT DISTINCT FROM $19 AND
    lv.last_edited_date IS NOT DISTINCT FROM $20 AND
    lv.processed IS NOT DISTINCT FROM $21 AND
    lv.creationdate IS NOT DISTINCT FROM $22 AND
    lv.creator IS NOT DISTINCT FROM $23 AND
    lv.editdate IS NOT DISTINCT FROM $24 AND
    lv.editor IS NOT DISTINCT FROM $25
  )
RETURNING *;

-- Example usage: EXECUTE insert_inspectionsampledetail_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: INSPSAMPLE_ID (uuid)
-- $3: FIELDSPECIES (fieldseeker.inspectionsampledetail_mosquitofieldspecies_enum)
-- $4: FLARVCOUNT (smallint)
-- $5: FPUPCOUNT (smallint)
-- $6: FEGGCOUNT (smallint)
-- $7: FLSTAGES (varchar)
-- $8: FDOMSTAGE (fieldseeker.inspectionsampledetail_mosquitodominantstage_enum)
-- $9: FADULTACT (fieldseeker.inspectionsampledetail_mosquitoadultactivity_enum)
-- $10: LABSPECIES (varchar)
-- $11: LLARVCOUNT (smallint)
-- $12: LPUPCOUNT (smallint)
-- $13: LEGGCOUNT (smallint)
-- $14: LDOMSTAGE (fieldseeker.inspectionsampledetail_mosquitodominantstage_enum)
-- $15: COMMENTS (varchar)
-- $16: GlobalID (uuid)
-- $17: created_user (varchar)
-- $18: created_date (timestamp)
-- $19: last_edited_user (varchar)
-- $20: last_edited_date (timestamp)
-- $21: PROCESSED (smallint)
-- $22: CreationDate (timestamp)
-- $23: Creator (varchar)
-- $24: EditDate (timestamp)
-- $25: Editor (varchar)
