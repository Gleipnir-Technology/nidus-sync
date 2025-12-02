-- Prepared statement for conditional insert with versioning for fieldseeker.samplecollection
-- Only inserts a new version if data has changed

PREPARE insert_samplecollection_versioned(bigint, uuid, timestamp, timestamp, fieldseeker.samplecollection_mosquitositecondition_enum, varchar, varchar, timestamp, timestamp, varchar, varchar, fieldseeker.samplecollection_notinuit_f_enum, fieldseeker.samplecollection_mosquitosampletype_enum, fieldseeker.samplecollection_mosquitosamplecondition_enum, fieldseeker.samplecollection_mosquitosamplespecies_enum, fieldseeker.samplecollection_notinuisex_enum, double precision, double precision, fieldseeker.samplecollection_notinuiwinddirection_enum, double precision, fieldseeker.samplecollection_mosquitoactivity_enum, fieldseeker.samplecollection_mosquitotestmethod_enum, fieldseeker.samplecollection_mosquitodisease_enum, fieldseeker.samplecollection_mosquitodisease_enum, fieldseeker.samplecollection_notinuit_f_enum, varchar, timestamp, varchar, varchar, smallint, varchar, uuid, varchar, timestamp, varchar, timestamp, fieldseeker.samplecollection_mosquitolabname_enum, varchar, uuid, smallint, uuid, smallint, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.samplecollection
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.samplecollection
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.samplecollection (
  objectid, loc_id, startdatetime, enddatetime, sitecond, sampleid, survtech, datesent, datetested, testtech, comments, processed, sampletype, samplecond, species, sex, avetemp, windspeed, winddir, raingauge, activity, testmethod, diseasetested, diseasepos, reviewed, reviewedby, revieweddate, locationname, zone, recordstatus, zone2, globalid, created_user, created_date, last_edited_user, last_edited_date, lab, fieldtech, flockid, samplecount, chickenid, gatewaysync, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.loc_id IS NOT DISTINCT FROM $2 AND
    lv.startdatetime IS NOT DISTINCT FROM $3 AND
    lv.enddatetime IS NOT DISTINCT FROM $4 AND
    lv.sitecond IS NOT DISTINCT FROM $5 AND
    lv.sampleid IS NOT DISTINCT FROM $6 AND
    lv.survtech IS NOT DISTINCT FROM $7 AND
    lv.datesent IS NOT DISTINCT FROM $8 AND
    lv.datetested IS NOT DISTINCT FROM $9 AND
    lv.testtech IS NOT DISTINCT FROM $10 AND
    lv.comments IS NOT DISTINCT FROM $11 AND
    lv.processed IS NOT DISTINCT FROM $12 AND
    lv.sampletype IS NOT DISTINCT FROM $13 AND
    lv.samplecond IS NOT DISTINCT FROM $14 AND
    lv.species IS NOT DISTINCT FROM $15 AND
    lv.sex IS NOT DISTINCT FROM $16 AND
    lv.avetemp IS NOT DISTINCT FROM $17 AND
    lv.windspeed IS NOT DISTINCT FROM $18 AND
    lv.winddir IS NOT DISTINCT FROM $19 AND
    lv.raingauge IS NOT DISTINCT FROM $20 AND
    lv.activity IS NOT DISTINCT FROM $21 AND
    lv.testmethod IS NOT DISTINCT FROM $22 AND
    lv.diseasetested IS NOT DISTINCT FROM $23 AND
    lv.diseasepos IS NOT DISTINCT FROM $24 AND
    lv.reviewed IS NOT DISTINCT FROM $25 AND
    lv.reviewedby IS NOT DISTINCT FROM $26 AND
    lv.revieweddate IS NOT DISTINCT FROM $27 AND
    lv.locationname IS NOT DISTINCT FROM $28 AND
    lv.zone IS NOT DISTINCT FROM $29 AND
    lv.recordstatus IS NOT DISTINCT FROM $30 AND
    lv.zone2 IS NOT DISTINCT FROM $31 AND
    lv.globalid IS NOT DISTINCT FROM $32 AND
    lv.created_user IS NOT DISTINCT FROM $33 AND
    lv.created_date IS NOT DISTINCT FROM $34 AND
    lv.last_edited_user IS NOT DISTINCT FROM $35 AND
    lv.last_edited_date IS NOT DISTINCT FROM $36 AND
    lv.lab IS NOT DISTINCT FROM $37 AND
    lv.fieldtech IS NOT DISTINCT FROM $38 AND
    lv.flockid IS NOT DISTINCT FROM $39 AND
    lv.samplecount IS NOT DISTINCT FROM $40 AND
    lv.chickenid IS NOT DISTINCT FROM $41 AND
    lv.gatewaysync IS NOT DISTINCT FROM $42 AND
    lv.creationdate IS NOT DISTINCT FROM $43 AND
    lv.creator IS NOT DISTINCT FROM $44 AND
    lv.editdate IS NOT DISTINCT FROM $45 AND
    lv.editor IS NOT DISTINCT FROM $46
  )
RETURNING *;

-- Example usage: EXECUTE insert_samplecollection_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: LOC_ID (uuid)
-- $3: STARTDATETIME (timestamp)
-- $4: ENDDATETIME (timestamp)
-- $5: SITECOND (fieldseeker.samplecollection_mosquitositecondition_enum)
-- $6: SAMPLEID (varchar)
-- $7: SURVTECH (varchar)
-- $8: DATESENT (timestamp)
-- $9: DATETESTED (timestamp)
-- $10: TESTTECH (varchar)
-- $11: COMMENTS (varchar)
-- $12: PROCESSED (fieldseeker.samplecollection_notinuit_f_enum)
-- $13: SAMPLETYPE (fieldseeker.samplecollection_mosquitosampletype_enum)
-- $14: SAMPLECOND (fieldseeker.samplecollection_mosquitosamplecondition_enum)
-- $15: SPECIES (fieldseeker.samplecollection_mosquitosamplespecies_enum)
-- $16: SEX (fieldseeker.samplecollection_notinuisex_enum)
-- $17: AVETEMP (double precision)
-- $18: WINDSPEED (double precision)
-- $19: WINDDIR (fieldseeker.samplecollection_notinuiwinddirection_enum)
-- $20: RAINGAUGE (double precision)
-- $21: ACTIVITY (fieldseeker.samplecollection_mosquitoactivity_enum)
-- $22: TESTMETHOD (fieldseeker.samplecollection_mosquitotestmethod_enum)
-- $23: DISEASETESTED (fieldseeker.samplecollection_mosquitodisease_enum)
-- $24: DISEASEPOS (fieldseeker.samplecollection_mosquitodisease_enum)
-- $25: REVIEWED (fieldseeker.samplecollection_notinuit_f_enum)
-- $26: REVIEWEDBY (varchar)
-- $27: REVIEWEDDATE (timestamp)
-- $28: LOCATIONNAME (varchar)
-- $29: ZONE (varchar)
-- $30: RECORDSTATUS (smallint)
-- $31: ZONE2 (varchar)
-- $32: GlobalID (uuid)
-- $33: created_user (varchar)
-- $34: created_date (timestamp)
-- $35: last_edited_user (varchar)
-- $36: last_edited_date (timestamp)
-- $37: LAB (fieldseeker.samplecollection_mosquitolabname_enum)
-- $38: FIELDTECH (varchar)
-- $39: FLOCKID (uuid)
-- $40: SAMPLECOUNT (smallint)
-- $41: CHICKENID (uuid)
-- $42: GATEWAYSYNC (smallint)
-- $43: CreationDate (timestamp)
-- $44: Creator (varchar)
-- $45: EditDate (timestamp)
-- $46: Editor (varchar)
