-- Prepared statement for conditional insert with versioning for fieldseeker.mosquitoinspection
-- Only inserts a new version if data has changed

PREPARE insert_mosquitoinspection_versioned(bigint, smallint, fieldseeker.mosquitoinspection_mosquitoactivity_enum, fieldseeker.mosquitoinspection_mosquitobreeding_enum, smallint, smallint, smallint, smallint, fieldseeker.mosquitoinspection_mosquitoadultactivity_enum, varchar, fieldseeker.mosquitoinspection_mosquitoinspection_domstage_b7a6c36bccde49a292020de4812cf5ae_enum, fieldseeker.mosquitoinspection_mosquitoinspection_actiontaken_252243d69b0b44ddbdc229c04ec3a8d5_enum, varchar, double precision, double precision, double precision, timestamp, timestamp, fieldseeker.mosquitoinspection_notinuiwinddirection_enum, double precision, double precision, fieldseeker.mosquitoinspection_notinuit_f_enum, varchar, timestamp, varchar, varchar, smallint, varchar, fieldseeker.mosquitoinspection_notinuit_f_enum, smallint, smallint, smallint, fieldseeker.mosquitoinspection_mosquitofieldspecies_enum, uuid, varchar, timestamp, varchar, timestamp, uuid, uuid, uuid, uuid, varchar, fieldseeker.mosquitoinspection_notinuit_f_enum, fieldseeker.mosquitoinspection_notinuit_f_enum, uuid, fieldseeker.mosquitoinspection_mosquitoinspection_sitecond_db7350bc_81e5_401e_858f_cd3e5e5d8a34_enum, smallint, timestamp, varchar, timestamp, varchar, varchar, fieldseeker.mosquitoinspection_notinuit_f_enum, varchar, fieldseeker.mosquitoinspection_mosquitoinspection_adminaction_b74ae1bb_c98b_40f6_8cfa_40e4fd16c270_enum, uuid) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.mosquitoinspection
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.mosquitoinspection
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.mosquitoinspection (
  objectid, numdips, activity, breeding, totlarvae, totpupae, eggs, posdips, adultact, lstages, domstage, actiontaken, comments, avetemp, windspeed, raingauge, startdatetime, enddatetime, winddir, avglarvae, avgpupae, reviewed, reviewedby, revieweddate, locationname, zone, recordstatus, zone2, personalcontact, tirecount, cbcount, containercount, fieldspecies, globalid, created_user, created_date, last_edited_user, last_edited_date, linelocid, pointlocid, polygonlocid, srid, fieldtech, larvaepresent, pupaepresent, sdid, sitecond, positivecontainercount, creationdate, creator, editdate, editor, jurisdiction, visualmonitoring, vmcomments, adminaction, ptaid,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.numdips IS NOT DISTINCT FROM $2 AND
    lv.activity IS NOT DISTINCT FROM $3 AND
    lv.breeding IS NOT DISTINCT FROM $4 AND
    lv.totlarvae IS NOT DISTINCT FROM $5 AND
    lv.totpupae IS NOT DISTINCT FROM $6 AND
    lv.eggs IS NOT DISTINCT FROM $7 AND
    lv.posdips IS NOT DISTINCT FROM $8 AND
    lv.adultact IS NOT DISTINCT FROM $9 AND
    lv.lstages IS NOT DISTINCT FROM $10 AND
    lv.domstage IS NOT DISTINCT FROM $11 AND
    lv.actiontaken IS NOT DISTINCT FROM $12 AND
    lv.comments IS NOT DISTINCT FROM $13 AND
    lv.avetemp IS NOT DISTINCT FROM $14 AND
    lv.windspeed IS NOT DISTINCT FROM $15 AND
    lv.raingauge IS NOT DISTINCT FROM $16 AND
    lv.startdatetime IS NOT DISTINCT FROM $17 AND
    lv.enddatetime IS NOT DISTINCT FROM $18 AND
    lv.winddir IS NOT DISTINCT FROM $19 AND
    lv.avglarvae IS NOT DISTINCT FROM $20 AND
    lv.avgpupae IS NOT DISTINCT FROM $21 AND
    lv.reviewed IS NOT DISTINCT FROM $22 AND
    lv.reviewedby IS NOT DISTINCT FROM $23 AND
    lv.revieweddate IS NOT DISTINCT FROM $24 AND
    lv.locationname IS NOT DISTINCT FROM $25 AND
    lv.zone IS NOT DISTINCT FROM $26 AND
    lv.recordstatus IS NOT DISTINCT FROM $27 AND
    lv.zone2 IS NOT DISTINCT FROM $28 AND
    lv.personalcontact IS NOT DISTINCT FROM $29 AND
    lv.tirecount IS NOT DISTINCT FROM $30 AND
    lv.cbcount IS NOT DISTINCT FROM $31 AND
    lv.containercount IS NOT DISTINCT FROM $32 AND
    lv.fieldspecies IS NOT DISTINCT FROM $33 AND
    lv.globalid IS NOT DISTINCT FROM $34 AND
    lv.created_user IS NOT DISTINCT FROM $35 AND
    lv.created_date IS NOT DISTINCT FROM $36 AND
    lv.last_edited_user IS NOT DISTINCT FROM $37 AND
    lv.last_edited_date IS NOT DISTINCT FROM $38 AND
    lv.linelocid IS NOT DISTINCT FROM $39 AND
    lv.pointlocid IS NOT DISTINCT FROM $40 AND
    lv.polygonlocid IS NOT DISTINCT FROM $41 AND
    lv.srid IS NOT DISTINCT FROM $42 AND
    lv.fieldtech IS NOT DISTINCT FROM $43 AND
    lv.larvaepresent IS NOT DISTINCT FROM $44 AND
    lv.pupaepresent IS NOT DISTINCT FROM $45 AND
    lv.sdid IS NOT DISTINCT FROM $46 AND
    lv.sitecond IS NOT DISTINCT FROM $47 AND
    lv.positivecontainercount IS NOT DISTINCT FROM $48 AND
    lv.creationdate IS NOT DISTINCT FROM $49 AND
    lv.creator IS NOT DISTINCT FROM $50 AND
    lv.editdate IS NOT DISTINCT FROM $51 AND
    lv.editor IS NOT DISTINCT FROM $52 AND
    lv.jurisdiction IS NOT DISTINCT FROM $53 AND
    lv.visualmonitoring IS NOT DISTINCT FROM $54 AND
    lv.vmcomments IS NOT DISTINCT FROM $55 AND
    lv.adminaction IS NOT DISTINCT FROM $56 AND
    lv.ptaid IS NOT DISTINCT FROM $57
  )
