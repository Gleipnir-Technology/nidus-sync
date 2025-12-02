-- Prepared statement for conditional insert with versioning for fieldseeker.pointlocation
-- Only inserts a new version if data has changed

PREPARE insert_pointlocation_versioned(bigint, varchar, varchar, fieldseeker.pointlocation_pointlocation_habitat_b4d8135a_4979_49c8_8bb3_67ec7230e661_enum, fieldseeker.pointlocation_locationpriority_enum, fieldseeker.pointlocation_pointlocation_usetype_58d62d18ef4f47fc8cb9874df867f89e_enum, fieldseeker.pointlocation_notinuit_f_enum, varchar, varchar, varchar, fieldseeker.pointlocation_locationsymbology_enum, varchar, timestamp, smallint, varchar, integer, uuid, varchar, timestamp, varchar, double precision, double precision, varchar, varchar, varchar, timestamp, varchar, double precision, varchar, varchar, varchar, varchar, fieldseeker.pointlocation_pointlocation_waterorigin_197b22bf_f3eb_4dad_8899_986460f6ea97_enum, double precision, double precision, fieldseeker.pointlocation_pointlocation_assignedtech_9393a162_2474_429d_85be_daa44e4c091f_enum, timestamp, varchar, timestamp, varchar, varchar, fieldseeker.pointlocation_pointlocation_deactivate_reason_dd303085_b33c_4894_8c47_fa847dd9d7c5_enum, integer, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.pointlocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.pointlocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.pointlocation (
  objectid, name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, nextactiondatescheduled, larvinspectinterval, zone2, locationnumber, globalid, stype, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, lastinspectactivity, lasttreatactivity, lastinspectconditions, waterorigin, x, y, assignedtech, creationdate, creator, editdate, editor, jurisdiction, deactivate_reason, scalarpriority, sourcestatus,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44,
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
    lv.nextactiondatescheduled IS NOT DISTINCT FROM $13 AND
    lv.larvinspectinterval IS NOT DISTINCT FROM $14 AND
    lv.zone2 IS NOT DISTINCT FROM $15 AND
    lv.locationnumber IS NOT DISTINCT FROM $16 AND
    lv.globalid IS NOT DISTINCT FROM $17 AND
    lv.stype IS NOT DISTINCT FROM $18 AND
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
    lv.lastinspectactivity IS NOT DISTINCT FROM $30 AND
    lv.lasttreatactivity IS NOT DISTINCT FROM $31 AND
    lv.lastinspectconditions IS NOT DISTINCT FROM $32 AND
    lv.waterorigin IS NOT DISTINCT FROM $33 AND
    lv.x IS NOT DISTINCT FROM $34 AND
    lv.y IS NOT DISTINCT FROM $35 AND
    lv.assignedtech IS NOT DISTINCT FROM $36 AND
    lv.creationdate IS NOT DISTINCT FROM $37 AND
    lv.creator IS NOT DISTINCT FROM $38 AND
    lv.editdate IS NOT DISTINCT FROM $39 AND
    lv.editor IS NOT DISTINCT FROM $40 AND
    lv.jurisdiction IS NOT DISTINCT FROM $41 AND
    lv.deactivate_reason IS NOT DISTINCT FROM $42 AND
    lv.scalarpriority IS NOT DISTINCT FROM $43 AND
    lv.sourcestatus IS NOT DISTINCT FROM $44
  )
RETURNING *;

-- Example usage: EXECUTE insert_pointlocation_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: NAME (varchar)
-- $3: ZONE (varchar)
-- $4: HABITAT (fieldseeker.pointlocation_pointlocation_habitat_b4d8135a_4979_49c8_8bb3_67ec7230e661_enum)
-- $5: PRIORITY (fieldseeker.pointlocation_locationpriority_enum)
-- $6: USETYPE (fieldseeker.pointlocation_pointlocation_usetype_58d62d18ef4f47fc8cb9874df867f89e_enum)
-- $7: ACTIVE (fieldseeker.pointlocation_notinuit_f_enum)
-- $8: DESCRIPTION (varchar)
-- $9: ACCESSDESC (varchar)
-- $10: COMMENTS (varchar)
-- $11: SYMBOLOGY (fieldseeker.pointlocation_locationsymbology_enum)
-- $12: EXTERNALID (varchar)
-- $13: NEXTACTIONDATESCHEDULED (timestamp)
-- $14: LARVINSPECTINTERVAL (smallint)
-- $15: ZONE2 (varchar)
-- $16: LOCATIONNUMBER (integer)
-- $17: GlobalID (uuid)
-- $18: STYPE (varchar)
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
-- $30: LASTINSPECTACTIVITY (varchar)
-- $31: LASTTREATACTIVITY (varchar)
-- $32: LASTINSPECTCONDITIONS (varchar)
-- $33: WATERORIGIN (fieldseeker.pointlocation_pointlocation_waterorigin_197b22bf_f3eb_4dad_8899_986460f6ea97_enum)
-- $34: X (double precision)
-- $35: Y (double precision)
-- $36: ASSIGNEDTECH (fieldseeker.pointlocation_pointlocation_assignedtech_9393a162_2474_429d_85be_daa44e4c091f_enum)
-- $37: CreationDate (timestamp)
-- $38: Creator (varchar)
-- $39: EditDate (timestamp)
-- $40: Editor (varchar)
-- $41: JURISDICTION (varchar)
-- $42: deactivate_reason (fieldseeker.pointlocation_pointlocation_deactivate_reason_dd303085_b33c_4894_8c47_fa847dd9d7c5_enum)
-- $43: scalarPriority (integer)
-- $44: sourceStatus (varchar)
