-- Prepared statement for conditional insert with versioning for fieldseeker.timecard
-- Only inserts a new version if data has changed

PREPARE insert_timecard_versioned(bigint, fieldseeker.timecard_timecard_activity_451e67260c084304a35457170dc13366_enum, timestamp, timestamp, varchar, varchar, fieldseeker.timecard_timecardequipmenttype_enum, varchar, varchar, varchar, uuid, varchar, timestamp, varchar, timestamp, uuid, uuid, uuid, uuid, uuid, uuid, uuid, varchar, timestamp, varchar, timestamp, varchar, uuid) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.timecard
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.timecard
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.timecard (
  objectid, activity, startdatetime, enddatetime, comments, externalid, equiptype, locationname, zone, zone2, globalid, created_user, created_date, last_edited_user, last_edited_date, linelocid, pointlocid, polygonlocid, lclocid, samplelocid, srid, traplocid, fieldtech, creationdate, creator, editdate, editor, rodentlocid,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.activity IS NOT DISTINCT FROM $2 AND
    lv.startdatetime IS NOT DISTINCT FROM $3 AND
    lv.enddatetime IS NOT DISTINCT FROM $4 AND
    lv.comments IS NOT DISTINCT FROM $5 AND
    lv.externalid IS NOT DISTINCT FROM $6 AND
    lv.equiptype IS NOT DISTINCT FROM $7 AND
    lv.locationname IS NOT DISTINCT FROM $8 AND
    lv.zone IS NOT DISTINCT FROM $9 AND
    lv.zone2 IS NOT DISTINCT FROM $10 AND
    lv.globalid IS NOT DISTINCT FROM $11 AND
    lv.created_user IS NOT DISTINCT FROM $12 AND
    lv.created_date IS NOT DISTINCT FROM $13 AND
    lv.last_edited_user IS NOT DISTINCT FROM $14 AND
    lv.last_edited_date IS NOT DISTINCT FROM $15 AND
    lv.linelocid IS NOT DISTINCT FROM $16 AND
    lv.pointlocid IS NOT DISTINCT FROM $17 AND
    lv.polygonlocid IS NOT DISTINCT FROM $18 AND
    lv.lclocid IS NOT DISTINCT FROM $19 AND
    lv.samplelocid IS NOT DISTINCT FROM $20 AND
    lv.srid IS NOT DISTINCT FROM $21 AND
    lv.traplocid IS NOT DISTINCT FROM $22 AND
    lv.fieldtech IS NOT DISTINCT FROM $23 AND
    lv.creationdate IS NOT DISTINCT FROM $24 AND
    lv.creator IS NOT DISTINCT FROM $25 AND
    lv.editdate IS NOT DISTINCT FROM $26 AND
    lv.editor IS NOT DISTINCT FROM $27 AND
    lv.rodentlocid IS NOT DISTINCT FROM $28
  )
RETURNING *;

-- Example usage: EXECUTE insert_timecard_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: ACTIVITY (fieldseeker.timecard_timecard_activity_451e67260c084304a35457170dc13366_enum)
-- $3: STARTDATETIME (timestamp)
-- $4: ENDDATETIME (timestamp)
-- $5: COMMENTS (varchar)
-- $6: EXTERNALID (varchar)
-- $7: EQUIPTYPE (fieldseeker.timecard_timecardequipmenttype_enum)
-- $8: LOCATIONNAME (varchar)
-- $9: ZONE (varchar)
-- $10: ZONE2 (varchar)
-- $11: GlobalID (uuid)
-- $12: created_user (varchar)
-- $13: created_date (timestamp)
-- $14: last_edited_user (varchar)
-- $15: last_edited_date (timestamp)
-- $16: LINELOCID (uuid)
-- $17: POINTLOCID (uuid)
-- $18: POLYGONLOCID (uuid)
-- $19: LCLOCID (uuid)
-- $20: SAMPLELOCID (uuid)
-- $21: SRID (uuid)
-- $22: TRAPLOCID (uuid)
-- $23: FIELDTECH (varchar)
-- $24: CreationDate (timestamp)
-- $25: Creator (varchar)
-- $26: EditDate (timestamp)
-- $27: Editor (varchar)
-- $28: RODENTLOCID (uuid)
