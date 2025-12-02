-- Prepared statement for conditional insert with versioning for fieldseeker.pool
-- Only inserts a new version if data has changed

PREPARE insert_pool_versioned(bigint, uuid, timestamp, varchar, timestamp, varchar, varchar, varchar, fieldseeker.pool_notinuit_f_enum, uuid, fieldseeker.pool_pool_testmethod_670efbfba86d41ba8e2d3cab5d749e7f_enum, fieldseeker.pool_pool_diseasetested_0f02232949c04c7e8de820b9b515ed97_enum, fieldseeker.pool_pool_diseasepos_6889f8dd00074874aa726907e78497fa_enum, uuid, varchar, timestamp, varchar, timestamp, fieldseeker.pool_mosquitolabname_enum, smallint, smallint, varchar, varchar, varchar, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.pool
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.pool
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.pool (
  objectid, trapdata_id, datesent, survtech, datetested, testtech, comments, sampleid, processed, lab_id, testmethod, diseasetested, diseasepos, globalid, created_user, created_date, last_edited_user, last_edited_date, lab, poolyear, gatewaysync, vectorsurvcollectionid, vectorsurvpoolid, vectorsurvtrapdataid, creationdate, creator, editdate, editor,
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
    lv.trapdata_id IS NOT DISTINCT FROM $2 AND
    lv.datesent IS NOT DISTINCT FROM $3 AND
    lv.survtech IS NOT DISTINCT FROM $4 AND
    lv.datetested IS NOT DISTINCT FROM $5 AND
    lv.testtech IS NOT DISTINCT FROM $6 AND
    lv.comments IS NOT DISTINCT FROM $7 AND
    lv.sampleid IS NOT DISTINCT FROM $8 AND
    lv.processed IS NOT DISTINCT FROM $9 AND
    lv.lab_id IS NOT DISTINCT FROM $10 AND
    lv.testmethod IS NOT DISTINCT FROM $11 AND
    lv.diseasetested IS NOT DISTINCT FROM $12 AND
    lv.diseasepos IS NOT DISTINCT FROM $13 AND
    lv.globalid IS NOT DISTINCT FROM $14 AND
    lv.created_user IS NOT DISTINCT FROM $15 AND
    lv.created_date IS NOT DISTINCT FROM $16 AND
    lv.last_edited_user IS NOT DISTINCT FROM $17 AND
    lv.last_edited_date IS NOT DISTINCT FROM $18 AND
    lv.lab IS NOT DISTINCT FROM $19 AND
    lv.poolyear IS NOT DISTINCT FROM $20 AND
    lv.gatewaysync IS NOT DISTINCT FROM $21 AND
    lv.vectorsurvcollectionid IS NOT DISTINCT FROM $22 AND
    lv.vectorsurvpoolid IS NOT DISTINCT FROM $23 AND
    lv.vectorsurvtrapdataid IS NOT DISTINCT FROM $24 AND
    lv.creationdate IS NOT DISTINCT FROM $25 AND
    lv.creator IS NOT DISTINCT FROM $26 AND
    lv.editdate IS NOT DISTINCT FROM $27 AND
    lv.editor IS NOT DISTINCT FROM $28
  )
RETURNING *;

-- Example usage: EXECUTE insert_pool_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: TRAPDATA_ID (uuid)
-- $3: DATESENT (timestamp)
-- $4: SURVTECH (varchar)
-- $5: DATETESTED (timestamp)
-- $6: TESTTECH (varchar)
-- $7: COMMENTS (varchar)
-- $8: SAMPLEID (varchar)
-- $9: PROCESSED (fieldseeker.pool_notinuit_f_enum)
-- $10: LAB_ID (uuid)
-- $11: TESTMETHOD (fieldseeker.pool_pool_testmethod_670efbfba86d41ba8e2d3cab5d749e7f_enum)
-- $12: DISEASETESTED (fieldseeker.pool_pool_diseasetested_0f02232949c04c7e8de820b9b515ed97_enum)
-- $13: DISEASEPOS (fieldseeker.pool_pool_diseasepos_6889f8dd00074874aa726907e78497fa_enum)
-- $14: GlobalID (uuid)
-- $15: created_user (varchar)
-- $16: created_date (timestamp)
-- $17: last_edited_user (varchar)
-- $18: last_edited_date (timestamp)
-- $19: LAB (fieldseeker.pool_mosquitolabname_enum)
-- $20: POOLYEAR (smallint)
-- $21: GATEWAYSYNC (smallint)
-- $22: VECTORSURVCOLLECTIONID (varchar)
-- $23: VECTORSURVPOOLID (varchar)
-- $24: VECTORSURVTRAPDATAID (varchar)
-- $25: CreationDate (timestamp)
-- $26: Creator (varchar)
-- $27: EditDate (timestamp)
-- $28: Editor (varchar)
