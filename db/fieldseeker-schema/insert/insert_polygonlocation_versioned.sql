-- Prepared statement for conditional insert with versioning for fieldseeker.polygonlocation
-- Only inserts a new version if data has changed

PREPARE insert_polygonlocation_versioned(bigint, varchar, varchar, fieldseeker.polygonlocation_polygonlocation_habitat_45e9dde79ac84d959df8b65ba7d5dafd_enum, fieldseeker.polygonlocation_locationpriority_enum, fieldseeker.polygonlocation_polygonlocation_usetype_e546154cb9544b9aa8e7b13e8e258b27_enum, fieldseeker.polygonlocation_notinuit_f_enum, varchar, varchar, varchar, fieldseeker.polygonlocation_locationsymbology_enum, varchar, double precision, timestamp, smallint, varchar, integer, uuid, timestamp, varchar, double precision, double precision, varchar, varchar, varchar, timestamp, varchar, double precision, varchar, double precision, varchar, varchar, varchar, fieldseeker.polygonlocation_polygonlocation_waterorigin_e9018e92_5f47_4ff9_8a7c_b818d848dc7a_enum, varchar, timestamp, varchar, timestamp, varchar, varchar, double precision, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.polygonlocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.polygonlocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.polygonlocation (
  objectid, name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, acres, nextactiondatescheduled, larvinspectinterval, zone2, locationnumber, globalid, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, hectares, lastinspectactivity, lasttreatactivity, lastinspectconditions, waterorigin, filter, creationdate, creator, editdate, editor, jurisdiction, shape__area, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42,
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
    lv.zone2 IS NOT DISTINCT FROM $16 AND
    lv.locationnumber IS NOT DISTINCT FROM $17 AND
    lv.globalid IS NOT DISTINCT FROM $18 AND
    lv.lastinspectdate IS NOT DISTINCT FROM $19 AND
    lv.lastinspectbreeding IS NOT DISTINCT FROM $20 AND
    lv.lastinspectavglarvae IS NOT DISTINCT FROM $21 AND
    lv.lastinspectavgpupae IS NOT DISTINCT FROM $22 AND
    lv.lastinspectlstages IS NOT DISTINCT FROM $23 AND
    lv.lastinspectactiontaken IS NOT DISTINCT FROM $24 AND
    lv.lastinspectfieldspecies IS NOT DISTINCT FROM $25 AND
    lv.lasttreatdate IS NOT DISTINCT FROM $26 AND
    lv.lasttreatproduct IS NOT DISTINCT FROM $27 AND
    lv.lasttreatqty IS NOT DISTINCT FROM $28 AND
    lv.lasttreatqtyunit IS NOT DISTINCT FROM $29 AND
    lv.hectares IS NOT DISTINCT FROM $30 AND
    lv.lastinspectactivity IS NOT DISTINCT FROM $31 AND
    lv.lasttreatactivity IS NOT DISTINCT FROM $32 AND
    lv.lastinspectconditions IS NOT DISTINCT FROM $33 AND
    lv.waterorigin IS NOT DISTINCT FROM $34 AND
    lv.filter IS NOT DISTINCT FROM $35 AND
    lv.creationdate IS NOT DISTINCT FROM $36 AND
    lv.creator IS NOT DISTINCT FROM $37 AND
    lv.editdate IS NOT DISTINCT FROM $38 AND
    lv.editor IS NOT DISTINCT FROM $39 AND
    lv.jurisdiction IS NOT DISTINCT FROM $40 AND
    lv.shape__area IS NOT DISTINCT FROM $41 AND
    lv.shape__length IS NOT DISTINCT FROM $42
  )
RETURNING *;

-- Example usage: EXECUTE insert_polygonlocation_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: NAME (varchar)
-- $3: ZONE (varchar)
-- $4: HABITAT (fieldseeker.polygonlocation_polygonlocation_habitat_45e9dde79ac84d959df8b65ba7d5dafd_enum)
-- $5: PRIORITY (fieldseeker.polygonlocation_locationpriority_enum)
-- $6: USETYPE (fieldseeker.polygonlocation_polygonlocation_usetype_e546154cb9544b9aa8e7b13e8e258b27_enum)
-- $7: ACTIVE (fieldseeker.polygonlocation_notinuit_f_enum)
-- $8: DESCRIPTION (varchar)
-- $9: ACCESSDESC (varchar)
-- $10: COMMENTS (varchar)
-- $11: SYMBOLOGY (fieldseeker.polygonlocation_locationsymbology_enum)
-- $12: EXTERNALID (varchar)
-- $13: ACRES (double precision)
-- $14: NEXTACTIONDATESCHEDULED (timestamp)
-- $15: LARVINSPECTINTERVAL (smallint)
-- $16: ZONE2 (varchar)
-- $17: LOCATIONNUMBER (integer)
-- $18: GlobalID (uuid)
-- $19: LASTINSPECTDATE (timestamp)
-- $20: LASTINSPECTBREEDING (varchar)
-- $21: LASTINSPECTAVGLARVAE (double precision)
-- $22: LASTINSPECTAVGPUPAE (double precision)
-- $23: LASTINSPECTLSTAGES (varchar)
-- $24: LASTINSPECTACTIONTAKEN (varchar)
-- $25: LASTINSPECTFIELDSPECIES (varchar)
-- $26: LASTTREATDATE (timestamp)
-- $27: LASTTREATPRODUCT (varchar)
-- $28: LASTTREATQTY (double precision)
-- $29: LASTTREATQTYUNIT (varchar)
-- $30: HECTARES (double precision)
-- $31: LASTINSPECTACTIVITY (varchar)
-- $32: LASTTREATACTIVITY (varchar)
-- $33: LASTINSPECTCONDITIONS (varchar)
-- $34: WATERORIGIN (fieldseeker.polygonlocation_polygonlocation_waterorigin_e9018e92_5f47_4ff9_8a7c_b818d848dc7a_enum)
-- $35: Filter (varchar)
-- $36: CreationDate (timestamp)
-- $37: Creator (varchar)
-- $38: EditDate (timestamp)
-- $39: Editor (varchar)
-- $40: JURISDICTION (varchar)
-- $41: Shape__Area (double precision)
-- $42: Shape__Length (double precision)
