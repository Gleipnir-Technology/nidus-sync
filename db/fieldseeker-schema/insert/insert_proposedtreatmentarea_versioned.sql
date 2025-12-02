-- Prepared statement for conditional insert with versioning for fieldseeker.proposedtreatmentarea
-- Only inserts a new version if data has changed

PREPARE insert_proposedtreatmentarea_versioned(bigint, fieldseeker.proposedtreatmentarea_mosquitotreatmentmethod_enum, varchar, varchar, fieldseeker.proposedtreatmentarea_notinuit_f_enum, varchar, timestamp, varchar, timestamp, varchar, fieldseeker.proposedtreatmentarea_notinuit_f_enum, fieldseeker.proposedtreatmentarea_notinuit_f_enum, varchar, double precision, uuid, fieldseeker.proposedtreatmentarea_notinuit_f_enum, varchar, double precision, double precision, varchar, timestamp, varchar, double precision, varchar, fieldseeker.proposedtreatmentarea_locationpriority_enum, timestamp, timestamp, varchar, timestamp, varchar, varchar, double precision, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.proposedtreatmentarea
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.proposedtreatmentarea
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.proposedtreatmentarea (
  objectid, method, comments, zone, reviewed, reviewedby, revieweddate, zone2, completeddate, completedby, completed, issprayroute, name, acres, globalid, exported, targetproduct, targetapprate, hectares, lasttreatactivity, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, priority, duedate, creationdate, creator, editdate, editor, targetspecies, shape__area, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.method IS NOT DISTINCT FROM $2 AND
    lv.comments IS NOT DISTINCT FROM $3 AND
    lv.zone IS NOT DISTINCT FROM $4 AND
    lv.reviewed IS NOT DISTINCT FROM $5 AND
    lv.reviewedby IS NOT DISTINCT FROM $6 AND
    lv.revieweddate IS NOT DISTINCT FROM $7 AND
    lv.zone2 IS NOT DISTINCT FROM $8 AND
    lv.completeddate IS NOT DISTINCT FROM $9 AND
    lv.completedby IS NOT DISTINCT FROM $10 AND
    lv.completed IS NOT DISTINCT FROM $11 AND
    lv.issprayroute IS NOT DISTINCT FROM $12 AND
    lv.name IS NOT DISTINCT FROM $13 AND
    lv.acres IS NOT DISTINCT FROM $14 AND
    lv.globalid IS NOT DISTINCT FROM $15 AND
    lv.exported IS NOT DISTINCT FROM $16 AND
    lv.targetproduct IS NOT DISTINCT FROM $17 AND
    lv.targetapprate IS NOT DISTINCT FROM $18 AND
    lv.hectares IS NOT DISTINCT FROM $19 AND
    lv.lasttreatactivity IS NOT DISTINCT FROM $20 AND
    lv.lasttreatdate IS NOT DISTINCT FROM $21 AND
    lv.lasttreatproduct IS NOT DISTINCT FROM $22 AND
    lv.lasttreatqty IS NOT DISTINCT FROM $23 AND
    lv.lasttreatqtyunit IS NOT DISTINCT FROM $24 AND
    lv.priority IS NOT DISTINCT FROM $25 AND
    lv.duedate IS NOT DISTINCT FROM $26 AND
    lv.creationdate IS NOT DISTINCT FROM $27 AND
    lv.creator IS NOT DISTINCT FROM $28 AND
    lv.editdate IS NOT DISTINCT FROM $29 AND
    lv.editor IS NOT DISTINCT FROM $30 AND
    lv.targetspecies IS NOT DISTINCT FROM $31 AND
    lv.shape__area IS NOT DISTINCT FROM $32 AND
    lv.shape__length IS NOT DISTINCT FROM $33
  )
RETURNING *;

-- Example usage: EXECUTE insert_proposedtreatmentarea_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: METHOD (fieldseeker.proposedtreatmentarea_mosquitotreatmentmethod_enum)
-- $3: COMMENTS (varchar)
-- $4: ZONE (varchar)
-- $5: REVIEWED (fieldseeker.proposedtreatmentarea_notinuit_f_enum)
-- $6: REVIEWEDBY (varchar)
-- $7: REVIEWEDDATE (timestamp)
-- $8: ZONE2 (varchar)
-- $9: COMPLETEDDATE (timestamp)
-- $10: COMPLETEDBY (varchar)
-- $11: COMPLETED (fieldseeker.proposedtreatmentarea_notinuit_f_enum)
-- $12: ISSPRAYROUTE (fieldseeker.proposedtreatmentarea_notinuit_f_enum)
-- $13: NAME (varchar)
-- $14: ACRES (double precision)
-- $15: GlobalID (uuid)
-- $16: EXPORTED (fieldseeker.proposedtreatmentarea_notinuit_f_enum)
-- $17: TARGETPRODUCT (varchar)
-- $18: TARGETAPPRATE (double precision)
-- $19: HECTARES (double precision)
-- $20: LASTTREATACTIVITY (varchar)
-- $21: LASTTREATDATE (timestamp)
-- $22: LASTTREATPRODUCT (varchar)
-- $23: LASTTREATQTY (double precision)
-- $24: LASTTREATQTYUNIT (varchar)
-- $25: PRIORITY (fieldseeker.proposedtreatmentarea_locationpriority_enum)
-- $26: DUEDATE (timestamp)
-- $27: CreationDate (timestamp)
-- $28: Creator (varchar)
-- $29: EditDate (timestamp)
-- $30: Editor (varchar)
-- $31: TARGETSPECIES (varchar)
-- $32: Shape__Area (double precision)
-- $33: Shape__Length (double precision)