RETURNING *;

-- Example usage: EXECUTE insert_mosquitoinspection_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: NUMDIPS (smallint)
-- $3: ACTIVITY (fieldseeker.mosquitoinspection_mosquitoactivity_enum)
-- $4: BREEDING (fieldseeker.mosquitoinspection_mosquitobreeding_enum)
-- $5: TOTLARVAE (smallint)
-- $6: TOTPUPAE (smallint)
-- $7: EGGS (smallint)
-- $8: POSDIPS (smallint)
-- $9: ADULTACT (fieldseeker.mosquitoinspection_mosquitoadultactivity_enum)
-- $10: LSTAGES (varchar)
-- $11: DOMSTAGE (fieldseeker.mosquitoinspection_mosquitoinspection_domstage_b7a6c36bccde49a292020de4812cf5ae_enum)
-- $12: ACTIONTAKEN (fieldseeker.mosquitoinspection_mosquitoinspection_actiontaken_252243d69b0b44ddbdc229c04ec3a8d5_enum)
-- $13: COMMENTS (varchar)
-- $14: AVETEMP (double precision)
-- $15: WINDSPEED (double precision)
-- $16: RAINGAUGE (double precision)
-- $17: STARTDATETIME (timestamp)
-- $18: ENDDATETIME (timestamp)
-- $19: WINDDIR (fieldseeker.mosquitoinspection_notinuiwinddirection_enum)
-- $20: AVGLARVAE (double precision)
-- $21: AVGPUPAE (double precision)
-- $22: REVIEWED (fieldseeker.mosquitoinspection_notinuit_f_enum)
-- $23: REVIEWEDBY (varchar)
-- $24: REVIEWEDDATE (timestamp)
-- $25: LOCATIONNAME (varchar)
-- $26: ZONE (varchar)
-- $27: RECORDSTATUS (smallint)
-- $28: ZONE2 (varchar)
-- $29: PERSONALCONTACT (fieldseeker.mosquitoinspection_notinuit_f_enum)
-- $30: TIRECOUNT (smallint)
-- $31: CBCOUNT (smallint)
-- $32: CONTAINERCOUNT (smallint)
-- $33: FIELDSPECIES (fieldseeker.mosquitoinspection_mosquitofieldspecies_enum)
-- $34: GlobalID (uuid)
-- $35: created_user (varchar)
-- $36: created_date (timestamp)
-- $37: last_edited_user (varchar)
-- $38: last_edited_date (timestamp)
-- $39: LINELOCID (uuid)
-- $40: POINTLOCID (uuid)
-- $41: POLYGONLOCID (uuid)
-- $42: SRID (uuid)
-- $43: FIELDTECH (varchar)
-- $44: LARVAEPRESENT (fieldseeker.mosquitoinspection_notinuit_f_enum)
-- $45: PUPAEPRESENT (fieldseeker.mosquitoinspection_notinuit_f_enum)
-- $46: SDID (uuid)
-- $47: SITECOND (fieldseeker.mosquitoinspection_mosquitoinspection_sitecond_db7350bc_81e5_401e_858f_cd3e5e5d8a34_enum)
-- $48: POSITIVECONTAINERCOUNT (smallint)
-- $49: CreationDate (timestamp)
-- $50: Creator (varchar)
-- $51: EditDate (timestamp)
-- $52: Editor (varchar)
-- $53: JURISDICTION (varchar)
-- $54: VISUALMONITORING (fieldseeker.mosquitoinspection_notinuit_f_enum)
-- $55: VMCOMMENTS (varchar)
-- $56: adminAction (fieldseeker.mosquitoinspection_mosquitoinspection_adminaction_b74ae1bb_c98b_40f6_8cfa_40e4fd16c270_enum)
-- $57: PTAID (uuid)
