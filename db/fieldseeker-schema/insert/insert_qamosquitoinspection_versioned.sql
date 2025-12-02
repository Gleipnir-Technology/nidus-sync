-- Prepared statement for conditional insert with versioning for fieldseeker.qamosquitoinspection
-- Only inserts a new version if data has changed

PREPARE insert_qamosquitoinspection_versioned(bigint, smallint, fieldseeker.qamosquitoinspection_mosquitoaction_enum, varchar, double precision, double precision, double precision, uuid, timestamp, timestamp, varchar, fieldseeker.qamosquitoinspection_notinuit_f_enum, varchar, timestamp, varchar, varchar, smallint, varchar, smallint, smallint, double precision, double precision, fieldseeker.qamosquitoinspection_notinuit_f_enum, fieldseeker.qamosquitoinspection_qasitetype_enum, fieldseeker.qamosquitoinspection_qabreedingpotential_enum, fieldseeker.qamosquitoinspection_notinuit_f_enum, fieldseeker.qamosquitoinspection_notinuit_f_enum, fieldseeker.qamosquitoinspection_qamosquitohabitat_enum, smallint, smallint, smallint, smallint, smallint, fieldseeker.qamosquitoinspection_notinuit_f_enum, fieldseeker.qamosquitoinspection_notinuit_f_enum, fieldseeker.qamosquitoinspection_notinuit_f_enum, fieldseeker.qamosquitoinspection_qalarvaereason_enum, fieldseeker.qamosquitoinspection_qaaquaticorganisms_enum, fieldseeker.qamosquitoinspection_qavegetation_enum, fieldseeker.qamosquitoinspection_qasourcereduction_enum, fieldseeker.qamosquitoinspection_notinuit_f_enum, fieldseeker.qamosquitoinspection_qawatermovement_enum, smallint, fieldseeker.qamosquitoinspection_qawatermovement_enum, smallint, fieldseeker.qamosquitoinspection_qasoilcondition_enum, fieldseeker.qamosquitoinspection_qawaterduration_enum, fieldseeker.qamosquitoinspection_qawatersource_enum, fieldseeker.qamosquitoinspection_qawaterconditions_enum, fieldseeker.qamosquitoinspection_notinuit_f_enum, uuid, uuid, uuid, varchar, timestamp, varchar, timestamp, varchar, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.qamosquitoinspection
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.qamosquitoinspection
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.qamosquitoinspection (
  objectid, posdips, actiontaken, comments, avetemp, windspeed, raingauge, globalid, startdatetime, enddatetime, winddir, reviewed, reviewedby, revieweddate, locationname, zone, recordstatus, zone2, lr, negdips, totalacres, acresbreeding, fish, sitetype, breedingpotential, movingwater, nowaterever, mosquitohabitat, habvalue1, habvalue1percent, habvalue2, habvalue2percent, potential, larvaepresent, larvaeinsidetreatedarea, larvaeoutsidetreatedarea, larvaereason, aquaticorganisms, vegetation, sourcereduction, waterpresent, watermovement1, watermovement1percent, watermovement2, watermovement2percent, soilconditions, waterduration, watersource, waterconditions, adultactivity, linelocid, pointlocid, polygonlocid, created_user, created_date, last_edited_user, last_edited_date, fieldtech, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.posdips IS NOT DISTINCT FROM $2 AND
    lv.actiontaken IS NOT DISTINCT FROM $3 AND
    lv.comments IS NOT DISTINCT FROM $4 AND
    lv.avetemp IS NOT DISTINCT FROM $5 AND
    lv.windspeed IS NOT DISTINCT FROM $6 AND
    lv.raingauge IS NOT DISTINCT FROM $7 AND
    lv.globalid IS NOT DISTINCT FROM $8 AND
    lv.startdatetime IS NOT DISTINCT FROM $9 AND
    lv.enddatetime IS NOT DISTINCT FROM $10 AND
    lv.winddir IS NOT DISTINCT FROM $11 AND
    lv.reviewed IS NOT DISTINCT FROM $12 AND
    lv.reviewedby IS NOT DISTINCT FROM $13 AND
    lv.revieweddate IS NOT DISTINCT FROM $14 AND
    lv.locationname IS NOT DISTINCT FROM $15 AND
    lv.zone IS NOT DISTINCT FROM $16 AND
    lv.recordstatus IS NOT DISTINCT FROM $17 AND
    lv.zone2 IS NOT DISTINCT FROM $18 AND
    lv.lr IS NOT DISTINCT FROM $19 AND
    lv.negdips IS NOT DISTINCT FROM $20 AND
    lv.totalacres IS NOT DISTINCT FROM $21 AND
    lv.acresbreeding IS NOT DISTINCT FROM $22 AND
    lv.fish IS NOT DISTINCT FROM $23 AND
    lv.sitetype IS NOT DISTINCT FROM $24 AND
    lv.breedingpotential IS NOT DISTINCT FROM $25 AND
    lv.movingwater IS NOT DISTINCT FROM $26 AND
    lv.nowaterever IS NOT DISTINCT FROM $27 AND
    lv.mosquitohabitat IS NOT DISTINCT FROM $28 AND
    lv.habvalue1 IS NOT DISTINCT FROM $29 AND
    lv.habvalue1percent IS NOT DISTINCT FROM $30 AND
    lv.habvalue2 IS NOT DISTINCT FROM $31 AND
    lv.habvalue2percent IS NOT DISTINCT FROM $32 AND
    lv.potential IS NOT DISTINCT FROM $33 AND
    lv.larvaepresent IS NOT DISTINCT FROM $34 AND
    lv.larvaeinsidetreatedarea IS NOT DISTINCT FROM $35 AND
    lv.larvaeoutsidetreatedarea IS NOT DISTINCT FROM $36 AND
    lv.larvaereason IS NOT DISTINCT FROM $37 AND
    lv.aquaticorganisms IS NOT DISTINCT FROM $38 AND
    lv.vegetation IS NOT DISTINCT FROM $39 AND
    lv.sourcereduction IS NOT DISTINCT FROM $40 AND
    lv.waterpresent IS NOT DISTINCT FROM $41 AND
    lv.watermovement1 IS NOT DISTINCT FROM $42 AND
    lv.watermovement1percent IS NOT DISTINCT FROM $43 AND
    lv.watermovement2 IS NOT DISTINCT FROM $44 AND
    lv.watermovement2percent IS NOT DISTINCT FROM $45 AND
    lv.soilconditions IS NOT DISTINCT FROM $46 AND
    lv.waterduration IS NOT DISTINCT FROM $47 AND
    lv.watersource IS NOT DISTINCT FROM $48 AND
    lv.waterconditions IS NOT DISTINCT FROM $49 AND
    lv.adultactivity IS NOT DISTINCT FROM $50 AND
    lv.linelocid IS NOT DISTINCT FROM $51 AND
    lv.pointlocid IS NOT DISTINCT FROM $52 AND
    lv.polygonlocid IS NOT DISTINCT FROM $53 AND
    lv.created_user IS NOT DISTINCT FROM $54 AND
    lv.created_date IS NOT DISTINCT FROM $55 AND
    lv.last_edited_user IS NOT DISTINCT FROM $56 AND
    lv.last_edited_date IS NOT DISTINCT FROM $57 AND
    lv.fieldtech IS NOT DISTINCT FROM $58 AND
    lv.creationdate IS NOT DISTINCT FROM $59 AND
    lv.creator IS NOT DISTINCT FROM $60 AND
    lv.editdate IS NOT DISTINCT FROM $61 AND
    lv.editor IS NOT DISTINCT FROM $62
  )
RETURNING *;

-- Example usage: EXECUTE insert_qamosquitoinspection_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: POSDIPS (smallint)
-- $3: ACTIONTAKEN (fieldseeker.qamosquitoinspection_mosquitoaction_enum)
-- $4: COMMENTS (varchar)
-- $5: AVETEMP (double precision)
-- $6: WINDSPEED (double precision)
-- $7: RAINGAUGE (double precision)
-- $8: GlobalID (uuid)
-- $9: STARTDATETIME (timestamp)
-- $10: ENDDATETIME (timestamp)
-- $11: WINDDIR (varchar)
-- $12: REVIEWED (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $13: REVIEWEDBY (varchar)
-- $14: REVIEWEDDATE (timestamp)
-- $15: LOCATIONNAME (varchar)
-- $16: ZONE (varchar)
-- $17: RECORDSTATUS (smallint)
-- $18: ZONE2 (varchar)
-- $19: LR (smallint)
-- $20: NEGDIPS (smallint)
-- $21: TOTALACRES (double precision)
-- $22: ACRESBREEDING (double precision)
-- $23: FISH (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $24: SITETYPE (fieldseeker.qamosquitoinspection_qasitetype_enum)
-- $25: BREEDINGPOTENTIAL (fieldseeker.qamosquitoinspection_qabreedingpotential_enum)
-- $26: MOVINGWATER (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $27: NOWATEREVER (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $28: MOSQUITOHABITAT (fieldseeker.qamosquitoinspection_qamosquitohabitat_enum)
-- $29: HABVALUE1 (smallint)
-- $30: HABVALUE1PERCENT (smallint)
-- $31: HABVALUE2 (smallint)
-- $32: HABVALUE2PERCENT (smallint)
-- $33: POTENTIAL (smallint)
-- $34: LARVAEPRESENT (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $35: LARVAEINSIDETREATEDAREA (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $36: LARVAEOUTSIDETREATEDAREA (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $37: LARVAEREASON (fieldseeker.qamosquitoinspection_qalarvaereason_enum)
-- $38: AQUATICORGANISMS (fieldseeker.qamosquitoinspection_qaaquaticorganisms_enum)
-- $39: VEGETATION (fieldseeker.qamosquitoinspection_qavegetation_enum)
-- $40: SOURCEREDUCTION (fieldseeker.qamosquitoinspection_qasourcereduction_enum)
-- $41: WATERPRESENT (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $42: WATERMOVEMENT1 (fieldseeker.qamosquitoinspection_qawatermovement_enum)
-- $43: WATERMOVEMENT1PERCENT (smallint)
-- $44: WATERMOVEMENT2 (fieldseeker.qamosquitoinspection_qawatermovement_enum)
-- $45: WATERMOVEMENT2PERCENT (smallint)
-- $46: SOILCONDITIONS (fieldseeker.qamosquitoinspection_qasoilcondition_enum)
-- $47: WATERDURATION (fieldseeker.qamosquitoinspection_qawaterduration_enum)
-- $48: WATERSOURCE (fieldseeker.qamosquitoinspection_qawatersource_enum)
-- $49: WATERCONDITIONS (fieldseeker.qamosquitoinspection_qawaterconditions_enum)
-- $50: ADULTACTIVITY (fieldseeker.qamosquitoinspection_notinuit_f_enum)
-- $51: LINELOCID (uuid)
-- $52: POINTLOCID (uuid)
-- $53: POLYGONLOCID (uuid)
-- $54: created_user (varchar)
-- $55: created_date (timestamp)
-- $56: last_edited_user (varchar)
-- $57: last_edited_date (timestamp)
-- $58: FIELDTECH (varchar)
-- $59: CreationDate (timestamp)
-- $60: Creator (varchar)
-- $61: EditDate (timestamp)
-- $62: Editor (varchar)
