-- Prepared statement for conditional insert with versioning for fieldseeker.linelocation
-- Only inserts a new version if data has changed

PREPARE insert_linelocation_versioned(bigint, varchar, varchar, fieldseeker.linelocation_linelocation_habitat_fc51bdc4f1954df58206d69ce14182f3_enum, fieldseeker.linelocation_locationpriority_enum, fieldseeker.linelocation_linelocation_usetype_2aeca2e60d2f455c86fc34895dc80a02_enum, fieldseeker.linelocation_notinuit_f_enum, varchar, varchar, varchar, fieldseeker.linelocation_locationsymbology_enum, varchar, double precision, timestamp, smallint, double precision, double precision, varchar, integer, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, double precision, double precision, varchar, varchar, varchar, timestamp, varchar, double precision, varchar, double precision, varchar, varchar, double precision, double precision, varchar, fieldseeker.linelocation_linelocation_waterorigin_84723d92_306a_46f4_8ef1_69b55a916008_enum, timestamp, varchar, timestamp, varchar, varchar, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.linelocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.linelocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.linelocation (
  objectid, name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, acres, nextactiondatescheduled, larvinspectinterval, length_ft, width_ft, zone2, locationnumber, globalid, created_user, created_date, last_edited_user, last_edited_date, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, hectares, lastinspectactivity, lasttreatactivity, length_meters, width_meters, lastinspectconditions, waterorigin, creationdate, creator, editdate, editor, jurisdiction, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.name IS NOT DISTINCT FROM $2 AND
    lv.zone IS NOT DISTINCT FROM $3 AND
    lv.habitat IS NOT DISTINCT FROM $4 AND
    lv.priority IS NOT DISTINCT FROM $5 AND
    lv.usetype IS NOT DISTINCT FROM $6 AND
    lv.active IS NOT DISTINCT FROM $7 AND
    lv.description IS NOT DISTINCT FROM $8 AND
    lv.accessdesc IS NOT DISTINCT FROM $9 AND
    lv.comments IS NOT DISTINCT FROM $10 AND
    lv.symbology IS NOT DISTINCT FROM $11 AND
    lv.externalid IS NOT DISTINCT FROM $12 AND
    lv.acres IS NOT DISTINCT FROM $13 AND
    lv.nextactiondatescheduled IS NOT DISTINCT FROM $14 AND
    lv.larvinspectinterval IS NOT DISTINCT FROM $15 AND
    lv.length_ft IS NOT DISTINCT FROM $16 AND
    lv.width_ft IS NOT DISTINCT FROM $17 AND
    lv.zone2 IS NOT DISTINCT FROM $18 AND
    lv.locationnumber IS NOT DISTINCT FROM $19 AND
    lv.globalid IS NOT DISTINCT FROM $20 AND
    lv.created_user IS NOT DISTINCT FROM $21 AND
    lv.created_date IS NOT DISTINCT FROM $22 AND
    lv.last_edited_user IS NOT DISTINCT FROM $23 AND
    lv.last_edited_date IS NOT DISTINCT FROM $24 AND
    lv.lastinspectdate IS NOT DISTINCT FROM $25 AND
    lv.lastinspectbreeding IS NOT DISTINCT FROM $26 AND
    lv.lastinspectavglarvae IS NOT DISTINCT FROM $27 AND
    lv.lastinspectavgpupae IS NOT DISTINCT FROM $28 AND
    lv.lastinspectlstages IS NOT DISTINCT FROM $29 AND
    lv.lastinspectactiontaken IS NOT DISTINCT FROM $30 AND
    lv.lastinspectfieldspecies IS NOT DISTINCT FROM $31 AND
    lv.lasttreatdate IS NOT DISTINCT FROM $32 AND
    lv.lasttreatproduct IS NOT DISTINCT FROM $33 AND
    lv.lasttreatqty IS NOT DISTINCT FROM $34 AND
    lv.lasttreatqtyunit IS NOT DISTINCT FROM $35 AND
    lv.hectares IS NOT DISTINCT FROM $36 AND
    lv.lastinspectactivity IS NOT DISTINCT FROM $37 AND
    lv.lasttreatactivity IS NOT DISTINCT FROM $38 AND
    lv.length_meters IS NOT DISTINCT FROM $39 AND
    lv.width_meters IS NOT DISTINCT FROM $40 AND
    lv.lastinspectconditions IS NOT DISTINCT FROM $41 AND
    lv.waterorigin IS NOT DISTINCT FROM $42 AND
    lv.creationdate IS NOT DISTINCT FROM $43 AND
    lv.creator IS NOT DISTINCT FROM $44 AND
    lv.editdate IS NOT DISTINCT FROM $45 AND
    lv.editor IS NOT DISTINCT FROM $46 AND
    lv.jurisdiction IS NOT DISTINCT FROM $47 AND
    lv.shape__length IS NOT DISTINCT FROM $48
  )
RETURNING *;

-- Example usage: EXECUTE insert_linelocation_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: NAME (varchar)
-- $3: ZONE (varchar)
-- $4: HABITAT (fieldseeker.linelocation_linelocation_habitat_fc51bdc4f1954df58206d69ce14182f3_enum)
-- $5: PRIORITY (fieldseeker.linelocation_locationpriority_enum)
-- $6: USETYPE (fieldseeker.linelocation_linelocation_usetype_2aeca2e60d2f455c86fc34895dc80a02_enum)
-- $7: ACTIVE (fieldseeker.linelocation_notinuit_f_enum)
-- $8: DESCRIPTION (varchar)
-- $9: ACCESSDESC (varchar)
-- $10: COMMENTS (varchar)
-- $11: SYMBOLOGY (fieldseeker.linelocation_locationsymbology_enum)
-- $12: EXTERNALID (varchar)
-- $13: ACRES (double precision)
-- $14: NEXTACTIONDATESCHEDULED (timestamp)
-- $15: LARVINSPECTINTERVAL (smallint)
-- $16: LENGTH_FT (double precision)
-- $17: WIDTH_FT (double precision)
-- $18: ZONE2 (varchar)
-- $19: LOCATIONNUMBER (integer)
-- $20: GlobalID (uuid)
-- $21: created_user (varchar)
-- $22: created_date (timestamp)
-- $23: last_edited_user (varchar)
-- $24: last_edited_date (timestamp)
-- $25: LASTINSPECTDATE (timestamp)
-- $26: LASTINSPECTBREEDING (varchar)
-- $27: LASTINSPECTAVGLARVAE (double precision)
-- $28: LASTINSPECTAVGPUPAE (double precision)
-- $29: LASTINSPECTLSTAGES (varchar)
-- $30: LASTINSPECTACTIONTAKEN (varchar)
-- $31: LASTINSPECTFIELDSPECIES (varchar)
-- $32: LASTTREATDATE (timestamp)
-- $33: LASTTREATPRODUCT (varchar)
-- $34: LASTTREATQTY (double precision)
-- $35: LASTTREATQTYUNIT (varchar)
-- $36: HECTARES (double precision)
-- $37: LASTINSPECTACTIVITY (varchar)
-- $38: LASTTREATACTIVITY (varchar)
-- $39: LENGTH_METERS (double precision)
-- $40: WIDTH_METERS (double precision)
-- $41: LASTINSPECTCONDITIONS (varchar)
-- $42: WATERORIGIN (fieldseeker.linelocation_linelocation_waterorigin_84723d92_306a_46f4_8ef1_69b55a916008_enum)
-- $43: CreationDate (timestamp)
-- $44: Creator (varchar)
-- $45: EditDate (timestamp)
-- $46: Editor (varchar)
-- $47: JURISDICTION (varchar)
-- $48: Shape__Length (double precision)
