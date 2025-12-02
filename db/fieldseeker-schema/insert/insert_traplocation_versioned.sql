-- Prepared statement for conditional insert with versioning for fieldseeker.traplocation
-- Only inserts a new version if data has changed

PREPARE insert_traplocation_versioned(bigint, varchar, varchar, fieldseeker.traplocation_traplocation_habitat_5c349680f5ff40b1aeca88c17993e8f3_enum, fieldseeker.traplocation_traplocation_priority_680fb011063b41d59f39271c959b857f_enum, fieldseeker.traplocation_traplocation_usetype_5e0eff9231fb404c98cc53c1d49a2193_enum, fieldseeker.traplocation_notinuit_f_enum, varchar, fieldseeker.traplocation_traplocation_accessdesc_154cbd10_4524_4e3a_8ca0_f099ec86556a_enum, varchar, varchar, timestamp, varchar, integer, uuid, varchar, timestamp, varchar, timestamp, smallint, integer, integer, integer, varchar, timestamp, varchar, timestamp, varchar, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.traplocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.traplocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.traplocation (
  objectid, name, zone, habitat, priority, usetype, active, description, accessdesc, comments, externalid, nextactiondatescheduled, zone2, locationnumber, globalid, created_user, created_date, last_edited_user, last_edited_date, gatewaysync, route, set_dow, route_order, vectorsurvsiteid, creationdate, creator, editdate, editor, h3r7, h3r8,
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
    lv.name IS NOT DISTINCT FROM $2 AND
    lv.zone IS NOT DISTINCT FROM $3 AND
    lv.habitat IS NOT DISTINCT FROM $4 AND
    lv.priority IS NOT DISTINCT FROM $5 AND
    lv.usetype IS NOT DISTINCT FROM $6 AND
    lv.active IS NOT DISTINCT FROM $7 AND
    lv.description IS NOT DISTINCT FROM $8 AND
    lv.accessdesc IS NOT DISTINCT FROM $9 AND
    lv.comments IS NOT DISTINCT FROM $10 AND
    lv.externalid IS NOT DISTINCT FROM $11 AND
    lv.nextactiondatescheduled IS NOT DISTINCT FROM $12 AND
    lv.zone2 IS NOT DISTINCT FROM $13 AND
    lv.locationnumber IS NOT DISTINCT FROM $14 AND
    lv.globalid IS NOT DISTINCT FROM $15 AND
    lv.created_user IS NOT DISTINCT FROM $16 AND
    lv.created_date IS NOT DISTINCT FROM $17 AND
    lv.last_edited_user IS NOT DISTINCT FROM $18 AND
    lv.last_edited_date IS NOT DISTINCT FROM $19 AND
    lv.gatewaysync IS NOT DISTINCT FROM $20 AND
    lv.route IS NOT DISTINCT FROM $21 AND
    lv.set_dow IS NOT DISTINCT FROM $22 AND
    lv.route_order IS NOT DISTINCT FROM $23 AND
    lv.vectorsurvsiteid IS NOT DISTINCT FROM $24 AND
    lv.creationdate IS NOT DISTINCT FROM $25 AND
    lv.creator IS NOT DISTINCT FROM $26 AND
    lv.editdate IS NOT DISTINCT FROM $27 AND
    lv.editor IS NOT DISTINCT FROM $28 AND
    lv.h3r7 IS NOT DISTINCT FROM $29 AND
    lv.h3r8 IS NOT DISTINCT FROM $30
  )
RETURNING *;

-- Example usage: EXECUTE insert_traplocation_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: NAME (varchar)
-- $3: ZONE (varchar)
-- $4: HABITAT (fieldseeker.traplocation_traplocation_habitat_5c349680f5ff40b1aeca88c17993e8f3_enum)
-- $5: PRIORITY (fieldseeker.traplocation_traplocation_priority_680fb011063b41d59f39271c959b857f_enum)
-- $6: USETYPE (fieldseeker.traplocation_traplocation_usetype_5e0eff9231fb404c98cc53c1d49a2193_enum)
-- $7: ACTIVE (fieldseeker.traplocation_notinuit_f_enum)
-- $8: DESCRIPTION (varchar)
-- $9: ACCESSDESC (fieldseeker.traplocation_traplocation_accessdesc_154cbd10_4524_4e3a_8ca0_f099ec86556a_enum)
-- $10: COMMENTS (varchar)
-- $11: EXTERNALID (varchar)
-- $12: NEXTACTIONDATESCHEDULED (timestamp)
-- $13: ZONE2 (varchar)
-- $14: LOCATIONNUMBER (integer)
-- $15: GlobalID (uuid)
-- $16: created_user (varchar)
-- $17: created_date (timestamp)
-- $18: last_edited_user (varchar)
-- $19: last_edited_date (timestamp)
-- $20: GATEWAYSYNC (smallint)
-- $21: route (integer)
-- $22: set_dow (integer)
-- $23: route_order (integer)
-- $24: VECTORSURVSITEID (varchar)
-- $25: CreationDate (timestamp)
-- $26: Creator (varchar)
-- $27: EditDate (timestamp)
-- $28: Editor (varchar)
-- $29: h3r7 (varchar)
-- $30: h3r8 (varchar)
