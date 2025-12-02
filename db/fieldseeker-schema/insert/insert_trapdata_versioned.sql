-- Prepared statement for conditional insert with versioning for fieldseeker.trapdata
-- Only inserts a new version if data has changed

PREPARE insert_trapdata_versioned(bigint, fieldseeker.trapdata_mosquitotraptype_enum, fieldseeker.trapdata_notinuitrapactivitytype_enum, timestamp, timestamp, varchar, varchar, varchar, fieldseeker.trapdata_notinuit_f_enum, fieldseeker.trapdata_mosquitositecondition_enum, varchar, smallint, fieldseeker.trapdata_notinuit_f_enum, varchar, timestamp, fieldseeker.trapdata_mosquitotrapcondition_enum, smallint, varchar, varchar, uuid, varchar, timestamp, varchar, timestamp, uuid, varchar, smallint, uuid, double precision, fieldseeker.trapdata_trapdata_winddir_c1a31e05_d0b9_4b22_8800_be127bb3f166_enum, double precision, double precision, double precision, smallint, integer, varchar, varchar, timestamp, varchar, timestamp, varchar, fieldseeker.trapdata_trapdata_lure_25fe542f_077f_4254_8681_76e8f436354b_enum) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.trapdata
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.trapdata
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.trapdata (
  objectid, traptype, trapactivitytype, startdatetime, enddatetime, comments, idbytech, sortbytech, processed, sitecond, locationname, recordstatus, reviewed, reviewedby, revieweddate, trapcondition, trapnights, zone, zone2, globalid, created_user, created_date, last_edited_user, last_edited_date, srid, fieldtech, gatewaysync, loc_id, voltage, winddir, windspeed, avetemp, raingauge, lr, field, vectorsurvtrapdataid, vectorsurvtraplocationid, creationdate, creator, editdate, editor, lure,
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
    lv.traptype IS NOT DISTINCT FROM $2 AND
    lv.trapactivitytype IS NOT DISTINCT FROM $3 AND
    lv.startdatetime IS NOT DISTINCT FROM $4 AND
    lv.enddatetime IS NOT DISTINCT FROM $5 AND
    lv.comments IS NOT DISTINCT FROM $6 AND
    lv.idbytech IS NOT DISTINCT FROM $7 AND
    lv.sortbytech IS NOT DISTINCT FROM $8 AND
    lv.processed IS NOT DISTINCT FROM $9 AND
    lv.sitecond IS NOT DISTINCT FROM $10 AND
    lv.locationname IS NOT DISTINCT FROM $11 AND
    lv.recordstatus IS NOT DISTINCT FROM $12 AND
    lv.reviewed IS NOT DISTINCT FROM $13 AND
    lv.reviewedby IS NOT DISTINCT FROM $14 AND
    lv.revieweddate IS NOT DISTINCT FROM $15 AND
    lv.trapcondition IS NOT DISTINCT FROM $16 AND
    lv.trapnights IS NOT DISTINCT FROM $17 AND
    lv.zone IS NOT DISTINCT FROM $18 AND
    lv.zone2 IS NOT DISTINCT FROM $19 AND
    lv.globalid IS NOT DISTINCT FROM $20 AND
    lv.created_user IS NOT DISTINCT FROM $21 AND
    lv.created_date IS NOT DISTINCT FROM $22 AND
    lv.last_edited_user IS NOT DISTINCT FROM $23 AND
    lv.last_edited_date IS NOT DISTINCT FROM $24 AND
    lv.srid IS NOT DISTINCT FROM $25 AND
    lv.fieldtech IS NOT DISTINCT FROM $26 AND
    lv.gatewaysync IS NOT DISTINCT FROM $27 AND
    lv.loc_id IS NOT DISTINCT FROM $28 AND
    lv.voltage IS NOT DISTINCT FROM $29 AND
    lv.winddir IS NOT DISTINCT FROM $30 AND
    lv.windspeed IS NOT DISTINCT FROM $31 AND
    lv.avetemp IS NOT DISTINCT FROM $32 AND
    lv.raingauge IS NOT DISTINCT FROM $33 AND
    lv.lr IS NOT DISTINCT FROM $34 AND
    lv.field IS NOT DISTINCT FROM $35 AND
    lv.vectorsurvtrapdataid IS NOT DISTINCT FROM $36 AND
    lv.vectorsurvtraplocationid IS NOT DISTINCT FROM $37 AND
    lv.creationdate IS NOT DISTINCT FROM $38 AND
    lv.creator IS NOT DISTINCT FROM $39 AND
    lv.editdate IS NOT DISTINCT FROM $40 AND
    lv.editor IS NOT DISTINCT FROM $41 AND
    lv.lure IS NOT DISTINCT FROM $42
  )
RETURNING *;

-- Example usage: EXECUTE insert_trapdata_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: TRAPTYPE (fieldseeker.trapdata_mosquitotraptype_enum)
-- $3: TRAPACTIVITYTYPE (fieldseeker.trapdata_notinuitrapactivitytype_enum)
-- $4: STARTDATETIME (timestamp)
-- $5: ENDDATETIME (timestamp)
-- $6: COMMENTS (varchar)
-- $7: IDBYTECH (varchar)
-- $8: SORTBYTECH (varchar)
-- $9: PROCESSED (fieldseeker.trapdata_notinuit_f_enum)
-- $10: SITECOND (fieldseeker.trapdata_mosquitositecondition_enum)
-- $11: LOCATIONNAME (varchar)
-- $12: RECORDSTATUS (smallint)
-- $13: REVIEWED (fieldseeker.trapdata_notinuit_f_enum)
-- $14: REVIEWEDBY (varchar)
-- $15: REVIEWEDDATE (timestamp)
-- $16: TRAPCONDITION (fieldseeker.trapdata_mosquitotrapcondition_enum)
-- $17: TRAPNIGHTS (smallint)
-- $18: ZONE (varchar)
-- $19: ZONE2 (varchar)
-- $20: GlobalID (uuid)
-- $21: created_user (varchar)
-- $22: created_date (timestamp)
-- $23: last_edited_user (varchar)
-- $24: last_edited_date (timestamp)
-- $25: SRID (uuid)
-- $26: FIELDTECH (varchar)
-- $27: GATEWAYSYNC (smallint)
-- $28: LOC_ID (uuid)
-- $29: VOLTAGE (double precision)
-- $30: WINDDIR (fieldseeker.trapdata_trapdata_winddir_c1a31e05_d0b9_4b22_8800_be127bb3f166_enum)
-- $31: WINDSPEED (double precision)
-- $32: AVETEMP (double precision)
-- $33: RAINGAUGE (double precision)
-- $34: LR (smallint)
-- $35: Field (integer)
-- $36: VECTORSURVTRAPDATAID (varchar)
-- $37: VECTORSURVTRAPLOCATIONID (varchar)
-- $38: CreationDate (timestamp)
-- $39: Creator (varchar)
-- $40: EditDate (timestamp)
-- $41: Editor (varchar)
-- $42: Lure (fieldseeker.trapdata_trapdata_lure_25fe542f_077f_4254_8681_76e8f436354b_enum)
