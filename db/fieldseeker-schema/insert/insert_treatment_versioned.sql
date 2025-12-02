-- Prepared statement for conditional insert with versioning for fieldseeker.treatment
-- Only inserts a new version if data has changed

PREPARE insert_treatment_versioned(bigint, fieldseeker.treatment_mosquitoactivity_enum, double precision, fieldseeker.treatment_mosquitoproductareaunit_enum, varchar, double precision, fieldseeker.treatment_mosquitoproductmeasureunit_enum, fieldseeker.treatment_treatment_method_d558ca3ccf43440c8160758253967621_enum, fieldseeker.treatment_treatment_equiptype_45694d79_ff21_42cc_be4f_a0d1def4fba0_enum, varchar, double precision, double precision, fieldseeker.treatment_notinuiwinddirection_enum, double precision, timestamp, timestamp, uuid, fieldseeker.treatment_notinuit_f_enum, varchar, timestamp, varchar, varchar, fieldseeker.treatment_notinuit_f_enum, smallint, varchar, double precision, smallint, smallint, smallint, uuid, double precision, double precision, varchar, uuid, uuid, uuid, uuid, uuid, uuid, uuid, varchar, uuid, double precision, fieldseeker.treatment_treatment_habitat_0afee7eb_f9ea_4707_8483_cccfe60f0d16_enum, double precision, varchar, fieldseeker.treatment_treatment_sitecond_f812e1f64dcb4dc9a75da9d00abe6169_enum, fieldseeker.treatment_treatment_sitecond_5a15bf36fa124280b961f31cd1a9b571_enum, double precision, timestamp, varchar, timestamp, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.treatment
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.treatment
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.treatment (
  objectid, activity, treatarea, areaunit, product, qty, qtyunit, method, equiptype, comments, avetemp, windspeed, winddir, raingauge, startdatetime, enddatetime, insp_id, reviewed, reviewedby, revieweddate, locationname, zone, warningoverride, recordstatus, zone2, treatacres, tirecount, cbcount, containercount, globalid, treatmentlength, treatmenthours, treatmentlengthunits, linelocid, pointlocid, polygonlocid, srid, sdid, barrierrouteid, ulvrouteid, fieldtech, ptaid, flowrate, habitat, treathectares, invloc, temp_sitecond, sitecond, totalcostprodcut, creationdate, creator, editdate, editor, targetspecies,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.activity IS NOT DISTINCT FROM $2 AND
    lv.treatarea IS NOT DISTINCT FROM $3 AND
    lv.areaunit IS NOT DISTINCT FROM $4 AND
    lv.product IS NOT DISTINCT FROM $5 AND
    lv.qty IS NOT DISTINCT FROM $6 AND
    lv.qtyunit IS NOT DISTINCT FROM $7 AND
    lv.method IS NOT DISTINCT FROM $8 AND
    lv.equiptype IS NOT DISTINCT FROM $9 AND
    lv.comments IS NOT DISTINCT FROM $10 AND
    lv.avetemp IS NOT DISTINCT FROM $11 AND
    lv.windspeed IS NOT DISTINCT FROM $12 AND
    lv.winddir IS NOT DISTINCT FROM $13 AND
    lv.raingauge IS NOT DISTINCT FROM $14 AND
    lv.startdatetime IS NOT DISTINCT FROM $15 AND
    lv.enddatetime IS NOT DISTINCT FROM $16 AND
    lv.insp_id IS NOT DISTINCT FROM $17 AND
    lv.reviewed IS NOT DISTINCT FROM $18 AND
    lv.reviewedby IS NOT DISTINCT FROM $19 AND
    lv.revieweddate IS NOT DISTINCT FROM $20 AND
    lv.locationname IS NOT DISTINCT FROM $21 AND
    lv.zone IS NOT DISTINCT FROM $22 AND
    lv.warningoverride IS NOT DISTINCT FROM $23 AND
    lv.recordstatus IS NOT DISTINCT FROM $24 AND
    lv.zone2 IS NOT DISTINCT FROM $25 AND
    lv.treatacres IS NOT DISTINCT FROM $26 AND
    lv.tirecount IS NOT DISTINCT FROM $27 AND
    lv.cbcount IS NOT DISTINCT FROM $28 AND
    lv.containercount IS NOT DISTINCT FROM $29 AND
    lv.globalid IS NOT DISTINCT FROM $30 AND
    lv.treatmentlength IS NOT DISTINCT FROM $31 AND
    lv.treatmenthours IS NOT DISTINCT FROM $32 AND
    lv.treatmentlengthunits IS NOT DISTINCT FROM $33 AND
    lv.linelocid IS NOT DISTINCT FROM $34 AND
    lv.pointlocid IS NOT DISTINCT FROM $35 AND
    lv.polygonlocid IS NOT DISTINCT FROM $36 AND
    lv.srid IS NOT DISTINCT FROM $37 AND
    lv.sdid IS NOT DISTINCT FROM $38 AND
    lv.barrierrouteid IS NOT DISTINCT FROM $39 AND
    lv.ulvrouteid IS NOT DISTINCT FROM $40 AND
    lv.fieldtech IS NOT DISTINCT FROM $41 AND
    lv.ptaid IS NOT DISTINCT FROM $42 AND
    lv.flowrate IS NOT DISTINCT FROM $43 AND
    lv.habitat IS NOT DISTINCT FROM $44 AND
    lv.treathectares IS NOT DISTINCT FROM $45 AND
    lv.invloc IS NOT DISTINCT FROM $46 AND
    lv.temp_sitecond IS NOT DISTINCT FROM $47 AND
    lv.sitecond IS NOT DISTINCT FROM $48 AND
    lv.totalcostprodcut IS NOT DISTINCT FROM $49 AND
    lv.creationdate IS NOT DISTINCT FROM $50 AND
    lv.creator IS NOT DISTINCT FROM $51 AND
    lv.editdate IS NOT DISTINCT FROM $52 AND
    lv.editor IS NOT DISTINCT FROM $53 AND
    lv.targetspecies IS NOT DISTINCT FROM $54
  )
RETURNING *;

-- Example usage: EXECUTE insert_treatment_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: ACTIVITY (fieldseeker.treatment_mosquitoactivity_enum)
-- $3: TREATAREA (double precision)
-- $4: AREAUNIT (fieldseeker.treatment_mosquitoproductareaunit_enum)
-- $5: PRODUCT (varchar)
-- $6: QTY (double precision)
-- $7: QTYUNIT (fieldseeker.treatment_mosquitoproductmeasureunit_enum)
-- $8: METHOD (fieldseeker.treatment_treatment_method_d558ca3ccf43440c8160758253967621_enum)
-- $9: EQUIPTYPE (fieldseeker.treatment_treatment_equiptype_45694d79_ff21_42cc_be4f_a0d1def4fba0_enum)
-- $10: COMMENTS (varchar)
-- $11: AVETEMP (double precision)
-- $12: WINDSPEED (double precision)
-- $13: WINDDIR (fieldseeker.treatment_notinuiwinddirection_enum)
-- $14: RAINGAUGE (double precision)
-- $15: STARTDATETIME (timestamp)
-- $16: ENDDATETIME (timestamp)
-- $17: INSP_ID (uuid)
-- $18: REVIEWED (fieldseeker.treatment_notinuit_f_enum)
-- $19: REVIEWEDBY (varchar)
-- $20: REVIEWEDDATE (timestamp)
-- $21: LOCATIONNAME (varchar)
-- $22: ZONE (varchar)
-- $23: WARNINGOVERRIDE (fieldseeker.treatment_notinuit_f_enum)
-- $24: RECORDSTATUS (smallint)
-- $25: ZONE2 (varchar)
-- $26: TREATACRES (double precision)
-- $27: TIRECOUNT (smallint)
-- $28: CBCOUNT (smallint)
-- $29: CONTAINERCOUNT (smallint)
-- $30: GlobalID (uuid)
-- $31: TREATMENTLENGTH (double precision)
-- $32: TREATMENTHOURS (double precision)
-- $33: TREATMENTLENGTHUNITS (varchar)
-- $34: LINELOCID (uuid)
-- $35: POINTLOCID (uuid)
-- $36: POLYGONLOCID (uuid)
-- $37: SRID (uuid)
-- $38: SDID (uuid)
-- $39: BARRIERROUTEID (uuid)
-- $40: ULVROUTEID (uuid)
-- $41: FIELDTECH (varchar)
-- $42: PTAID (uuid)
-- $43: FLOWRATE (double precision)
-- $44: HABITAT (fieldseeker.treatment_treatment_habitat_0afee7eb_f9ea_4707_8483_cccfe60f0d16_enum)
-- $45: TREATHECTARES (double precision)
-- $46: INVLOC (varchar)
-- $47: temp_SITECOND (fieldseeker.treatment_treatment_sitecond_f812e1f64dcb4dc9a75da9d00abe6169_enum)
-- $48: SITECOND (fieldseeker.treatment_treatment_sitecond_5a15bf36fa124280b961f31cd1a9b571_enum)
-- $49: TotalCostProdcut (double precision)
-- $50: CreationDate (timestamp)
-- $51: Creator (varchar)
-- $52: EditDate (timestamp)
-- $53: Editor (varchar)
-- $54: TARGETSPECIES (varchar)
