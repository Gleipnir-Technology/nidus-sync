-- Prepared statement for conditional insert with versioning for fieldseeker.rodentlocation
-- Only inserts a new version if data has changed

PREPARE insert_rodentlocation_versioned(bigint, varchar, varchar, varchar, fieldseeker.rodentlocation_rodentlocationhabitat_enum, fieldseeker.rodentlocation_locationpriority_1_enum, fieldseeker.rodentlocation_locationusetype_1_enum, fieldseeker.rodentlocation_notinuit_f_1_enum, varchar, varchar, varchar, fieldseeker.rodentlocation_rodentlocation_symbology_enum, varchar, timestamp, integer, timestamp, varchar, varchar, varchar, varchar, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.rodentlocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.rodentlocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.rodentlocation (
  objectid, locationname, zone, zone2, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, nextactiondatescheduled, locationnumber, lastinspectdate, lastinspectspecies, lastinspectaction, lastinspectconditions, lastinspectrodentevidence, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, jurisdiction,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.locationname IS NOT DISTINCT FROM $2 AND
    lv.zone IS NOT DISTINCT FROM $3 AND
    lv.zone2 IS NOT DISTINCT FROM $4 AND
    lv.habitat IS NOT DISTINCT FROM $5 AND
    lv.priority IS NOT DISTINCT FROM $6 AND
    lv.usetype IS NOT DISTINCT FROM $7 AND
    lv.active IS NOT DISTINCT FROM $8 AND
    lv.description IS NOT DISTINCT FROM $9 AND
    lv.accessdesc IS NOT DISTINCT FROM $10 AND
    lv.comments IS NOT DISTINCT FROM $11 AND
    lv.symbology IS NOT DISTINCT FROM $12 AND
    lv.externalid IS NOT DISTINCT FROM $13 AND
    lv.nextactiondatescheduled IS NOT DISTINCT FROM $14 AND
    lv.locationnumber IS NOT DISTINCT FROM $15 AND
    lv.lastinspectdate IS NOT DISTINCT FROM $16 AND
    lv.lastinspectspecies IS NOT DISTINCT FROM $17 AND
    lv.lastinspectaction IS NOT DISTINCT FROM $18 AND
    lv.lastinspectconditions IS NOT DISTINCT FROM $19 AND
    lv.lastinspectrodentevidence IS NOT DISTINCT FROM $20 AND
    lv.globalid IS NOT DISTINCT FROM $21 AND
    lv.created_user IS NOT DISTINCT FROM $22 AND
    lv.created_date IS NOT DISTINCT FROM $23 AND
    lv.last_edited_user IS NOT DISTINCT FROM $24 AND
    lv.last_edited_date IS NOT DISTINCT FROM $25 AND
    lv.creationdate IS NOT DISTINCT FROM $26 AND
    lv.creator IS NOT DISTINCT FROM $27 AND
    lv.editdate IS NOT DISTINCT FROM $28 AND
    lv.editor IS NOT DISTINCT FROM $29 AND
    lv.jurisdiction IS NOT DISTINCT FROM $30
  )
RETURNING *;

-- Example usage: EXECUTE insert_rodentlocation_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: LOCATIONNAME (varchar)
-- $3: ZONE (varchar)
-- $4: ZONE2 (varchar)
-- $5: HABITAT (fieldseeker.rodentlocation_rodentlocationhabitat_enum)
-- $6: PRIORITY (fieldseeker.rodentlocation_locationpriority_1_enum)
-- $7: USETYPE (fieldseeker.rodentlocation_locationusetype_1_enum)
-- $8: ACTIVE (fieldseeker.rodentlocation_notinuit_f_1_enum)
-- $9: DESCRIPTION (varchar)
-- $10: ACCESSDESC (varchar)
-- $11: COMMENTS (varchar)
-- $12: SYMBOLOGY (fieldseeker.rodentlocation_rodentlocation_symbology_enum)
-- $13: EXTERNALID (varchar)
-- $14: NEXTACTIONDATESCHEDULED (timestamp)
-- $15: LOCATIONNUMBER (integer)
-- $16: LASTINSPECTDATE (timestamp)
-- $17: LASTINSPECTSPECIES (varchar)
-- $18: LASTINSPECTACTION (varchar)
-- $19: LASTINSPECTCONDITIONS (varchar)
-- $20: LASTINSPECTRODENTEVIDENCE (varchar)
-- $21: GlobalID (uuid)
-- $22: created_user (varchar)
-- $23: created_date (timestamp)
-- $24: last_edited_user (varchar)
-- $25: last_edited_date (timestamp)
-- $26: CreationDate (timestamp)
-- $27: Creator (varchar)
-- $28: EditDate (timestamp)
-- $29: Editor (varchar)
-- $30: JURISDICTION (varchar)
