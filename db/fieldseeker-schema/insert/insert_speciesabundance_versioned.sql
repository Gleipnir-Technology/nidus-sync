-- Prepared statement for conditional insert with versioning for fieldseeker.speciesabundance
-- Only inserts a new version if data has changed

PREPARE insert_speciesabundance_versioned(bigint, uuid, varchar, smallint, smallint, smallint, smallint, smallint, smallint, fieldseeker.speciesabundance_notinuit_f_enum, uuid, varchar, timestamp, varchar, timestamp, smallint, smallint, integer, integer, timestamp, varchar, timestamp, varchar, integer, double precision, double precision, double precision, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.speciesabundance
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.speciesabundance
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.speciesabundance (
  objectid, trapdata_id, species, males, unknown, bloodedfem, gravidfem, larvae, poolstogen, processed, globalid, created_user, created_date, last_edited_user, last_edited_date, pupae, eggs, females, total, creationdate, creator, editdate, editor, yearweek, globalzscore, r7score, r8score, h3r7, h3r8,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.trapdata_id IS NOT DISTINCT FROM $2 AND
    lv.species IS NOT DISTINCT FROM $3 AND
    lv.males IS NOT DISTINCT FROM $4 AND
    lv.unknown IS NOT DISTINCT FROM $5 AND
    lv.bloodedfem IS NOT DISTINCT FROM $6 AND
    lv.gravidfem IS NOT DISTINCT FROM $7 AND
    lv.larvae IS NOT DISTINCT FROM $8 AND
    lv.poolstogen IS NOT DISTINCT FROM $9 AND
    lv.processed IS NOT DISTINCT FROM $10 AND
    lv.globalid IS NOT DISTINCT FROM $11 AND
    lv.created_user IS NOT DISTINCT FROM $12 AND
    lv.created_date IS NOT DISTINCT FROM $13 AND
    lv.last_edited_user IS NOT DISTINCT FROM $14 AND
    lv.last_edited_date IS NOT DISTINCT FROM $15 AND
    lv.pupae IS NOT DISTINCT FROM $16 AND
    lv.eggs IS NOT DISTINCT FROM $17 AND
    lv.females IS NOT DISTINCT FROM $18 AND
    lv.total IS NOT DISTINCT FROM $19 AND
    lv.creationdate IS NOT DISTINCT FROM $20 AND
    lv.creator IS NOT DISTINCT FROM $21 AND
    lv.editdate IS NOT DISTINCT FROM $22 AND
    lv.editor IS NOT DISTINCT FROM $23 AND
    lv.yearweek IS NOT DISTINCT FROM $24 AND
    lv.globalzscore IS NOT DISTINCT FROM $25 AND
    lv.r7score IS NOT DISTINCT FROM $26 AND
    lv.r8score IS NOT DISTINCT FROM $27 AND
    lv.h3r7 IS NOT DISTINCT FROM $28 AND
    lv.h3r8 IS NOT DISTINCT FROM $29
  )
RETURNING *;

-- Example usage: EXECUTE insert_speciesabundance_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: TRAPDATA_ID (uuid)
-- $3: SPECIES (varchar)
-- $4: MALES (smallint)
-- $5: UNKNOWN (smallint)
-- $6: BLOODEDFEM (smallint)
-- $7: GRAVIDFEM (smallint)
-- $8: LARVAE (smallint)
-- $9: POOLSTOGEN (smallint)
-- $10: PROCESSED (fieldseeker.speciesabundance_notinuit_f_enum)
-- $11: GlobalID (uuid)
-- $12: created_user (varchar)
-- $13: created_date (timestamp)
-- $14: last_edited_user (varchar)
-- $15: last_edited_date (timestamp)
-- $16: PUPAE (smallint)
-- $17: EGGS (smallint)
-- $18: FEMALES (integer)
-- $19: TOTAL (integer)
-- $20: CreationDate (timestamp)
-- $21: Creator (varchar)
-- $22: EditDate (timestamp)
-- $23: Editor (varchar)
-- $24: yearWeek (integer)
-- $25: globalZScore (double precision)
-- $26: r7Score (double precision)
-- $27: r8Score (double precision)
-- $28: h3r7 (varchar)
-- $29: h3r8 (varchar)
