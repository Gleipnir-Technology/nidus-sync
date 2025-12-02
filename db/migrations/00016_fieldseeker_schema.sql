-- +goose Up
-- Table definition for fieldseeker.ContainerRelate
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.containerrelate_mosquitocontainertype_enum AS ENUM (
  'Aquarium',
  'Flower pot',
  '5 gallon bucket',
  'Fountain',
  'Bird bath'
);

CREATE TABLE fieldseeker.containerrelate (
  objectid BIGSERIAL NOT NULL,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  inspsampleid UUID,
  mosquitoinspid UUID,
  treatmentid UUID,
  containertype fieldseeker.containerrelate_mosquitocontainertype_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.containerrelate.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.containerrelate.containertype IS 'Container Type';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_containerrelate_versioned(bigint, uuid, varchar, timestamp, varchar, timestamp, uuid, uuid, uuid, fieldseeker.containerrelate_mosquitocontainertype_enum, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.containerrelate
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.containerrelate
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.containerrelate (
  objectid, globalid, created_user, created_date, last_edited_user, last_edited_date, inspsampleid, mosquitoinspid, treatmentid, containertype, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.globalid IS NOT DISTINCT FROM $2 AND
    lv.created_user IS NOT DISTINCT FROM $3 AND
    lv.created_date IS NOT DISTINCT FROM $4 AND
    lv.last_edited_user IS NOT DISTINCT FROM $5 AND
    lv.last_edited_date IS NOT DISTINCT FROM $6 AND
    lv.inspsampleid IS NOT DISTINCT FROM $7 AND
    lv.mosquitoinspid IS NOT DISTINCT FROM $8 AND
    lv.treatmentid IS NOT DISTINCT FROM $9 AND
    lv.containertype IS NOT DISTINCT FROM $10 AND
    lv.creationdate IS NOT DISTINCT FROM $11 AND
    lv.creator IS NOT DISTINCT FROM $12 AND
    lv.editdate IS NOT DISTINCT FROM $13 AND
    lv.editor IS NOT DISTINCT FROM $14
  )
RETURNING *;

-- Example usage: EXECUTE insert_containerrelate_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.FieldScoutingLog
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.fieldscoutinglog_fieldscoutingsymbology_enum AS ENUM (
  '0',
  '1',
  '2',
  '3'
);

CREATE TABLE fieldseeker.fieldscoutinglog (
  objectid BIGSERIAL NOT NULL,
  status fieldseeker.fieldscoutinglog_fieldscoutingsymbology_enum,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.fieldscoutinglog.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.fieldscoutinglog.status IS 'Status';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_fieldscoutinglog_versioned(bigint, fieldseeker.fieldscoutinglog_fieldscoutingsymbology_enum, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.fieldscoutinglog
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.fieldscoutinglog
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.fieldscoutinglog (
  objectid, status, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.status IS NOT DISTINCT FROM $2 AND
    lv.globalid IS NOT DISTINCT FROM $3 AND
    lv.created_user IS NOT DISTINCT FROM $4 AND
    lv.created_date IS NOT DISTINCT FROM $5 AND
    lv.last_edited_user IS NOT DISTINCT FROM $6 AND
    lv.last_edited_date IS NOT DISTINCT FROM $7 AND
    lv.creationdate IS NOT DISTINCT FROM $8 AND
    lv.creator IS NOT DISTINCT FROM $9 AND
    lv.editdate IS NOT DISTINCT FROM $10 AND
    lv.editor IS NOT DISTINCT FROM $11
  )
RETURNING *;

-- Example usage: EXECUTE insert_fieldscoutinglog_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.HabitatRelate
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.habitatrelate_habitatrelate_habitattype_2e81cf2f550e400783cf284f3cec3953_enum AS ENUM (
  'orchard',
  'row_crops',
  'vine_crops',
  'ag_grasses_or_grain',
  'pasture',
  'irrigation_standpipe',
  'ditch',
  'pond',
  'sump',
  'drain',
  'dairy_lagoon',
  'wastewater_treatment',
  'trough',
  'depression',
  'gutter',
  'rain_gutter',
  'culvert',
  'utility',
  'catch_basin',
  'stream_or_creek',
  'slough',
  'river',
  'marsh_or_wetlands',
  'containers',
  'watering_bowl',
  'plant_saucer',
  'yard_drain',
  'plant_axil',
  'treehole',
  'fountain_or_water_feature',
  'bird_bath',
  'misc_water_accumulation',
  'tarp_or_cover',
  'swimming_pool',
  'aboveground_pool',
  'kid_pool',
  'hot_tub',
  'appliance',
  'tires',
  'flooded_structure',
  'low_point'
);

CREATE TABLE fieldseeker.habitatrelate (
  objectid BIGSERIAL NOT NULL,
  foreign_id UUID,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  habitattype fieldseeker.habitatrelate_habitatrelate_habitattype_2e81cf2f550e400783cf284f3cec3953_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.habitatrelate.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.habitatrelate.habitattype IS 'Habitat Type';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_habitatrelate_versioned(bigint, uuid, uuid, varchar, timestamp, varchar, timestamp, fieldseeker.habitatrelate_habitatrelate_habitattype_2e81cf2f550e400783cf284f3cec3953_enum, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.habitatrelate
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.habitatrelate
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.habitatrelate (
  objectid, foreign_id, globalid, created_user, created_date, last_edited_user, last_edited_date, habitattype, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.foreign_id IS NOT DISTINCT FROM $2 AND
    lv.globalid IS NOT DISTINCT FROM $3 AND
    lv.created_user IS NOT DISTINCT FROM $4 AND
    lv.created_date IS NOT DISTINCT FROM $5 AND
    lv.last_edited_user IS NOT DISTINCT FROM $6 AND
    lv.last_edited_date IS NOT DISTINCT FROM $7 AND
    lv.habitattype IS NOT DISTINCT FROM $8 AND
    lv.creationdate IS NOT DISTINCT FROM $9 AND
    lv.creator IS NOT DISTINCT FROM $10 AND
    lv.editdate IS NOT DISTINCT FROM $11 AND
    lv.editor IS NOT DISTINCT FROM $12
  )
RETURNING *;

-- Example usage: EXECUTE insert_habitatrelate_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.InspectionSample
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.inspectionsample_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TABLE fieldseeker.inspectionsample (
  objectid BIGSERIAL NOT NULL,
  insp_id UUID,
  sampleid VARCHAR(25),
  processed fieldseeker.inspectionsample_notinuit_f_enum,
  idbytech VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.inspectionsample.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.inspectionsample.sampleid IS 'Sample ID';
COMMENT ON COLUMN fieldseeker.inspectionsample.processed IS 'Processed';
COMMENT ON COLUMN fieldseeker.inspectionsample.idbytech IS 'Tech Identifying Species in Lab';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_inspectionsample_versioned(bigint, uuid, varchar, fieldseeker.inspectionsample_notinuit_f_enum, varchar, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.inspectionsample
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.inspectionsample
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.inspectionsample (
  objectid, insp_id, sampleid, processed, idbytech, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.insp_id IS NOT DISTINCT FROM $2 AND
    lv.sampleid IS NOT DISTINCT FROM $3 AND
    lv.processed IS NOT DISTINCT FROM $4 AND
    lv.idbytech IS NOT DISTINCT FROM $5 AND
    lv.globalid IS NOT DISTINCT FROM $6 AND
    lv.created_user IS NOT DISTINCT FROM $7 AND
    lv.created_date IS NOT DISTINCT FROM $8 AND
    lv.last_edited_user IS NOT DISTINCT FROM $9 AND
    lv.last_edited_date IS NOT DISTINCT FROM $10 AND
    lv.creationdate IS NOT DISTINCT FROM $11 AND
    lv.creator IS NOT DISTINCT FROM $12 AND
    lv.editdate IS NOT DISTINCT FROM $13 AND
    lv.editor IS NOT DISTINCT FROM $14
  )
RETURNING *;

-- Example usage: EXECUTE insert_inspectionsample_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.InspectionSampleDetail
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.inspectionsampledetail_mosquitofieldspecies_enum AS ENUM (
  'Aedes',
  'Culex'
);

CREATE TYPE fieldseeker.inspectionsampledetail_mosquitodominantstage_enum AS ENUM (
  '1',
  '2',
  '3',
  '4',
  '1-2',
  '3-4'
);

CREATE TYPE fieldseeker.inspectionsampledetail_mosquitoadultactivity_enum AS ENUM (
  'None',
  'Light',
  'Moderate',
  'Intense'
);

CREATE TABLE fieldseeker.inspectionsampledetail (
  objectid BIGSERIAL NOT NULL,
  inspsample_id UUID,
  fieldspecies fieldseeker.inspectionsampledetail_mosquitofieldspecies_enum,
  flarvcount SMALLINT,
  fpupcount SMALLINT,
  feggcount SMALLINT,
  flstages VARCHAR(25),
  fdomstage fieldseeker.inspectionsampledetail_mosquitodominantstage_enum,
  fadultact fieldseeker.inspectionsampledetail_mosquitoadultactivity_enum,
  labspecies VARCHAR(50),
  llarvcount SMALLINT,
  lpupcount SMALLINT,
  leggcount SMALLINT,
  ldomstage fieldseeker.inspectionsampledetail_mosquitodominantstage_enum,
  comments VARCHAR(250),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  processed SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.inspectionsampledetail.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.inspectionsampledetail.fieldspecies IS 'Field Species';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.flarvcount IS 'Field Larva Count';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.fpupcount IS 'Field Pupa Count';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.feggcount IS 'Field Egg Count';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.flstages IS 'Field Larval Stages';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.fdomstage IS 'Field Dominant Stage';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.fadultact IS 'Field Adult Activity';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.labspecies IS 'Lab Species';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.llarvcount IS 'Lab Larva Count';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.lpupcount IS 'Lab Pupa Count';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.leggcount IS 'Lab Egg Count';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.ldomstage IS 'Lab Dominant Stage';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.comments IS 'Comments';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_inspectionsampledetail_versioned(bigint, uuid, fieldseeker.inspectionsampledetail_mosquitofieldspecies_enum, smallint, smallint, smallint, varchar, fieldseeker.inspectionsampledetail_mosquitodominantstage_enum, fieldseeker.inspectionsampledetail_mosquitoadultactivity_enum, varchar, smallint, smallint, smallint, fieldseeker.inspectionsampledetail_mosquitodominantstage_enum, varchar, uuid, varchar, timestamp, varchar, timestamp, smallint, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.inspectionsampledetail
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.inspectionsampledetail
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.inspectionsampledetail (
  objectid, inspsample_id, fieldspecies, flarvcount, fpupcount, feggcount, flstages, fdomstage, fadultact, labspecies, llarvcount, lpupcount, leggcount, ldomstage, comments, globalid, created_user, created_date, last_edited_user, last_edited_date, processed, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.inspsample_id IS NOT DISTINCT FROM $2 AND
    lv.fieldspecies IS NOT DISTINCT FROM $3 AND
    lv.flarvcount IS NOT DISTINCT FROM $4 AND
    lv.fpupcount IS NOT DISTINCT FROM $5 AND
    lv.feggcount IS NOT DISTINCT FROM $6 AND
    lv.flstages IS NOT DISTINCT FROM $7 AND
    lv.fdomstage IS NOT DISTINCT FROM $8 AND
    lv.fadultact IS NOT DISTINCT FROM $9 AND
    lv.labspecies IS NOT DISTINCT FROM $10 AND
    lv.llarvcount IS NOT DISTINCT FROM $11 AND
    lv.lpupcount IS NOT DISTINCT FROM $12 AND
    lv.leggcount IS NOT DISTINCT FROM $13 AND
    lv.ldomstage IS NOT DISTINCT FROM $14 AND
    lv.comments IS NOT DISTINCT FROM $15 AND
    lv.globalid IS NOT DISTINCT FROM $16 AND
    lv.created_user IS NOT DISTINCT FROM $17 AND
    lv.created_date IS NOT DISTINCT FROM $18 AND
    lv.last_edited_user IS NOT DISTINCT FROM $19 AND
    lv.last_edited_date IS NOT DISTINCT FROM $20 AND
    lv.processed IS NOT DISTINCT FROM $21 AND
    lv.creationdate IS NOT DISTINCT FROM $22 AND
    lv.creator IS NOT DISTINCT FROM $23 AND
    lv.editdate IS NOT DISTINCT FROM $24 AND
    lv.editor IS NOT DISTINCT FROM $25
  )
RETURNING *;

-- Example usage: EXECUTE insert_inspectionsampledetail_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.LineLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.linelocation_linelocation_habitat_fc51bdc4f1954df58206d69ce14182f3_enum AS ENUM (
  'orchard',
  'row_crops',
  'vine_crops',
  'ag_grasses_or_grain',
  'pasture',
  'irrigation_standpipe',
  'ditch',
  'pond',
  'sump',
  'drain',
  'dairy_lagoon',
  'wastewater_treatment',
  'trough',
  'depression',
  'gutter',
  'rain_gutter',
  'culvert',
  'utility',
  'catch_basin',
  'stream_or_creek',
  'slough',
  'river',
  'marsh_or_wetland',
  'containers',
  'watering_bowl',
  'plant_saucer',
  'yard_drain',
  'plant_axil',
  'treehole',
  'fountain_or_water_feature',
  'bird_bath',
  'misc_water_accumulation',
  'tarp_or_cover',
  'swimming_pool',
  'aboveground_pool',
  'kid_pool',
  'hot_tub',
  'appliance',
  'tires',
  'flooded_structure',
  'low_point'
);

CREATE TYPE fieldseeker.linelocation_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TYPE fieldseeker.linelocation_linelocation_usetype_2aeca2e60d2f455c86fc34895dc80a02_enum AS ENUM (
  'residential',
  'commercial',
  'industrial',
  'agricultural',
  'mixed_use',
  'public_domain',
  'natural',
  'municipal'
);

CREATE TYPE fieldseeker.linelocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.linelocation_locationsymbology_enum AS ENUM (
  'ACTION',
  'INACTIVE',
  'NONE'
);

CREATE TYPE fieldseeker.linelocation_linelocation_waterorigin_84723d92_306a_46f4_8ef1_69b55a916008_enum AS ENUM (
  'flood_irrigation',
  'furrow_irrigation',
  'drip_irrigation',
  'sprinkler_irrigation',
  'wastewater_irrigation',
  'irrigation_runoff',
  'rainwater_accumulation',
  'leak',
  'seepage',
  'stored_water',
  'wastewater_system',
  'permanent_natural_water',
  'temporary_natural_water',
  'recreational_or_ornamental_water',
  'water_conveyance'
);

CREATE TABLE fieldseeker.linelocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.linelocation_linelocation_habitat_fc51bdc4f1954df58206d69ce14182f3_enum,
  priority fieldseeker.linelocation_locationpriority_enum,
  usetype fieldseeker.linelocation_linelocation_usetype_2aeca2e60d2f455c86fc34895dc80a02_enum,
  active fieldseeker.linelocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology fieldseeker.linelocation_locationsymbology_enum,
  externalid VARCHAR(50),
  acres DOUBLE PRECISION,
  nextactiondatescheduled TIMESTAMP,
  larvinspectinterval SMALLINT,
  length_ft DOUBLE PRECISION,
  width_ft DOUBLE PRECISION,
  zone2 VARCHAR(25),
  locationnumber INTEGER,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  lastinspectdate TIMESTAMP,
  lastinspectbreeding VARCHAR(25),
  lastinspectavglarvae DOUBLE PRECISION,
  lastinspectavgpupae DOUBLE PRECISION,
  lastinspectlstages VARCHAR(25),
  lastinspectactiontaken VARCHAR(50),
  lastinspectfieldspecies VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  hectares DOUBLE PRECISION,
  lastinspectactivity VARCHAR(25),
  lasttreatactivity VARCHAR(25),
  length_meters DOUBLE PRECISION,
  width_meters DOUBLE PRECISION,
  lastinspectconditions VARCHAR(250),
  waterorigin fieldseeker.linelocation_linelocation_waterorigin_84723d92_306a_46f4_8ef1_69b55a916008_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.linelocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.linelocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.linelocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.linelocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.linelocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.linelocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.linelocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.linelocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.linelocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.linelocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.linelocation.symbology IS 'Symbology';
COMMENT ON COLUMN fieldseeker.linelocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.linelocation.acres IS 'Acres';
COMMENT ON COLUMN fieldseeker.linelocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.linelocation.larvinspectinterval IS 'Larval Inspection Interval';
COMMENT ON COLUMN fieldseeker.linelocation.length_ft IS 'Length';
COMMENT ON COLUMN fieldseeker.linelocation.width_ft IS 'Width';
COMMENT ON COLUMN fieldseeker.linelocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectdate IS 'Last Inspection Date';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectbreeding IS 'Last Inspection Breeding';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectavglarvae IS 'Last Inspection Average Larvae';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectavgpupae IS 'Last Inspection Average Pupae';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectlstages IS 'Last Inspection Larval Stages';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectactiontaken IS 'Last Inspection Action';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectfieldspecies IS 'Last Inspection Field Species';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatdate IS 'Last Treatment Date';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatproduct IS 'Last Treatment Product';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatqty IS 'Last Treatment Quantity';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatqtyunit IS 'Last Treatment Quantity Unit';
COMMENT ON COLUMN fieldseeker.linelocation.hectares IS 'Hectares';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectactivity IS 'Last Inspection Activity';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatactivity IS 'Last Treatment Activity';
COMMENT ON COLUMN fieldseeker.linelocation.length_meters IS 'Length Meters';
COMMENT ON COLUMN fieldseeker.linelocation.width_meters IS 'Width Meters';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectconditions IS 'Last Inspection Conditions';
COMMENT ON COLUMN fieldseeker.linelocation.waterorigin IS 'Water Origin';
COMMENT ON COLUMN fieldseeker.linelocation.jurisdiction IS 'Jurisdiction';

-- Field active has default value: 1

-- Field symbology has default value: 'NONE'

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_linelocation_versioned(bigint, varchar, varchar, fieldseeker.linelocation_linelocation_habitat_fc51bdc4f1954df58206d69ce14182f3_enum, fieldseeker.linelocation_locationpriority_enum, fieldseeker.linelocation_linelocation_usetype_2aeca2e60d2f455c86fc34895dc80a02_enum, fieldseeker.linelocation_notinuit_f_enum, varchar, varchar, varchar, fieldseeker.linelocation_locationsymbology_enum, varchar, double precision, timestamp, smallint, double precision, double precision, varchar, integer, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, double precision, double precision, varchar, varchar, varchar, timestamp, varchar, double precision, varchar, double precision, varchar, varchar, double precision, double precision, varchar, fieldseeker.linelocation_linelocation_waterorigin_84723d92_306a_46f4_8ef1_69b55a916008_enum, timestamp, varchar, timestamp, varchar, varchar, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.linelocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.linelocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.linelocation (
  objectid, name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, acres, nextactiondatescheduled, larvinspectinterval, length_ft, width_ft, zone2, locationnumber, globalid, created_user, created_date, last_edited_user, last_edited_date, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, hectares, lastinspectactivity, lasttreatactivity, length_meters, width_meters, lastinspectconditions, waterorigin, creationdate, creator, editdate, editor, jurisdiction, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48,
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
    lv.symbology IS NOT DISTINCT FROM $11 AND
    lv.externalid IS NOT DISTINCT FROM $12 AND
    lv.acres IS NOT DISTINCT FROM $13 AND
    lv.nextactiondatescheduled IS NOT DISTINCT FROM $14 AND
    lv.larvinspectinterval IS NOT DISTINCT FROM $15 AND
    lv.length_ft IS NOT DISTINCT FROM $16 AND
    lv.width_ft IS NOT DISTINCT FROM $17 AND
    lv.zone2 IS NOT DISTINCT FROM $18 AND
    lv.locationnumber IS NOT DISTINCT FROM $19 AND
    lv.globalid IS NOT DISTINCT FROM $20 AND
    lv.created_user IS NOT DISTINCT FROM $21 AND
    lv.created_date IS NOT DISTINCT FROM $22 AND
    lv.last_edited_user IS NOT DISTINCT FROM $23 AND
    lv.last_edited_date IS NOT DISTINCT FROM $24 AND
    lv.lastinspectdate IS NOT DISTINCT FROM $25 AND
    lv.lastinspectbreeding IS NOT DISTINCT FROM $26 AND
    lv.lastinspectavglarvae IS NOT DISTINCT FROM $27 AND
    lv.lastinspectavgpupae IS NOT DISTINCT FROM $28 AND
    lv.lastinspectlstages IS NOT DISTINCT FROM $29 AND
    lv.lastinspectactiontaken IS NOT DISTINCT FROM $30 AND
    lv.lastinspectfieldspecies IS NOT DISTINCT FROM $31 AND
    lv.lasttreatdate IS NOT DISTINCT FROM $32 AND
    lv.lasttreatproduct IS NOT DISTINCT FROM $33 AND
    lv.lasttreatqty IS NOT DISTINCT FROM $34 AND
    lv.lasttreatqtyunit IS NOT DISTINCT FROM $35 AND
    lv.hectares IS NOT DISTINCT FROM $36 AND
    lv.lastinspectactivity IS NOT DISTINCT FROM $37 AND
    lv.lasttreatactivity IS NOT DISTINCT FROM $38 AND
    lv.length_meters IS NOT DISTINCT FROM $39 AND
    lv.width_meters IS NOT DISTINCT FROM $40 AND
    lv.lastinspectconditions IS NOT DISTINCT FROM $41 AND
    lv.waterorigin IS NOT DISTINCT FROM $42 AND
    lv.creationdate IS NOT DISTINCT FROM $43 AND
    lv.creator IS NOT DISTINCT FROM $44 AND
    lv.editdate IS NOT DISTINCT FROM $45 AND
    lv.editor IS NOT DISTINCT FROM $46 AND
    lv.jurisdiction IS NOT DISTINCT FROM $47 AND
    lv.shape__length IS NOT DISTINCT FROM $48
  )
RETURNING *;

-- Example usage: EXECUTE insert_linelocation_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.LocationTracking
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.locationtracking (
  objectid BIGSERIAL NOT NULL,
  accuracy DOUBLE PRECISION,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  globalid UUID,
  fieldtech VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.locationtracking.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.locationtracking.accuracy IS 'Accuracy(m)';
COMMENT ON COLUMN fieldseeker.locationtracking.fieldtech IS 'Field Tech';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_locationtracking_versioned(bigint, double precision, varchar, timestamp, varchar, timestamp, uuid, varchar, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.locationtracking
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.locationtracking
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.locationtracking (
  objectid, accuracy, created_user, created_date, last_edited_user, last_edited_date, globalid, fieldtech, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.accuracy IS NOT DISTINCT FROM $2 AND
    lv.created_user IS NOT DISTINCT FROM $3 AND
    lv.created_date IS NOT DISTINCT FROM $4 AND
    lv.last_edited_user IS NOT DISTINCT FROM $5 AND
    lv.last_edited_date IS NOT DISTINCT FROM $6 AND
    lv.globalid IS NOT DISTINCT FROM $7 AND
    lv.fieldtech IS NOT DISTINCT FROM $8 AND
    lv.creationdate IS NOT DISTINCT FROM $9 AND
    lv.creator IS NOT DISTINCT FROM $10 AND
    lv.editdate IS NOT DISTINCT FROM $11 AND
    lv.editor IS NOT DISTINCT FROM $12
  )
RETURNING *;

-- Example usage: EXECUTE insert_locationtracking_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.MosquitoInspection
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoinspection_adminaction_b74ae1bb_c98b_40f6_8cfa_40e4fd16c270_enum AS ENUM (
  'yes',
  'no'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitobreeding_enum AS ENUM (
  'None',
  'Light',
  'Moderate',
  'Intense'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoinspection_domstage_b7a6c36bccde49a292020de4812cf5ae_enum AS ENUM (
  '1',
  '2',
  '3',
  '4',
  '5',
  '1-2',
  '2-3',
  '3-4'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoinspection_actiontaken_252243d69b0b44ddbdc229c04ec3a8d5_enum AS ENUM (
  'Treatment',
  'Mechanical or Biological Treatment',
  'Resident Schedule Request',
  'Administrative'
);

CREATE TYPE fieldseeker.mosquitoinspection_notinuiwinddirection_enum AS ENUM (
  'N',
  'NE',
  'E',
  'SE',
  'S',
  'SW',
  'W',
  'NW'
);

CREATE TYPE fieldseeker.mosquitoinspection_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitofieldspecies_enum AS ENUM (
  'Aedes',
  'Culex'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoactivity_enum AS ENUM (
  'Routine inspection',
  'Pre-treatment',
  'Maintenance',
  'ULV',
  'BARRIER',
  'LOGIN',
  'TREATSD',
  'SD',
  'SITEVISIT',
  'ONLINE',
  'SYNC',
  'CREATESR',
  'LC',
  'ACCEPTSR',
  'POINT',
  'DOWNLOAD',
  'COMPLETESR',
  'POLYGON',
  'TRAP',
  'SAMPLE',
  'QA',
  'PTA',
  'FIELDSCOUTING',
  'OFFLINE',
  'LINE',
  'TRAPLOCATION',
  'SAMPLELOCATION',
  'LCLOCATION'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoadultactivity_enum AS ENUM (
  'None',
  'Light',
  'Moderate',
  'Intense'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoinspection_sitecond_db7350bc_81e5_401e_858f_cd3e5e5d8a34_enum AS ENUM (
  'Dry',
  'Flowing',
  'Maintained',
  'Unmaintained',
  'High Organic',
  'Unknown',
  'Stagnant',
  'Needs Monitoring',
  'Drying Out',
  'Appears Vacant',
  'Entry Denied',
  'Pool Removed',
  'False Pool'
);

CREATE TABLE fieldseeker.mosquitoinspection (
  objectid BIGSERIAL NOT NULL,
  numdips SMALLINT,
  activity fieldseeker.mosquitoinspection_mosquitoactivity_enum,
  breeding fieldseeker.mosquitoinspection_mosquitobreeding_enum,
  totlarvae SMALLINT,
  totpupae SMALLINT,
  eggs SMALLINT,
  posdips SMALLINT,
  adultact fieldseeker.mosquitoinspection_mosquitoadultactivity_enum,
  lstages VARCHAR(25),
  domstage fieldseeker.mosquitoinspection_mosquitoinspection_domstage_b7a6c36bccde49a292020de4812cf5ae_enum,
  actiontaken fieldseeker.mosquitoinspection_mosquitoinspection_actiontaken_252243d69b0b44ddbdc229c04ec3a8d5_enum,
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  raingauge DOUBLE PRECISION,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  winddir fieldseeker.mosquitoinspection_notinuiwinddirection_enum,
  avglarvae DOUBLE PRECISION,
  avgpupae DOUBLE PRECISION,
  reviewed fieldseeker.mosquitoinspection_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  personalcontact fieldseeker.mosquitoinspection_notinuit_f_enum,
  tirecount SMALLINT,
  cbcount SMALLINT,
  containercount SMALLINT,
  fieldspecies fieldseeker.mosquitoinspection_mosquitofieldspecies_enum,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  linelocid UUID,
  pointlocid UUID,
  polygonlocid UUID,
  srid UUID,
  fieldtech VARCHAR(25),
  larvaepresent fieldseeker.mosquitoinspection_notinuit_f_enum,
  pupaepresent fieldseeker.mosquitoinspection_notinuit_f_enum,
  sdid UUID,
  sitecond fieldseeker.mosquitoinspection_mosquitoinspection_sitecond_db7350bc_81e5_401e_858f_cd3e5e5d8a34_enum,
  positivecontainercount SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  visualmonitoring fieldseeker.mosquitoinspection_notinuit_f_enum,
  vmcomments VARCHAR(250),
  adminaction fieldseeker.mosquitoinspection_mosquitoinspection_adminaction_b74ae1bb_c98b_40f6_8cfa_40e4fd16c270_enum,
  ptaid UUID,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.mosquitoinspection.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.mosquitoinspection.numdips IS '# Dips';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.activity IS 'Activity';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.breeding IS 'Breeding';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.totlarvae IS 'Total Larvae';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.totpupae IS 'Total Pupae';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.eggs IS 'Eggs';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.posdips IS 'Positive Dips';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.adultact IS 'Adult Activity';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.lstages IS 'Larval Stages';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.domstage IS 'Dominant Stage';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.actiontaken IS 'Action';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avetemp IS 'Average Temperature';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.windspeed IS 'Wind Speed';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.raingauge IS 'Rain Gauge';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.winddir IS 'Wind Direction';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avglarvae IS 'Average Larvae';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avgpupae IS 'Average Pupae';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.recordstatus IS 'RecordStatus';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.personalcontact IS 'Personal Contact';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.tirecount IS 'Tire Count';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.cbcount IS 'Catch Basin Count';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.containercount IS 'Container Count';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.fieldspecies IS 'Field Species';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.fieldtech IS 'Field Tech';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.larvaepresent IS 'Larvae Present';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.pupaepresent IS 'Pupae Present';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.sdid IS 'Storm Drain ID';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.sitecond IS 'Conditions';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.positivecontainercount IS 'Positive Container Count';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.jurisdiction IS 'Jurisdiction';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.visualmonitoring IS 'Visual Monitoring';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.vmcomments IS 'VM Comments';

-- Field adminaction has default value: 'no'

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.PointLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.pointlocation_pointlocation_habitat_b4d8135a_4979_49c8_8bb3_67ec7230e661_enum AS ENUM (
  'orchard',
  'row_crops',
  'vine_crops',
  'ag_grasses_or_grain',
  'pasture',
  'irrigation_standpipe',
  'ditch',
  'pond',
  'sump',
  'drain',
  'dairy_lagoon',
  'wastewater_treatment',
  'trough',
  'depression',
  'gutter',
  'rain_gutter',
  'culvert',
  'utility',
  'catch_basin',
  'stream_or_creek',
  'slough',
  'river',
  'marsh_or_wetlands',
  'containers',
  'watering_bowl',
  'plant_saucer',
  'yard_drain',
  'plant_axil',
  'treehole',
  'fountain_or_water_feature',
  'bird_bath',
  'misc_water_accumulation',
  'tarp_or_cover',
  'swimming_pool',
  'aboveground_pool',
  'kid_pool',
  'hot_tub',
  'appliance',
  'tires',
  'flooded_structure',
  'low_point'
);

CREATE TYPE fieldseeker.pointlocation_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TYPE fieldseeker.pointlocation_pointlocation_usetype_58d62d18ef4f47fc8cb9874df867f89e_enum AS ENUM (
  'residential',
  'commercial',
  'agricultural',
  'industrial',
  'mixed_use',
  'public_domain',
  'natural',
  'municipal'
);

CREATE TYPE fieldseeker.pointlocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.pointlocation_locationsymbology_enum AS ENUM (
  'ACTION',
  'INACTIVE',
  'NONE'
);

CREATE TYPE fieldseeker.pointlocation_pointlocation_waterorigin_197b22bf_f3eb_4dad_8899_986460f6ea97_enum AS ENUM (
  'flood_irrigation',
  'furrow_irrigation',
  'drip_irrigation',
  'sprinkler_irrigation',
  'wastewater_irrigation',
  'irrigation_runoff',
  'stormwater_or_municipal_runoff',
  'industrial_runoff',
  'rainwater_accumulation',
  'leak',
  'seepage',
  'stored_water',
  'wastewater_system',
  'permanent_natural_water',
  'temporary_natural_water',
  'recreational_or_orenamental_water',
  'water_conveyance'
);

CREATE TYPE fieldseeker.pointlocation_pointlocation_assignedtech_9393a162_2474_429d_85be_daa44e4c091f_enum AS ENUM (
  'Bryan Feguson',
  'Rick Alverez',
  'Alysia Davis',
  'Bryan Ruiz',
  'Kory Wilson',
  'Adrian Sifuentes',
  'Marco Martinez',
  'Carlos Rodriguez',
  'Landon McGill',
  'Ted McGill',
  'Mario Sanchez',
  'Jorge Perez',
  'Arturo Garcia-Trejo',
  'Lisa Salgado',
  'Lawrence Guzman',
  'Tricia Snowden',
  'Ryan Spratt',
  'Andrea Troupin',
  'Mark Nakata',
  'Pablo Ortega',
  'Benjamin Sperry',
  'Fatima Hidalgo',
  'Zackery Barragan',
  'Yajaira Godinez',
  'Jake Maldonado',
  'Rafael Ramirez',
  'Carlos Palacios',
  'Aaron Fredrick',
  'Josh Malone',
  'Alec Caposella',
  'Laura Ramos'
);

CREATE TYPE fieldseeker.pointlocation_pointlocation_deactivate_reason_dd303085_b33c_4894_8c47_fa847dd9d7c5_enum AS ENUM (
  'Source Removed',
  'Pool Maintained',
  'Source Screened',
  'Crop Change',
  'Low or No Mosquito Activity',
  'Consistent Fish Presence'
);

CREATE TABLE fieldseeker.pointlocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.pointlocation_pointlocation_habitat_b4d8135a_4979_49c8_8bb3_67ec7230e661_enum,
  priority fieldseeker.pointlocation_locationpriority_enum,
  usetype fieldseeker.pointlocation_pointlocation_usetype_58d62d18ef4f47fc8cb9874df867f89e_enum,
  active fieldseeker.pointlocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology fieldseeker.pointlocation_locationsymbology_enum,
  externalid VARCHAR(50),
  nextactiondatescheduled TIMESTAMP,
  larvinspectinterval SMALLINT,
  zone2 VARCHAR(25),
  locationnumber INTEGER,
  globalid UUID,
  stype VARCHAR(3),
  lastinspectdate TIMESTAMP,
  lastinspectbreeding VARCHAR(25),
  lastinspectavglarvae DOUBLE PRECISION,
  lastinspectavgpupae DOUBLE PRECISION,
  lastinspectlstages VARCHAR(25),
  lastinspectactiontaken VARCHAR(50),
  lastinspectfieldspecies VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  lastinspectactivity VARCHAR(25),
  lasttreatactivity VARCHAR(25),
  lastinspectconditions VARCHAR(250),
  waterorigin fieldseeker.pointlocation_pointlocation_waterorigin_197b22bf_f3eb_4dad_8899_986460f6ea97_enum,
  x DOUBLE PRECISION,
  y DOUBLE PRECISION,
  assignedtech fieldseeker.pointlocation_pointlocation_assignedtech_9393a162_2474_429d_85be_daa44e4c091f_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  deactivate_reason fieldseeker.pointlocation_pointlocation_deactivate_reason_dd303085_b33c_4894_8c47_fa847dd9d7c5_enum,
  scalarpriority INTEGER,
  sourcestatus VARCHAR(255),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.pointlocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.pointlocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.pointlocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.pointlocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.pointlocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.pointlocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.pointlocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.pointlocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.pointlocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.pointlocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.pointlocation.symbology IS 'Symbology';
COMMENT ON COLUMN fieldseeker.pointlocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.pointlocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.pointlocation.larvinspectinterval IS 'Larval Inspection Interval';
COMMENT ON COLUMN fieldseeker.pointlocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.pointlocation.stype IS 'SourceType';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectdate IS 'Last Inspection Date';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectbreeding IS 'Last Inspection Breeding';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectavglarvae IS 'Last Inspection Average Larvae';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectavgpupae IS 'Last Inspection Average Pupae';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectlstages IS 'Last Inspection Larval Stages';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectactiontaken IS 'Last Inspection Action';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectfieldspecies IS 'Last Inspection Field Species';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatdate IS 'Last Treatment Date';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatproduct IS 'Last Treatment Product';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatqty IS 'Last Treatment Quantity';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatqtyunit IS 'Last Treatment Quantity Unit';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectactivity IS 'Last Inspection Activity';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatactivity IS 'Last Treatment Activity';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectconditions IS 'Last Inspection Conditions';
COMMENT ON COLUMN fieldseeker.pointlocation.waterorigin IS 'Water Origin';
COMMENT ON COLUMN fieldseeker.pointlocation.assignedtech IS 'Assigned Tech';
COMMENT ON COLUMN fieldseeker.pointlocation.jurisdiction IS 'Jurisdiction';
COMMENT ON COLUMN fieldseeker.pointlocation.deactivate_reason IS 'Reason for Deactivation';

-- Field active has default value: 1

-- Field symbology has default value: 'NONE'

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_pointlocation_versioned(bigint, varchar, varchar, fieldseeker.pointlocation_pointlocation_habitat_b4d8135a_4979_49c8_8bb3_67ec7230e661_enum, fieldseeker.pointlocation_locationpriority_enum, fieldseeker.pointlocation_pointlocation_usetype_58d62d18ef4f47fc8cb9874df867f89e_enum, fieldseeker.pointlocation_notinuit_f_enum, varchar, varchar, varchar, fieldseeker.pointlocation_locationsymbology_enum, varchar, timestamp, smallint, varchar, integer, uuid, varchar, timestamp, varchar, double precision, double precision, varchar, varchar, varchar, timestamp, varchar, double precision, varchar, varchar, varchar, varchar, fieldseeker.pointlocation_pointlocation_waterorigin_197b22bf_f3eb_4dad_8899_986460f6ea97_enum, double precision, double precision, fieldseeker.pointlocation_pointlocation_assignedtech_9393a162_2474_429d_85be_daa44e4c091f_enum, timestamp, varchar, timestamp, varchar, varchar, fieldseeker.pointlocation_pointlocation_deactivate_reason_dd303085_b33c_4894_8c47_fa847dd9d7c5_enum, integer, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.pointlocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.pointlocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.pointlocation (
  objectid, name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, nextactiondatescheduled, larvinspectinterval, zone2, locationnumber, globalid, stype, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, lastinspectactivity, lasttreatactivity, lastinspectconditions, waterorigin, x, y, assignedtech, creationdate, creator, editdate, editor, jurisdiction, deactivate_reason, scalarpriority, sourcestatus,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44,
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
    lv.symbology IS NOT DISTINCT FROM $11 AND
    lv.externalid IS NOT DISTINCT FROM $12 AND
    lv.nextactiondatescheduled IS NOT DISTINCT FROM $13 AND
    lv.larvinspectinterval IS NOT DISTINCT FROM $14 AND
    lv.zone2 IS NOT DISTINCT FROM $15 AND
    lv.locationnumber IS NOT DISTINCT FROM $16 AND
    lv.globalid IS NOT DISTINCT FROM $17 AND
    lv.stype IS NOT DISTINCT FROM $18 AND
    lv.lastinspectdate IS NOT DISTINCT FROM $19 AND
    lv.lastinspectbreeding IS NOT DISTINCT FROM $20 AND
    lv.lastinspectavglarvae IS NOT DISTINCT FROM $21 AND
    lv.lastinspectavgpupae IS NOT DISTINCT FROM $22 AND
    lv.lastinspectlstages IS NOT DISTINCT FROM $23 AND
    lv.lastinspectactiontaken IS NOT DISTINCT FROM $24 AND
    lv.lastinspectfieldspecies IS NOT DISTINCT FROM $25 AND
    lv.lasttreatdate IS NOT DISTINCT FROM $26 AND
    lv.lasttreatproduct IS NOT DISTINCT FROM $27 AND
    lv.lasttreatqty IS NOT DISTINCT FROM $28 AND
    lv.lasttreatqtyunit IS NOT DISTINCT FROM $29 AND
    lv.lastinspectactivity IS NOT DISTINCT FROM $30 AND
    lv.lasttreatactivity IS NOT DISTINCT FROM $31 AND
    lv.lastinspectconditions IS NOT DISTINCT FROM $32 AND
    lv.waterorigin IS NOT DISTINCT FROM $33 AND
    lv.x IS NOT DISTINCT FROM $34 AND
    lv.y IS NOT DISTINCT FROM $35 AND
    lv.assignedtech IS NOT DISTINCT FROM $36 AND
    lv.creationdate IS NOT DISTINCT FROM $37 AND
    lv.creator IS NOT DISTINCT FROM $38 AND
    lv.editdate IS NOT DISTINCT FROM $39 AND
    lv.editor IS NOT DISTINCT FROM $40 AND
    lv.jurisdiction IS NOT DISTINCT FROM $41 AND
    lv.deactivate_reason IS NOT DISTINCT FROM $42 AND
    lv.scalarpriority IS NOT DISTINCT FROM $43 AND
    lv.sourcestatus IS NOT DISTINCT FROM $44
  )
RETURNING *;

-- Example usage: EXECUTE insert_pointlocation_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.PolygonLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.polygonlocation_polygonlocation_waterorigin_e9018e92_5f47_4ff9_8a7c_b818d848dc7a_enum AS ENUM (
  'flood_irrigation',
  'furrow_irrigation',
  'drip_irritation',
  'sprinkler_irrigation',
  'wastewater_irrigation',
  'irrigation_runoff',
  'rainwater_accumulation',
  'leak',
  'seepage',
  'stored_water',
  'wastewater_system',
  'permanent_natural_water',
  'temporary_natural_water',
  'recreational_or_ornamental_water',
  'water_conveyance'
);

CREATE TYPE fieldseeker.polygonlocation_polygonlocation_habitat_45e9dde79ac84d959df8b65ba7d5dafd_enum AS ENUM (
  'orchard',
  'row_crops',
  'vine_crops',
  'ag_grasses_or_grain',
  'pasture',
  'irrigation_standpipe',
  'ditch',
  'dairy_lagoon',
  'wastewater_treatment',
  'trough',
  'depression',
  'gutter',
  'rain_gutter',
  'culvert',
  'utility',
  'catch_basin',
  'stream_or_creek',
  'slough',
  'river',
  'marsh_or_wetland',
  'containers',
  'watering_bowl',
  'plant_saucer',
  'yard_drain',
  'plant_axil',
  'treehole',
  'fountain_or_water_feature',
  'bird_bath',
  'misc_water_accumulation',
  'tarp_or_cover',
  'swimming_pool',
  'aboveground_pool',
  'kid_pool',
  'hot_tub',
  'appliance',
  'tires',
  'flooded_structure',
  'low_point',
  'unknown'
);

CREATE TYPE fieldseeker.polygonlocation_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TYPE fieldseeker.polygonlocation_polygonlocation_usetype_e546154cb9544b9aa8e7b13e8e258b27_enum AS ENUM (
  'residential',
  'commercial',
  'industrial',
  'agricultural',
  'mixed_use',
  'public_domain',
  'natural',
  'municipal'
);

CREATE TYPE fieldseeker.polygonlocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.polygonlocation_locationsymbology_enum AS ENUM (
  'ACTION',
  'INACTIVE',
  'NONE'
);

CREATE TABLE fieldseeker.polygonlocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.polygonlocation_polygonlocation_habitat_45e9dde79ac84d959df8b65ba7d5dafd_enum,
  priority fieldseeker.polygonlocation_locationpriority_enum,
  usetype fieldseeker.polygonlocation_polygonlocation_usetype_e546154cb9544b9aa8e7b13e8e258b27_enum,
  active fieldseeker.polygonlocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology fieldseeker.polygonlocation_locationsymbology_enum,
  externalid VARCHAR(50),
  acres DOUBLE PRECISION,
  nextactiondatescheduled TIMESTAMP,
  larvinspectinterval SMALLINT,
  zone2 VARCHAR(25),
  locationnumber INTEGER,
  globalid UUID,
  lastinspectdate TIMESTAMP,
  lastinspectbreeding VARCHAR(25),
  lastinspectavglarvae DOUBLE PRECISION,
  lastinspectavgpupae DOUBLE PRECISION,
  lastinspectlstages VARCHAR(25),
  lastinspectactiontaken VARCHAR(50),
  lastinspectfieldspecies VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  hectares DOUBLE PRECISION,
  lastinspectactivity VARCHAR(25),
  lasttreatactivity VARCHAR(25),
  lastinspectconditions VARCHAR(250),
  waterorigin fieldseeker.polygonlocation_polygonlocation_waterorigin_e9018e92_5f47_4ff9_8a7c_b818d848dc7a_enum,
  filter VARCHAR(255),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.polygonlocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.polygonlocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.polygonlocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.polygonlocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.polygonlocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.polygonlocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.polygonlocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.polygonlocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.polygonlocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.polygonlocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.polygonlocation.symbology IS 'Symbology';
COMMENT ON COLUMN fieldseeker.polygonlocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.polygonlocation.acres IS 'Acres';
COMMENT ON COLUMN fieldseeker.polygonlocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.polygonlocation.larvinspectinterval IS 'Larval Inspection Interval';
COMMENT ON COLUMN fieldseeker.polygonlocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectdate IS 'Last Inspection Date';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectbreeding IS 'Last Inspection Breeding';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectavglarvae IS 'Last Inspection Average Larvae';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectavgpupae IS 'Last Inspection Average Pupae';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectlstages IS 'Last Inspection Larval Stages';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectactiontaken IS 'Last Inspection Action';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectfieldspecies IS 'Last Inspection Field Species';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatdate IS 'Last Treatment Date';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatproduct IS 'Last Treatment Product';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatqty IS 'Last Treatment Quantity';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatqtyunit IS 'Last Treatment Quantity Unit';
COMMENT ON COLUMN fieldseeker.polygonlocation.hectares IS 'Hectares';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectactivity IS 'Last Inspection Activity';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatactivity IS 'Last Treatment Activity';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectconditions IS 'Last Inspection Conditions';
COMMENT ON COLUMN fieldseeker.polygonlocation.waterorigin IS 'Water Origin';
COMMENT ON COLUMN fieldseeker.polygonlocation.jurisdiction IS 'Jurisdiction';

-- Field active has default value: 1

-- Field symbology has default value: 'NONE'

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_polygonlocation_versioned(bigint, varchar, varchar, fieldseeker.polygonlocation_polygonlocation_habitat_45e9dde79ac84d959df8b65ba7d5dafd_enum, fieldseeker.polygonlocation_locationpriority_enum, fieldseeker.polygonlocation_polygonlocation_usetype_e546154cb9544b9aa8e7b13e8e258b27_enum, fieldseeker.polygonlocation_notinuit_f_enum, varchar, varchar, varchar, fieldseeker.polygonlocation_locationsymbology_enum, varchar, double precision, timestamp, smallint, varchar, integer, uuid, timestamp, varchar, double precision, double precision, varchar, varchar, varchar, timestamp, varchar, double precision, varchar, double precision, varchar, varchar, varchar, fieldseeker.polygonlocation_polygonlocation_waterorigin_e9018e92_5f47_4ff9_8a7c_b818d848dc7a_enum, varchar, timestamp, varchar, timestamp, varchar, varchar, double precision, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.polygonlocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.polygonlocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.polygonlocation (
  objectid, name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, acres, nextactiondatescheduled, larvinspectinterval, zone2, locationnumber, globalid, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, hectares, lastinspectactivity, lasttreatactivity, lastinspectconditions, waterorigin, filter, creationdate, creator, editdate, editor, jurisdiction, shape__area, shape__length,
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
    lv.name IS NOT DISTINCT FROM $2 AND
    lv.zone IS NOT DISTINCT FROM $3 AND
    lv.habitat IS NOT DISTINCT FROM $4 AND
    lv.priority IS NOT DISTINCT FROM $5 AND
    lv.usetype IS NOT DISTINCT FROM $6 AND
    lv.active IS NOT DISTINCT FROM $7 AND
    lv.description IS NOT DISTINCT FROM $8 AND
    lv.accessdesc IS NOT DISTINCT FROM $9 AND
    lv.comments IS NOT DISTINCT FROM $10 AND
    lv.symbology IS NOT DISTINCT FROM $11 AND
    lv.externalid IS NOT DISTINCT FROM $12 AND
    lv.acres IS NOT DISTINCT FROM $13 AND
    lv.nextactiondatescheduled IS NOT DISTINCT FROM $14 AND
    lv.larvinspectinterval IS NOT DISTINCT FROM $15 AND
    lv.zone2 IS NOT DISTINCT FROM $16 AND
    lv.locationnumber IS NOT DISTINCT FROM $17 AND
    lv.globalid IS NOT DISTINCT FROM $18 AND
    lv.lastinspectdate IS NOT DISTINCT FROM $19 AND
    lv.lastinspectbreeding IS NOT DISTINCT FROM $20 AND
    lv.lastinspectavglarvae IS NOT DISTINCT FROM $21 AND
    lv.lastinspectavgpupae IS NOT DISTINCT FROM $22 AND
    lv.lastinspectlstages IS NOT DISTINCT FROM $23 AND
    lv.lastinspectactiontaken IS NOT DISTINCT FROM $24 AND
    lv.lastinspectfieldspecies IS NOT DISTINCT FROM $25 AND
    lv.lasttreatdate IS NOT DISTINCT FROM $26 AND
    lv.lasttreatproduct IS NOT DISTINCT FROM $27 AND
    lv.lasttreatqty IS NOT DISTINCT FROM $28 AND
    lv.lasttreatqtyunit IS NOT DISTINCT FROM $29 AND
    lv.hectares IS NOT DISTINCT FROM $30 AND
    lv.lastinspectactivity IS NOT DISTINCT FROM $31 AND
    lv.lasttreatactivity IS NOT DISTINCT FROM $32 AND
    lv.lastinspectconditions IS NOT DISTINCT FROM $33 AND
    lv.waterorigin IS NOT DISTINCT FROM $34 AND
    lv.filter IS NOT DISTINCT FROM $35 AND
    lv.creationdate IS NOT DISTINCT FROM $36 AND
    lv.creator IS NOT DISTINCT FROM $37 AND
    lv.editdate IS NOT DISTINCT FROM $38 AND
    lv.editor IS NOT DISTINCT FROM $39 AND
    lv.jurisdiction IS NOT DISTINCT FROM $40 AND
    lv.shape__area IS NOT DISTINCT FROM $41 AND
    lv.shape__length IS NOT DISTINCT FROM $42
  )
RETURNING *;

-- Example usage: EXECUTE insert_polygonlocation_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.Pool
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.pool_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.pool_pool_testmethod_670efbfba86d41ba8e2d3cab5d749e7f_enum AS ENUM (
  'RT-PCR'
);

CREATE TYPE fieldseeker.pool_pool_diseasetested_0f02232949c04c7e8de820b9b515ed97_enum AS ENUM (
  'WNV',
  'SLEV',
  'WEEV',
  'DENV',
  'ZIKV',
  'CHIKV'
);

CREATE TYPE fieldseeker.pool_pool_diseasepos_6889f8dd00074874aa726907e78497fa_enum AS ENUM (
  'WNV',
  'SLEV',
  'WEEV',
  'DENV',
  'ZIKV',
  'CHIKV',
  'WNV/SLEV'
);

CREATE TYPE fieldseeker.pool_mosquitolabname_enum AS ENUM (
  'Internal Lab',
  'State Lab'
);

CREATE TABLE fieldseeker.pool (
  objectid BIGSERIAL NOT NULL,
  trapdata_id UUID,
  datesent TIMESTAMP,
  survtech VARCHAR(25),
  datetested TIMESTAMP,
  testtech VARCHAR(25),
  comments VARCHAR(250),
  sampleid VARCHAR(50),
  processed fieldseeker.pool_notinuit_f_enum,
  lab_id UUID,
  testmethod fieldseeker.pool_pool_testmethod_670efbfba86d41ba8e2d3cab5d749e7f_enum,
  diseasetested fieldseeker.pool_pool_diseasetested_0f02232949c04c7e8de820b9b515ed97_enum,
  diseasepos fieldseeker.pool_pool_diseasepos_6889f8dd00074874aa726907e78497fa_enum,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  lab fieldseeker.pool_mosquitolabname_enum,
  poolyear SMALLINT,
  gatewaysync SMALLINT,
  vectorsurvcollectionid VARCHAR(50),
  vectorsurvpoolid VARCHAR(50),
  vectorsurvtrapdataid VARCHAR(50),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.pool.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.pool.trapdata_id IS 'Trap Data ID';
COMMENT ON COLUMN fieldseeker.pool.datesent IS 'Date Sent';
COMMENT ON COLUMN fieldseeker.pool.survtech IS 'Survey Tech';
COMMENT ON COLUMN fieldseeker.pool.datetested IS 'Date Tested';
COMMENT ON COLUMN fieldseeker.pool.testtech IS 'Test Tech';
COMMENT ON COLUMN fieldseeker.pool.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.pool.sampleid IS 'Sample ID';
COMMENT ON COLUMN fieldseeker.pool.processed IS 'Processed';
COMMENT ON COLUMN fieldseeker.pool.testmethod IS 'Test Methods';
COMMENT ON COLUMN fieldseeker.pool.diseasetested IS 'Diseases Tested';
COMMENT ON COLUMN fieldseeker.pool.diseasepos IS 'Diseases Positive';
COMMENT ON COLUMN fieldseeker.pool.poolyear IS 'Pool Year';
COMMENT ON COLUMN fieldseeker.pool.gatewaysync IS 'Gateway Sync';

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.PoolDetail
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.pooldetail (
  objectid BIGSERIAL NOT NULL,
  trapdata_id UUID,
  pool_id UUID,
  species VARCHAR(50),
  females SMALLINT,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.pooldetail.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.pooldetail.trapdata_id IS 'Trap Data ID';
COMMENT ON COLUMN fieldseeker.pooldetail.pool_id IS 'Pool ID';
COMMENT ON COLUMN fieldseeker.pooldetail.species IS 'Species';
COMMENT ON COLUMN fieldseeker.pooldetail.females IS 'Females';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_pooldetail_versioned(bigint, uuid, uuid, varchar, smallint, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.pooldetail
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.pooldetail
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.pooldetail (
  objectid, trapdata_id, pool_id, species, females, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.trapdata_id IS NOT DISTINCT FROM $2 AND
    lv.pool_id IS NOT DISTINCT FROM $3 AND
    lv.species IS NOT DISTINCT FROM $4 AND
    lv.females IS NOT DISTINCT FROM $5 AND
    lv.globalid IS NOT DISTINCT FROM $6 AND
    lv.created_user IS NOT DISTINCT FROM $7 AND
    lv.created_date IS NOT DISTINCT FROM $8 AND
    lv.last_edited_user IS NOT DISTINCT FROM $9 AND
    lv.last_edited_date IS NOT DISTINCT FROM $10 AND
    lv.creationdate IS NOT DISTINCT FROM $11 AND
    lv.creator IS NOT DISTINCT FROM $12 AND
    lv.editdate IS NOT DISTINCT FROM $13 AND
    lv.editor IS NOT DISTINCT FROM $14
  )
RETURNING *;

-- Example usage: EXECUTE insert_pooldetail_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.ProposedTreatmentArea
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.proposedtreatmentarea_mosquitotreatmentmethod_enum AS ENUM (
  'Argo',
  'ATV',
  'Backpack',
  'Drone',
  'Manual',
  'Truck',
  'ULV',
  'Enhanced_Surveillance'
);

CREATE TYPE fieldseeker.proposedtreatmentarea_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.proposedtreatmentarea_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TABLE fieldseeker.proposedtreatmentarea (
  objectid BIGSERIAL NOT NULL,
  method fieldseeker.proposedtreatmentarea_mosquitotreatmentmethod_enum,
  comments VARCHAR(250),
  zone VARCHAR(25),
  reviewed fieldseeker.proposedtreatmentarea_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  zone2 VARCHAR(25),
  completeddate TIMESTAMP,
  completedby VARCHAR(25),
  completed fieldseeker.proposedtreatmentarea_notinuit_f_enum,
  issprayroute fieldseeker.proposedtreatmentarea_notinuit_f_enum,
  name VARCHAR(25),
  acres DOUBLE PRECISION,
  globalid UUID,
  exported fieldseeker.proposedtreatmentarea_notinuit_f_enum,
  targetproduct VARCHAR(25),
  targetapprate DOUBLE PRECISION,
  hectares DOUBLE PRECISION,
  lasttreatactivity VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  priority fieldseeker.proposedtreatmentarea_locationpriority_enum,
  duedate TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  targetspecies VARCHAR(250),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.method IS 'Method';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completeddate IS 'Completed Date';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completedby IS 'Completed By';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completed IS 'Completed';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.issprayroute IS 'Is Spray Route';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.name IS 'Name';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.acres IS 'Acres';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetproduct IS 'Target Product';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetapprate IS 'Target App Rate';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.hectares IS 'Hectares';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatactivity IS 'Last Treatment Activity';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatdate IS 'Last Treatment Date';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatproduct IS 'Last Treatment Product';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatqty IS 'Last Treatment Quantity';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatqtyunit IS 'Last Treatment Quantity Unit';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.duedate IS 'Due Date';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetspecies IS 'Target Species';

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.QAMosquitoInspection
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.qamosquitoinspection_qalarvaereason_enum AS ENUM (
  'Missed Area',
  'New Site',
  'Not Visited',
  'Rate Low',
  'Treated Recently',
  'Unknown',
  'Wrong Product'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qaaquaticorganisms_enum AS ENUM (
  'fish',
  'scuds',
  'snails'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qawatermovement_enum AS ENUM (
  'Fast',
  'Medium',
  'None',
  'Slow',
  'Very Slow'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qawaterduration_enum AS ENUM (
  '~month',
  '~week',
  '<1 week',
  '<day',
  '<month',
  '>month',
  '>week'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qawatersource_enum AS ENUM (
  'Irrigation',
  'Manually Controlled',
  'Percolation',
  'Rain Runoff',
  'Tidal',
  'Water Table'
);

CREATE TYPE fieldseeker.qamosquitoinspection_mosquitoaction_enum AS ENUM (
  'Treatment',
  'Covered container',
  'Cleared debris',
  'Maintenance'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qabreedingpotential_enum AS ENUM (
  'High',
  'Low',
  'Medium',
  'Rare'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qamosquitohabitat_enum AS ENUM (
  'Depressions',
  'Detritus present',
  'Fast',
  'Few predators',
  'Fluctuating levels',
  'H20<6"',
  'Low wave potential',
  'No fish',
  'Shallow edges',
  'Still water edges',
  'Still water whole',
  'Veg. on edges'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qavegetation_enum AS ENUM (
  'Algae',
  'Cattails',
  'Duckweed',
  'Glasswort',
  'Grass on edge',
  'Mangrove',
  'Mosquito fern',
  'Muskgrass',
  'Myriophyllum',
  'Other',
  'Rotting vegetation',
  'Saltwort',
  'Sedges'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qasourcereduction_enum AS ENUM (
  '1 tractor < day',
  'adjust flood irrigation',
  'adjust turf irrigation',
  'clear outflow',
  'cut ditch',
  'hand grading',
  'laser leveling',
  'multiple loads soil',
  'none'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qasoilcondition_enum AS ENUM (
  'Cracked',
  'Dry',
  'Inundated',
  'Saturated',
  'Surface Moist'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qawaterconditions_enum AS ENUM (
  '"rust" material',
  'Clear',
  'Cloudy/fines',
  'Floating debris',
  'Submerged/decom. debris'
);

CREATE TYPE fieldseeker.qamosquitoinspection_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qasitetype_enum AS ENUM (
  'Detention Pond',
  'Ditch',
  'Low Area',
  'Mangrove Edge',
  'Pond',
  'Pond Edge',
  'Swale'
);

CREATE TABLE fieldseeker.qamosquitoinspection (
  objectid BIGSERIAL NOT NULL,
  posdips SMALLINT,
  actiontaken fieldseeker.qamosquitoinspection_mosquitoaction_enum,
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  raingauge DOUBLE PRECISION,
  globalid UUID,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  winddir VARCHAR(3),
  reviewed fieldseeker.qamosquitoinspection_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  lr SMALLINT,
  negdips SMALLINT,
  totalacres DOUBLE PRECISION,
  acresbreeding DOUBLE PRECISION,
  fish fieldseeker.qamosquitoinspection_notinuit_f_enum,
  sitetype fieldseeker.qamosquitoinspection_qasitetype_enum,
  breedingpotential fieldseeker.qamosquitoinspection_qabreedingpotential_enum,
  movingwater fieldseeker.qamosquitoinspection_notinuit_f_enum,
  nowaterever fieldseeker.qamosquitoinspection_notinuit_f_enum,
  mosquitohabitat fieldseeker.qamosquitoinspection_qamosquitohabitat_enum,
  habvalue1 SMALLINT,
  habvalue1percent SMALLINT,
  habvalue2 SMALLINT,
  habvalue2percent SMALLINT,
  potential SMALLINT,
  larvaepresent fieldseeker.qamosquitoinspection_notinuit_f_enum,
  larvaeinsidetreatedarea fieldseeker.qamosquitoinspection_notinuit_f_enum,
  larvaeoutsidetreatedarea fieldseeker.qamosquitoinspection_notinuit_f_enum,
  larvaereason fieldseeker.qamosquitoinspection_qalarvaereason_enum,
  aquaticorganisms fieldseeker.qamosquitoinspection_qaaquaticorganisms_enum,
  vegetation fieldseeker.qamosquitoinspection_qavegetation_enum,
  sourcereduction fieldseeker.qamosquitoinspection_qasourcereduction_enum,
  waterpresent fieldseeker.qamosquitoinspection_notinuit_f_enum,
  watermovement1 fieldseeker.qamosquitoinspection_qawatermovement_enum,
  watermovement1percent SMALLINT,
  watermovement2 fieldseeker.qamosquitoinspection_qawatermovement_enum,
  watermovement2percent SMALLINT,
  soilconditions fieldseeker.qamosquitoinspection_qasoilcondition_enum,
  waterduration fieldseeker.qamosquitoinspection_qawaterduration_enum,
  watersource fieldseeker.qamosquitoinspection_qawatersource_enum,
  waterconditions fieldseeker.qamosquitoinspection_qawaterconditions_enum,
  adultactivity fieldseeker.qamosquitoinspection_notinuit_f_enum,
  linelocid UUID,
  pointlocid UUID,
  polygonlocid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  fieldtech VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.qamosquitoinspection.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.qamosquitoinspection.posdips IS 'Positive Dips';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.actiontaken IS 'Action';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.avetemp IS 'Average Temperature';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.windspeed IS 'Wind Speed';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.raingauge IS 'Rain Gauge';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.winddir IS 'Wind Direction';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.lr IS 'Landing Rate';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.negdips IS 'Negative Dips';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.totalacres IS 'Total Acres';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.acresbreeding IS 'Acres Breeding';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.fish IS 'Fish Present?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.sitetype IS 'Site Type';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.breedingpotential IS 'Breeding Potential';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.movingwater IS 'Moving Water';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.nowaterever IS 'No Evidence of Water Ever';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.mosquitohabitat IS 'Mosquito Habitat Indicators';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue1 IS 'Habitat Value';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue1percent IS '%';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue2 IS 'Habitat Value';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue2percent IS '%';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.potential IS 'Potential';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaepresent IS 'Larvae Present';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaeinsidetreatedarea IS 'Larvae Inside Treated Area?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaeoutsidetreatedarea IS 'Larvae Outside Treated Area?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaereason IS 'Reason Larvae Present';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.aquaticorganisms IS 'Aquatic Organisms';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.vegetation IS 'Vegetation';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.sourcereduction IS 'Source Reduction';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterpresent IS 'Water Present?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement1 IS 'Water Movement';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement1percent IS '%';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement2 IS 'Water Movement';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement2percent IS '%';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.soilconditions IS 'Soil Conditions';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterduration IS 'How Long Water Present?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watersource IS 'Water Source';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterconditions IS 'Water Conditions';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.adultactivity IS 'Adult Activity';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.fieldtech IS 'Field Tech';

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.RodentLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.rodentlocation_notinuit_f_1_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.rodentlocation_rodentlocation_symbology_enum AS ENUM (
  'ACTION',
  'INACTIVE',
  'NONE'
);

CREATE TYPE fieldseeker.rodentlocation_rodentlocationhabitat_enum AS ENUM (
  'Commercial',
  'Industrial',
  'Residential',
  'Wood Pile'
);

CREATE TYPE fieldseeker.rodentlocation_locationpriority_1_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TYPE fieldseeker.rodentlocation_locationusetype_1_enum AS ENUM (
  'Residential',
  'Commercial',
  'Industrial',
  'Agricultural',
  'Mixed use'
);

CREATE TABLE fieldseeker.rodentlocation (
  objectid BIGSERIAL NOT NULL,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  zone2 VARCHAR(25),
  habitat fieldseeker.rodentlocation_rodentlocationhabitat_enum,
  priority fieldseeker.rodentlocation_locationpriority_1_enum,
  usetype fieldseeker.rodentlocation_locationusetype_1_enum,
  active fieldseeker.rodentlocation_notinuit_f_1_enum,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology fieldseeker.rodentlocation_rodentlocation_symbology_enum,
  externalid VARCHAR(50),
  nextactiondatescheduled TIMESTAMP,
  locationnumber INTEGER,
  lastinspectdate TIMESTAMP,
  lastinspectspecies VARCHAR(25),
  lastinspectaction VARCHAR(50),
  lastinspectconditions VARCHAR(250),
  lastinspectrodentevidence VARCHAR(250),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.rodentlocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.rodentlocation.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.rodentlocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectdate IS 'Last Inspection Date';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectspecies IS 'Last Inspection Species';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectaction IS 'Last Inspection Action';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectconditions IS 'Last Inspection Conditions';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectrodentevidence IS 'Last Inspection Rodent Evidence';
COMMENT ON COLUMN fieldseeker.rodentlocation.jurisdiction IS 'Jurisdiction';

-- Field symbology has default value: 'NONE'

-- Field active has default value: 1

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_rodentlocation_versioned(bigint, varchar, varchar, varchar, fieldseeker.rodentlocation_rodentlocationhabitat_enum, fieldseeker.rodentlocation_locationpriority_1_enum, fieldseeker.rodentlocation_locationusetype_1_enum, fieldseeker.rodentlocation_notinuit_f_1_enum, varchar, varchar, varchar, fieldseeker.rodentlocation_rodentlocation_symbology_enum, varchar, timestamp, integer, timestamp, varchar, varchar, varchar, varchar, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.rodentlocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.rodentlocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.rodentlocation (
  objectid, locationname, zone, zone2, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, nextactiondatescheduled, locationnumber, lastinspectdate, lastinspectspecies, lastinspectaction, lastinspectconditions, lastinspectrodentevidence, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, jurisdiction,
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
    lv.locationname IS NOT DISTINCT FROM $2 AND
    lv.zone IS NOT DISTINCT FROM $3 AND
    lv.zone2 IS NOT DISTINCT FROM $4 AND
    lv.habitat IS NOT DISTINCT FROM $5 AND
    lv.priority IS NOT DISTINCT FROM $6 AND
    lv.usetype IS NOT DISTINCT FROM $7 AND
    lv.active IS NOT DISTINCT FROM $8 AND
    lv.description IS NOT DISTINCT FROM $9 AND
    lv.accessdesc IS NOT DISTINCT FROM $10 AND
    lv.comments IS NOT DISTINCT FROM $11 AND
    lv.symbology IS NOT DISTINCT FROM $12 AND
    lv.externalid IS NOT DISTINCT FROM $13 AND
    lv.nextactiondatescheduled IS NOT DISTINCT FROM $14 AND
    lv.locationnumber IS NOT DISTINCT FROM $15 AND
    lv.lastinspectdate IS NOT DISTINCT FROM $16 AND
    lv.lastinspectspecies IS NOT DISTINCT FROM $17 AND
    lv.lastinspectaction IS NOT DISTINCT FROM $18 AND
    lv.lastinspectconditions IS NOT DISTINCT FROM $19 AND
    lv.lastinspectrodentevidence IS NOT DISTINCT FROM $20 AND
    lv.globalid IS NOT DISTINCT FROM $21 AND
    lv.created_user IS NOT DISTINCT FROM $22 AND
    lv.created_date IS NOT DISTINCT FROM $23 AND
    lv.last_edited_user IS NOT DISTINCT FROM $24 AND
    lv.last_edited_date IS NOT DISTINCT FROM $25 AND
    lv.creationdate IS NOT DISTINCT FROM $26 AND
    lv.creator IS NOT DISTINCT FROM $27 AND
    lv.editdate IS NOT DISTINCT FROM $28 AND
    lv.editor IS NOT DISTINCT FROM $29 AND
    lv.jurisdiction IS NOT DISTINCT FROM $30
  )
RETURNING *;

-- Example usage: EXECUTE insert_rodentlocation_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.SampleCollection
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.samplecollection_mosquitolabname_enum AS ENUM (
  'Internal Lab',
  'State Lab'
);

CREATE TYPE fieldseeker.samplecollection_mosquitosampletype_enum AS ENUM (
  'Blood',
  'Tissue',
  'Specimen',
  'Carcass'
);

CREATE TYPE fieldseeker.samplecollection_mosquitosamplecondition_enum AS ENUM (
  'Usable',
  'Unusable'
);

CREATE TYPE fieldseeker.samplecollection_mosquitosamplespecies_enum AS ENUM (
  'Chicken',
  'Wild bird',
  'Horse',
  'Human'
);

CREATE TYPE fieldseeker.samplecollection_notinuisex_enum AS ENUM (
  'M',
  'F',
  'U'
);

CREATE TYPE fieldseeker.samplecollection_mosquitoactivity_enum AS ENUM (
  'Routine inspection',
  'Pre-treatment',
  'Maintenance',
  'ULV',
  'BARRIER',
  'LOGIN',
  'TREATSD',
  'SD',
  'SITEVISIT',
  'ONLINE',
  'SYNC',
  'CREATESR',
  'LC',
  'ACCEPTSR',
  'POINT',
  'DOWNLOAD',
  'COMPLETESR',
  'POLYGON',
  'TRAP',
  'SAMPLE',
  'QA',
  'PTA',
  'FIELDSCOUTING',
  'OFFLINE',
  'LINE',
  'TRAPLOCATION',
  'SAMPLELOCATION',
  'LCLOCATION'
);

CREATE TYPE fieldseeker.samplecollection_mosquitodisease_enum AS ENUM (
  'EEE',
  'WNV',
  'Dengue',
  'Zika'
);

CREATE TYPE fieldseeker.samplecollection_mosquitositecondition_enum AS ENUM (
  'Dry',
  'Clean',
  'Full',
  'Low'
);

CREATE TYPE fieldseeker.samplecollection_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.samplecollection_notinuiwinddirection_enum AS ENUM (
  'N',
  'NE',
  'E',
  'SE',
  'S',
  'SW',
  'W',
  'NW'
);

CREATE TYPE fieldseeker.samplecollection_mosquitotestmethod_enum AS ENUM (
  'RAMP',
  'VecTest',
  'ELISA',
  'RT-PCR'
);

CREATE TABLE fieldseeker.samplecollection (
  objectid BIGSERIAL NOT NULL,
  loc_id UUID,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  sitecond fieldseeker.samplecollection_mosquitositecondition_enum,
  sampleid VARCHAR(25),
  survtech VARCHAR(25),
  datesent TIMESTAMP,
  datetested TIMESTAMP,
  testtech VARCHAR(25),
  comments VARCHAR(250),
  processed fieldseeker.samplecollection_notinuit_f_enum,
  sampletype fieldseeker.samplecollection_mosquitosampletype_enum,
  samplecond fieldseeker.samplecollection_mosquitosamplecondition_enum,
  species fieldseeker.samplecollection_mosquitosamplespecies_enum,
  sex fieldseeker.samplecollection_notinuisex_enum,
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  winddir fieldseeker.samplecollection_notinuiwinddirection_enum,
  raingauge DOUBLE PRECISION,
  activity fieldseeker.samplecollection_mosquitoactivity_enum,
  testmethod fieldseeker.samplecollection_mosquitotestmethod_enum,
  diseasetested fieldseeker.samplecollection_mosquitodisease_enum,
  diseasepos fieldseeker.samplecollection_mosquitodisease_enum,
  reviewed fieldseeker.samplecollection_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  lab fieldseeker.samplecollection_mosquitolabname_enum,
  fieldtech VARCHAR(25),
  flockid UUID,
  samplecount SMALLINT,
  chickenid UUID,
  gatewaysync SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.samplecollection.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.samplecollection.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.samplecollection.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.samplecollection.sitecond IS 'Conditions';
COMMENT ON COLUMN fieldseeker.samplecollection.sampleid IS 'Sample ID';
COMMENT ON COLUMN fieldseeker.samplecollection.survtech IS 'Surveillance Technician';
COMMENT ON COLUMN fieldseeker.samplecollection.datesent IS 'Sent';
COMMENT ON COLUMN fieldseeker.samplecollection.datetested IS 'Tested';
COMMENT ON COLUMN fieldseeker.samplecollection.testtech IS 'Test Technician';
COMMENT ON COLUMN fieldseeker.samplecollection.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.samplecollection.processed IS 'Processed';
COMMENT ON COLUMN fieldseeker.samplecollection.sampletype IS 'Sample Type';
COMMENT ON COLUMN fieldseeker.samplecollection.samplecond IS 'Sample Condition';
COMMENT ON COLUMN fieldseeker.samplecollection.species IS 'Species';
COMMENT ON COLUMN fieldseeker.samplecollection.sex IS 'Sex';
COMMENT ON COLUMN fieldseeker.samplecollection.avetemp IS 'Average Temperature';
COMMENT ON COLUMN fieldseeker.samplecollection.windspeed IS 'Wind Speed';
COMMENT ON COLUMN fieldseeker.samplecollection.winddir IS 'Wind Direction';
COMMENT ON COLUMN fieldseeker.samplecollection.raingauge IS 'Rain Gauge';
COMMENT ON COLUMN fieldseeker.samplecollection.activity IS 'Activity';
COMMENT ON COLUMN fieldseeker.samplecollection.testmethod IS 'Test Method';
COMMENT ON COLUMN fieldseeker.samplecollection.diseasetested IS 'Disease Tested';
COMMENT ON COLUMN fieldseeker.samplecollection.diseasepos IS 'Disease Positive';
COMMENT ON COLUMN fieldseeker.samplecollection.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.samplecollection.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.samplecollection.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.samplecollection.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.samplecollection.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.samplecollection.recordstatus IS 'RecordStatus';
COMMENT ON COLUMN fieldseeker.samplecollection.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.samplecollection.lab IS 'Lab';
COMMENT ON COLUMN fieldseeker.samplecollection.fieldtech IS 'Field Tech';
COMMENT ON COLUMN fieldseeker.samplecollection.samplecount IS 'Sample Count';
COMMENT ON COLUMN fieldseeker.samplecollection.chickenid IS 'ChickenID';
COMMENT ON COLUMN fieldseeker.samplecollection.gatewaysync IS 'Gateway Sync';

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.SampleLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.samplelocation_locationhabitattype_enum AS ENUM (
  'Catch basin',
  'Creek',
  'Ditch',
  'Field/Pasture',
  'Pond',
  'Pond fish',
  'Pond marshy',
  'Pond ornamental',
  'Pond retention',
  'Pond sewage',
  'Pond woodland',
  'Tree hole',
  'Swimming pool',
  'Park',
  'Unknown'
);

CREATE TYPE fieldseeker.samplelocation_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TYPE fieldseeker.samplelocation_samplelocationusetype_enum AS ENUM (
  'Flock Site',
  'Dead Bird'
);

CREATE TYPE fieldseeker.samplelocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TABLE fieldseeker.samplelocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.samplelocation_locationhabitattype_enum,
  priority fieldseeker.samplelocation_locationpriority_enum,
  usetype fieldseeker.samplelocation_samplelocationusetype_enum,
  active fieldseeker.samplelocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  externalid VARCHAR(50),
  nextactiondatescheduled TIMESTAMP,
  zone2 VARCHAR(25),
  locationnumber INTEGER,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  gatewaysync SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.samplelocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.samplelocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.samplelocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.samplelocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.samplelocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.samplelocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.samplelocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.samplelocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.samplelocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.samplelocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.samplelocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.samplelocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.samplelocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.samplelocation.gatewaysync IS 'Gateway Sync';

-- Field active has default value: 1

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_samplelocation_versioned(bigint, varchar, varchar, fieldseeker.samplelocation_locationhabitattype_enum, fieldseeker.samplelocation_locationpriority_enum, fieldseeker.samplelocation_samplelocationusetype_enum, fieldseeker.samplelocation_notinuit_f_enum, varchar, varchar, varchar, varchar, timestamp, varchar, integer, uuid, varchar, timestamp, varchar, timestamp, smallint, timestamp, varchar, timestamp, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.samplelocation
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.samplelocation
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.samplelocation (
  objectid, name, zone, habitat, priority, usetype, active, description, accessdesc, comments, externalid, nextactiondatescheduled, zone2, locationnumber, globalid, created_user, created_date, last_edited_user, last_edited_date, gatewaysync, creationdate, creator, editdate, editor,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24,
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
    lv.creationdate IS NOT DISTINCT FROM $21 AND
    lv.creator IS NOT DISTINCT FROM $22 AND
    lv.editdate IS NOT DISTINCT FROM $23 AND
    lv.editor IS NOT DISTINCT FROM $24
  )
RETURNING *;

-- Example usage: EXECUTE insert_samplelocation_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.ServiceRequest
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.servicerequest_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestcontactpreferences_enum AS ENUM (
  'None',
  'Call',
  'Email',
  'Text'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestrejectedreason_enum AS ENUM (
  'Distance',
  'Workload'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_spanish_aaa3dc66_9f9a_4527_8ecd_c9f76db33879_enum AS ENUM (
  '0',
  '1'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestissues_enum AS ENUM (
  'Beehive Related',
  'Unsanitary Accumulations',
  'Rooster or Noise',
  'Rats Attracted',
  'Odor',
  'Number of Animals Over Limit',
  'Location',
  'Violation',
  'Inadequate Enclosure',
  'Escaped Animal',
  'Illegal Animal'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestsource_enum AS ENUM (
  'Phone',
  'Email',
  'Website',
  'Drop-in',
  '2025_pools'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestregion_enum AS ENUM (
  'FL',
  'ID'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_schedule_period_3f40c046_afd1_4abd_8bf4_389650d29a49_enum AS ENUM (
  'AM',
  'PM'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'Follow up Visit',
  'HTC Response',
  'Disease Activity Response'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_supervisor_eba07b90_c885_4fe6_8080_7aa775403b9a_enum AS ENUM (
  'Rick Alverez',
  'Bryan Ferguson',
  'Bryan Ruiz',
  'Andrea Troupin',
  'Conlin Reis'
);

CREATE TYPE fieldseeker.servicerequest_servicerequesttarget_enum AS ENUM (
  'mosquitofish',
  'neglected pool or spa',
  'standing water',
  'mosquito presence',
  'biting mosquitoes',
  'event',
  'fish',
  'mosquito',
  'source',
  'bird'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_dog_2b95ec97_1286_4fcd_88f4_f0e31113f696_enum AS ENUM (
  '0',
  '1',
  '2',
  '3'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_assignedtech_71d0d685_868f_4b7a_87e2_3661a3ee67c5_enum AS ENUM (
  'Alysia Davis',
  'Alejandra Gill',
  'Andrea Troupin',
  'Brenda Rodriguez',
  'Bryan Ferguson',
  'Bryan Ruiz',
  'Conlin Reis',
  'Carlos Rodriguez',
  'Erick Arriga',
  'Landon McGill',
  'Marco Martinez',
  'Mark Nakata',
  'Mario Sanchez',
  'Juan Pablo Ortega',
  'Ryan Spratt',
  'Ted McGill',
  'Benjamin Sperry',
  'Zachery Barragan',
  'Arturo Garcia-Trejo',
  'Jesus Jolano',
  'Yajaira Godinez',
  'Jake Maldonado',
  'Rafael Ramirez',
  'Lisa Salgado',
  'Kory Wilson',
  'Carlos Palacios',
  'Fatima Hidalgo',
  'Aaron Fredrick',
  'Josh Malone',
  'Jorge Perez',
  'Laura Romos'
);

CREATE TYPE fieldseeker.servicerequest_servicerequeststatus_enum AS ENUM (
  'Assigned',
  'Closed',
  'FieldRectified',
  'Open',
  'Rejected',
  'Unverified',
  'Accepted'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestnextaction_enum AS ENUM (
  'Night spray',
  'Site visit'
);

CREATE TABLE fieldseeker.servicerequest (
  objectid BIGSERIAL NOT NULL,
  recdatetime TIMESTAMP,
  source fieldseeker.servicerequest_servicerequestsource_enum,
  entrytech VARCHAR(25),
  priority fieldseeker.servicerequest_servicerequestpriority_enum,
  supervisor fieldseeker.servicerequest_servicerequest_supervisor_eba07b90_c885_4fe6_8080_7aa775403b9a_enum,
  assignedtech fieldseeker.servicerequest_servicerequest_assignedtech_71d0d685_868f_4b7a_87e2_3661a3ee67c5_enum,
  status fieldseeker.servicerequest_servicerequeststatus_enum,
  clranon fieldseeker.servicerequest_notinuit_f_enum,
  clrfname VARCHAR(25),
  clrphone1 VARCHAR(25),
  clrphone2 VARCHAR(25),
  clremail VARCHAR(50),
  clrcompany VARCHAR(25),
  clraddr1 VARCHAR(50),
  clraddr2 VARCHAR(50),
  clrcity VARCHAR(25),
  clrstate fieldseeker.servicerequest_servicerequestregion_enum,
  clrzip VARCHAR(25),
  clrother VARCHAR(25),
  clrcontpref fieldseeker.servicerequest_servicerequestcontactpreferences_enum,
  reqcompany VARCHAR(25),
  reqaddr1 VARCHAR(50),
  reqaddr2 VARCHAR(50),
  reqcity VARCHAR(25),
  reqstate fieldseeker.servicerequest_servicerequestregion_enum,
  reqzip VARCHAR(25),
  reqcrossst VARCHAR(25),
  reqsubdiv VARCHAR(25),
  reqmapgrid VARCHAR(25),
  reqpermission fieldseeker.servicerequest_notinuit_f_enum,
  reqtarget fieldseeker.servicerequest_servicerequesttarget_enum,
  reqdescr VARCHAR(1000),
  reqnotesfortech VARCHAR(250),
  reqnotesforcust VARCHAR(250),
  reqfldnotes VARCHAR(250),
  reqprogramactions VARCHAR(250),
  datetimeclosed TIMESTAMP,
  techclosed VARCHAR(25),
  sr_number INTEGER,
  reviewed fieldseeker.servicerequest_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  accepted fieldseeker.servicerequest_notinuit_f_enum,
  accepteddate TIMESTAMP,
  rejectedby VARCHAR(25),
  rejecteddate TIMESTAMP,
  rejectedreason fieldseeker.servicerequest_servicerequestrejectedreason_enum,
  duedate TIMESTAMP,
  acceptedby VARCHAR(25),
  comments VARCHAR(2500),
  estcompletedate TIMESTAMP,
  nextaction fieldseeker.servicerequest_servicerequestnextaction_enum,
  recordstatus SMALLINT,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  firstresponsedate TIMESTAMP,
  responsedaycount SMALLINT,
  allowed VARCHAR(25),
  xvalue VARCHAR(25),
  yvalue VARCHAR(25),
  validx VARCHAR(25),
  validy VARCHAR(25),
  externalid VARCHAR(25),
  externalerror VARCHAR(500),
  pointlocid UUID,
  notified SMALLINT,
  notifieddate TIMESTAMP,
  scheduled SMALLINT,
  scheduleddate TIMESTAMP,
  dog fieldseeker.servicerequest_servicerequest_dog_2b95ec97_1286_4fcd_88f4_f0e31113f696_enum,
  schedule_period fieldseeker.servicerequest_servicerequest_schedule_period_3f40c046_afd1_4abd_8bf4_389650d29a49_enum,
  schedule_notes VARCHAR(256),
  spanish fieldseeker.servicerequest_servicerequest_spanish_aaa3dc66_9f9a_4527_8ecd_c9f76db33879_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  issuesreported fieldseeker.servicerequest_servicerequestissues_enum,
  jurisdiction VARCHAR(25),
  notificationtimestamp VARCHAR(250),
  zone VARCHAR(50),
  zone2 VARCHAR(50),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.servicerequest.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.servicerequest.recdatetime IS 'Received';
COMMENT ON COLUMN fieldseeker.servicerequest.source IS 'Source';
COMMENT ON COLUMN fieldseeker.servicerequest.entrytech IS 'Entered By';
COMMENT ON COLUMN fieldseeker.servicerequest.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.servicerequest.supervisor IS 'Supervisor';
COMMENT ON COLUMN fieldseeker.servicerequest.assignedtech IS 'Assigned To';
COMMENT ON COLUMN fieldseeker.servicerequest.status IS 'Status';
COMMENT ON COLUMN fieldseeker.servicerequest.clranon IS 'Anonymous Caller';
COMMENT ON COLUMN fieldseeker.servicerequest.clrfname IS 'Caller Name';
COMMENT ON COLUMN fieldseeker.servicerequest.clrphone1 IS 'Caller Phone';
COMMENT ON COLUMN fieldseeker.servicerequest.clrphone2 IS 'Caller Alternate Phone';
COMMENT ON COLUMN fieldseeker.servicerequest.clremail IS 'Caller Email';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcompany IS 'Caller Company';
COMMENT ON COLUMN fieldseeker.servicerequest.clraddr1 IS 'Caller Address';
COMMENT ON COLUMN fieldseeker.servicerequest.clraddr2 IS 'Caller Address 2';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcity IS 'Caller City';
COMMENT ON COLUMN fieldseeker.servicerequest.clrstate IS 'Caller State';
COMMENT ON COLUMN fieldseeker.servicerequest.clrzip IS 'Caller ZIP';
COMMENT ON COLUMN fieldseeker.servicerequest.clrother IS 'Caller Other';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcontpref IS 'Caller Contact Preference';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcompany IS 'Request Company';
COMMENT ON COLUMN fieldseeker.servicerequest.reqaddr1 IS 'Request Address';
COMMENT ON COLUMN fieldseeker.servicerequest.reqaddr2 IS 'Request Address 2';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcity IS 'Request City';
COMMENT ON COLUMN fieldseeker.servicerequest.reqstate IS 'Request State';
COMMENT ON COLUMN fieldseeker.servicerequest.reqzip IS 'Request ZIP';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcrossst IS 'Request Cross Street';
COMMENT ON COLUMN fieldseeker.servicerequest.reqsubdiv IS 'Request Subdivision';
COMMENT ON COLUMN fieldseeker.servicerequest.reqmapgrid IS 'Request Map Grid';
COMMENT ON COLUMN fieldseeker.servicerequest.reqpermission IS 'Permission to Enter';
COMMENT ON COLUMN fieldseeker.servicerequest.reqtarget IS 'Request Target';
COMMENT ON COLUMN fieldseeker.servicerequest.reqdescr IS 'Request Description';
COMMENT ON COLUMN fieldseeker.servicerequest.reqnotesfortech IS 'Notes for Field Technician';
COMMENT ON COLUMN fieldseeker.servicerequest.reqnotesforcust IS 'Notes for Customer';
COMMENT ON COLUMN fieldseeker.servicerequest.reqfldnotes IS 'Request Field Notes';
COMMENT ON COLUMN fieldseeker.servicerequest.reqprogramactions IS 'Request Program Actions';
COMMENT ON COLUMN fieldseeker.servicerequest.datetimeclosed IS 'Closed';
COMMENT ON COLUMN fieldseeker.servicerequest.techclosed IS 'Closed By';
COMMENT ON COLUMN fieldseeker.servicerequest.sr_number IS 'SR#';
COMMENT ON COLUMN fieldseeker.servicerequest.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.servicerequest.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.servicerequest.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.servicerequest.accepted IS 'Accepted';
COMMENT ON COLUMN fieldseeker.servicerequest.accepteddate IS 'Accepted Date';
COMMENT ON COLUMN fieldseeker.servicerequest.rejectedby IS 'Rejected By';
COMMENT ON COLUMN fieldseeker.servicerequest.rejecteddate IS 'Rejected Date';
COMMENT ON COLUMN fieldseeker.servicerequest.rejectedreason IS 'Rejected Reason';
COMMENT ON COLUMN fieldseeker.servicerequest.duedate IS 'Due Date';
COMMENT ON COLUMN fieldseeker.servicerequest.acceptedby IS 'Accepted By';
COMMENT ON COLUMN fieldseeker.servicerequest.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.servicerequest.estcompletedate IS 'Estimated Completion Date';
COMMENT ON COLUMN fieldseeker.servicerequest.nextaction IS 'Next Action';
COMMENT ON COLUMN fieldseeker.servicerequest.recordstatus IS 'Record Status';
COMMENT ON COLUMN fieldseeker.servicerequest.firstresponsedate IS 'First Response Date';
COMMENT ON COLUMN fieldseeker.servicerequest.responsedaycount IS 'Response Day Count';
COMMENT ON COLUMN fieldseeker.servicerequest.allowed IS 'Verify Correct Location';
COMMENT ON COLUMN fieldseeker.servicerequest.xvalue IS 'Xvalue';
COMMENT ON COLUMN fieldseeker.servicerequest.yvalue IS 'Yvalue';
COMMENT ON COLUMN fieldseeker.servicerequest.validx IS 'ValidX';
COMMENT ON COLUMN fieldseeker.servicerequest.validy IS 'ValidY';
COMMENT ON COLUMN fieldseeker.servicerequest.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.servicerequest.externalerror IS 'External Error';
COMMENT ON COLUMN fieldseeker.servicerequest.notified IS 'Notified';
COMMENT ON COLUMN fieldseeker.servicerequest.notifieddate IS 'Notified Date';
COMMENT ON COLUMN fieldseeker.servicerequest.scheduled IS 'Scheduled';
COMMENT ON COLUMN fieldseeker.servicerequest.scheduleddate IS 'Scheduled Date';
COMMENT ON COLUMN fieldseeker.servicerequest.dog IS 'Dog';
COMMENT ON COLUMN fieldseeker.servicerequest.schedule_period IS 'Schedule Period';
COMMENT ON COLUMN fieldseeker.servicerequest.schedule_notes IS 'Schedule Notes';
COMMENT ON COLUMN fieldseeker.servicerequest.spanish IS 'Prefer speaking Spanish';
COMMENT ON COLUMN fieldseeker.servicerequest.issuesreported IS 'Issues Reported';
COMMENT ON COLUMN fieldseeker.servicerequest.jurisdiction IS 'Jurisdiction';
COMMENT ON COLUMN fieldseeker.servicerequest.notificationtimestamp IS 'Notification Timestamp';
COMMENT ON COLUMN fieldseeker.servicerequest.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.servicerequest.zone2 IS 'Zone2';

-- Field dog has default value: 0

-- Field spanish has default value: 0

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_servicerequest_versioned(bigint, timestamp, fieldseeker.servicerequest_servicerequestsource_enum, varchar, fieldseeker.servicerequest_servicerequestpriority_enum, fieldseeker.servicerequest_servicerequest_supervisor_eba07b90_c885_4fe6_8080_7aa775403b9a_enum, fieldseeker.servicerequest_servicerequest_assignedtech_71d0d685_868f_4b7a_87e2_3661a3ee67c5_enum, fieldseeker.servicerequest_servicerequeststatus_enum, fieldseeker.servicerequest_notinuit_f_enum, varchar, varchar, varchar, varchar, varchar, varchar, varchar, varchar, fieldseeker.servicerequest_servicerequestregion_enum, varchar, varchar, fieldseeker.servicerequest_servicerequestcontactpreferences_enum, varchar, varchar, varchar, varchar, fieldseeker.servicerequest_servicerequestregion_enum, varchar, varchar, varchar, varchar, fieldseeker.servicerequest_notinuit_f_enum, fieldseeker.servicerequest_servicerequesttarget_enum, varchar, varchar, varchar, varchar, varchar, timestamp, varchar, integer, fieldseeker.servicerequest_notinuit_f_enum, varchar, timestamp, fieldseeker.servicerequest_notinuit_f_enum, timestamp, varchar, timestamp, fieldseeker.servicerequest_servicerequestrejectedreason_enum, timestamp, varchar, varchar, timestamp, fieldseeker.servicerequest_servicerequestnextaction_enum, smallint, uuid, varchar, timestamp, varchar, timestamp, timestamp, smallint, varchar, varchar, varchar, varchar, varchar, varchar, varchar, uuid, smallint, timestamp, smallint, timestamp, fieldseeker.servicerequest_servicerequest_dog_2b95ec97_1286_4fcd_88f4_f0e31113f696_enum, fieldseeker.servicerequest_servicerequest_schedule_period_3f40c046_afd1_4abd_8bf4_389650d29a49_enum, varchar, fieldseeker.servicerequest_servicerequest_spanish_aaa3dc66_9f9a_4527_8ecd_c9f76db33879_enum, timestamp, varchar, timestamp, varchar, fieldseeker.servicerequest_servicerequestissues_enum, varchar, varchar, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.servicerequest
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.servicerequest
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.servicerequest (
  objectid, recdatetime, source, entrytech, priority, supervisor, assignedtech, status, clranon, clrfname, clrphone1, clrphone2, clremail, clrcompany, clraddr1, clraddr2, clrcity, clrstate, clrzip, clrother, clrcontpref, reqcompany, reqaddr1, reqaddr2, reqcity, reqstate, reqzip, reqcrossst, reqsubdiv, reqmapgrid, reqpermission, reqtarget, reqdescr, reqnotesfortech, reqnotesforcust, reqfldnotes, reqprogramactions, datetimeclosed, techclosed, sr_number, reviewed, reviewedby, revieweddate, accepted, accepteddate, rejectedby, rejecteddate, rejectedreason, duedate, acceptedby, comments, estcompletedate, nextaction, recordstatus, globalid, created_user, created_date, last_edited_user, last_edited_date, firstresponsedate, responsedaycount, allowed, xvalue, yvalue, validx, validy, externalid, externalerror, pointlocid, notified, notifieddate, scheduled, scheduleddate, dog, schedule_period, schedule_notes, spanish, creationdate, creator, editdate, editor, issuesreported, jurisdiction, notificationtimestamp, zone, zone2,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.recdatetime IS NOT DISTINCT FROM $2 AND
    lv.source IS NOT DISTINCT FROM $3 AND
    lv.entrytech IS NOT DISTINCT FROM $4 AND
    lv.priority IS NOT DISTINCT FROM $5 AND
    lv.supervisor IS NOT DISTINCT FROM $6 AND
    lv.assignedtech IS NOT DISTINCT FROM $7 AND
    lv.status IS NOT DISTINCT FROM $8 AND
    lv.clranon IS NOT DISTINCT FROM $9 AND
    lv.clrfname IS NOT DISTINCT FROM $10 AND
    lv.clrphone1 IS NOT DISTINCT FROM $11 AND
    lv.clrphone2 IS NOT DISTINCT FROM $12 AND
    lv.clremail IS NOT DISTINCT FROM $13 AND
    lv.clrcompany IS NOT DISTINCT FROM $14 AND
    lv.clraddr1 IS NOT DISTINCT FROM $15 AND
    lv.clraddr2 IS NOT DISTINCT FROM $16 AND
    lv.clrcity IS NOT DISTINCT FROM $17 AND
    lv.clrstate IS NOT DISTINCT FROM $18 AND
    lv.clrzip IS NOT DISTINCT FROM $19 AND
    lv.clrother IS NOT DISTINCT FROM $20 AND
    lv.clrcontpref IS NOT DISTINCT FROM $21 AND
    lv.reqcompany IS NOT DISTINCT FROM $22 AND
    lv.reqaddr1 IS NOT DISTINCT FROM $23 AND
    lv.reqaddr2 IS NOT DISTINCT FROM $24 AND
    lv.reqcity IS NOT DISTINCT FROM $25 AND
    lv.reqstate IS NOT DISTINCT FROM $26 AND
    lv.reqzip IS NOT DISTINCT FROM $27 AND
    lv.reqcrossst IS NOT DISTINCT FROM $28 AND
    lv.reqsubdiv IS NOT DISTINCT FROM $29 AND
    lv.reqmapgrid IS NOT DISTINCT FROM $30 AND
    lv.reqpermission IS NOT DISTINCT FROM $31 AND
    lv.reqtarget IS NOT DISTINCT FROM $32 AND
    lv.reqdescr IS NOT DISTINCT FROM $33 AND
    lv.reqnotesfortech IS NOT DISTINCT FROM $34 AND
    lv.reqnotesforcust IS NOT DISTINCT FROM $35 AND
    lv.reqfldnotes IS NOT DISTINCT FROM $36 AND
    lv.reqprogramactions IS NOT DISTINCT FROM $37 AND
    lv.datetimeclosed IS NOT DISTINCT FROM $38 AND
    lv.techclosed IS NOT DISTINCT FROM $39 AND
    lv.sr_number IS NOT DISTINCT FROM $40 AND
    lv.reviewed IS NOT DISTINCT FROM $41 AND
    lv.reviewedby IS NOT DISTINCT FROM $42 AND
    lv.revieweddate IS NOT DISTINCT FROM $43 AND
    lv.accepted IS NOT DISTINCT FROM $44 AND
    lv.accepteddate IS NOT DISTINCT FROM $45 AND
    lv.rejectedby IS NOT DISTINCT FROM $46 AND
    lv.rejecteddate IS NOT DISTINCT FROM $47 AND
    lv.rejectedreason IS NOT DISTINCT FROM $48 AND
    lv.duedate IS NOT DISTINCT FROM $49 AND
    lv.acceptedby IS NOT DISTINCT FROM $50 AND
    lv.comments IS NOT DISTINCT FROM $51 AND
    lv.estcompletedate IS NOT DISTINCT FROM $52 AND
    lv.nextaction IS NOT DISTINCT FROM $53 AND
    lv.recordstatus IS NOT DISTINCT FROM $54 AND
    lv.globalid IS NOT DISTINCT FROM $55 AND
    lv.created_user IS NOT DISTINCT FROM $56 AND
    lv.created_date IS NOT DISTINCT FROM $57 AND
    lv.last_edited_user IS NOT DISTINCT FROM $58 AND
    lv.last_edited_date IS NOT DISTINCT FROM $59 AND
    lv.firstresponsedate IS NOT DISTINCT FROM $60 AND
    lv.responsedaycount IS NOT DISTINCT FROM $61 AND
    lv.allowed IS NOT DISTINCT FROM $62 AND
    lv.xvalue IS NOT DISTINCT FROM $63 AND
    lv.yvalue IS NOT DISTINCT FROM $64 AND
    lv.validx IS NOT DISTINCT FROM $65 AND
    lv.validy IS NOT DISTINCT FROM $66 AND
    lv.externalid IS NOT DISTINCT FROM $67 AND
    lv.externalerror IS NOT DISTINCT FROM $68 AND
    lv.pointlocid IS NOT DISTINCT FROM $69 AND
    lv.notified IS NOT DISTINCT FROM $70 AND
    lv.notifieddate IS NOT DISTINCT FROM $71 AND
    lv.scheduled IS NOT DISTINCT FROM $72 AND
    lv.scheduleddate IS NOT DISTINCT FROM $73 AND
    lv.dog IS NOT DISTINCT FROM $74 AND
    lv.schedule_period IS NOT DISTINCT FROM $75 AND
    lv.schedule_notes IS NOT DISTINCT FROM $76 AND
    lv.spanish IS NOT DISTINCT FROM $77 AND
    lv.creationdate IS NOT DISTINCT FROM $78 AND
    lv.creator IS NOT DISTINCT FROM $79 AND
    lv.editdate IS NOT DISTINCT FROM $80 AND
    lv.editor IS NOT DISTINCT FROM $81 AND
    lv.issuesreported IS NOT DISTINCT FROM $82 AND
    lv.jurisdiction IS NOT DISTINCT FROM $83 AND
    lv.notificationtimestamp IS NOT DISTINCT FROM $84 AND
    lv.zone IS NOT DISTINCT FROM $85 AND
    lv.zone2 IS NOT DISTINCT FROM $86
  )
RETURNING *;

-- Example usage: EXECUTE insert_servicerequest_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.SpeciesAbundance
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.speciesabundance_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TABLE fieldseeker.speciesabundance (
  objectid BIGSERIAL NOT NULL,
  trapdata_id UUID,
  species VARCHAR(50),
  males SMALLINT,
  unknown SMALLINT,
  bloodedfem SMALLINT,
  gravidfem SMALLINT,
  larvae SMALLINT,
  poolstogen SMALLINT,
  processed fieldseeker.speciesabundance_notinuit_f_enum,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  pupae SMALLINT,
  eggs SMALLINT,
  females INTEGER,
  total INTEGER,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  yearweek INTEGER,
  globalzscore DOUBLE PRECISION,
  r7score DOUBLE PRECISION,
  r8score DOUBLE PRECISION,
  h3r7 VARCHAR(256),
  h3r8 VARCHAR(256),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.speciesabundance.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.speciesabundance.trapdata_id IS 'Trap Data ID';
COMMENT ON COLUMN fieldseeker.speciesabundance.species IS 'Species';
COMMENT ON COLUMN fieldseeker.speciesabundance.males IS 'Males';
COMMENT ON COLUMN fieldseeker.speciesabundance.unknown IS 'Unknown';
COMMENT ON COLUMN fieldseeker.speciesabundance.bloodedfem IS 'Blooded Females';
COMMENT ON COLUMN fieldseeker.speciesabundance.gravidfem IS 'Gravid Females';
COMMENT ON COLUMN fieldseeker.speciesabundance.larvae IS 'Larvae';
COMMENT ON COLUMN fieldseeker.speciesabundance.poolstogen IS 'Pools to Generate';
COMMENT ON COLUMN fieldseeker.speciesabundance.processed IS 'Processed';
COMMENT ON COLUMN fieldseeker.speciesabundance.pupae IS 'Pupae';
COMMENT ON COLUMN fieldseeker.speciesabundance.eggs IS 'Eggs';
COMMENT ON COLUMN fieldseeker.speciesabundance.females IS 'Females';
COMMENT ON COLUMN fieldseeker.speciesabundance.total IS 'Total Adults';

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.StormDrain
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.stormdrain_stormdrainsymbology_enum AS ENUM (
  'Dry',
  'Needs Treatment',
  'Treated'
);

CREATE TABLE fieldseeker.stormdrain (
  objectid BIGSERIAL NOT NULL,
  nexttreatmentdate TIMESTAMP,
  lasttreatdate TIMESTAMP,
  lastaction VARCHAR(25),
  symbology fieldseeker.stormdrain_stormdrainsymbology_enum,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  laststatus VARCHAR(25),
  zone VARCHAR(25),
  zone2 VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  type VARCHAR(25),
  jurisdiction VARCHAR(25),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.stormdrain.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.stormdrain.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.stormdrain.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.stormdrain.type IS 'Type';
COMMENT ON COLUMN fieldseeker.stormdrain.jurisdiction IS 'Jurisdiction';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_stormdrain_versioned(bigint, timestamp, timestamp, varchar, fieldseeker.stormdrain_stormdrainsymbology_enum, uuid, varchar, timestamp, varchar, timestamp, varchar, varchar, varchar, timestamp, varchar, timestamp, varchar, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.stormdrain
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.stormdrain
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.stormdrain (
  objectid, nexttreatmentdate, lasttreatdate, lastaction, symbology, globalid, created_user, created_date, last_edited_user, last_edited_date, laststatus, zone, zone2, creationdate, creator, editdate, editor, type, jurisdiction,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.nexttreatmentdate IS NOT DISTINCT FROM $2 AND
    lv.lasttreatdate IS NOT DISTINCT FROM $3 AND
    lv.lastaction IS NOT DISTINCT FROM $4 AND
    lv.symbology IS NOT DISTINCT FROM $5 AND
    lv.globalid IS NOT DISTINCT FROM $6 AND
    lv.created_user IS NOT DISTINCT FROM $7 AND
    lv.created_date IS NOT DISTINCT FROM $8 AND
    lv.last_edited_user IS NOT DISTINCT FROM $9 AND
    lv.last_edited_date IS NOT DISTINCT FROM $10 AND
    lv.laststatus IS NOT DISTINCT FROM $11 AND
    lv.zone IS NOT DISTINCT FROM $12 AND
    lv.zone2 IS NOT DISTINCT FROM $13 AND
    lv.creationdate IS NOT DISTINCT FROM $14 AND
    lv.creator IS NOT DISTINCT FROM $15 AND
    lv.editdate IS NOT DISTINCT FROM $16 AND
    lv.editor IS NOT DISTINCT FROM $17 AND
    lv.type IS NOT DISTINCT FROM $18 AND
    lv.jurisdiction IS NOT DISTINCT FROM $19
  )
RETURNING *;

-- Example usage: EXECUTE insert_stormdrain_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.TimeCard
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.timecard_timecard_activity_451e67260c084304a35457170dc13366_enum AS ENUM (
  'Routine inspection',
  'Pre-treatment',
  'Maintenance',
  'ULV',
  'BARRIER',
  'LOGIN',
  'TREATSD',
  'SD',
  'SITEVISIT',
  'ONLINE',
  'SYNC',
  'CREATESR',
  'LC',
  'ACCEPTSR',
  'POINT',
  'DOWNLOAD',
  'COMPLETESR',
  'POLYGON',
  'TRAP',
  'SAMPLE',
  'QA',
  'PTA',
  'FIELDSCOUTING',
  'OFFLINE',
  'LINE',
  'TRAPLOCATION',
  'SAMPLELOCATION',
  'LCLOCATION'
);

CREATE TYPE fieldseeker.timecard_timecardequipmenttype_enum AS ENUM (
  'Spreader',
  'ATV',
  'Truck'
);

CREATE TABLE fieldseeker.timecard (
  objectid BIGSERIAL NOT NULL,
  activity fieldseeker.timecard_timecard_activity_451e67260c084304a35457170dc13366_enum,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  comments VARCHAR(250),
  externalid VARCHAR(25),
  equiptype fieldseeker.timecard_timecardequipmenttype_enum,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  zone2 VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  linelocid UUID,
  pointlocid UUID,
  polygonlocid UUID,
  lclocid UUID,
  samplelocid UUID,
  srid UUID,
  traplocid UUID,
  fieldtech VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  rodentlocid UUID,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.timecard.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.timecard.activity IS 'Activity';
COMMENT ON COLUMN fieldseeker.timecard.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.timecard.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.timecard.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.timecard.equiptype IS 'Equipment Type';
COMMENT ON COLUMN fieldseeker.timecard.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.timecard.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.timecard.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.timecard.fieldtech IS 'Field Tech';

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.TrapData
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.trapdata_mosquitotraptype_enum AS ENUM (
  'GRVD',
  'BGSENT',
  'CO2'
);

CREATE TYPE fieldseeker.trapdata_notinuitrapactivitytype_enum AS ENUM (
  'S',
  'R'
);

CREATE TYPE fieldseeker.trapdata_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.trapdata_mosquitositecondition_enum AS ENUM (
  'Dry',
  'Clean',
  'Full',
  'Low'
);

CREATE TYPE fieldseeker.trapdata_mosquitotrapcondition_enum AS ENUM (
  'Damaged',
  'Missing',
  'Fan Off',
  'Fan Slow'
);

CREATE TYPE fieldseeker.trapdata_trapdata_winddir_c1a31e05_d0b9_4b22_8800_be127bb3f166_enum AS ENUM (
  'E',
  'N',
  'NE',
  'NW',
  'S',
  'SE',
  'SW',
  'W'
);

CREATE TYPE fieldseeker.trapdata_trapdata_lure_25fe542f_077f_4254_8681_76e8f436354b_enum AS ENUM (
  'CO2 (Dry Ice)',
  'CO2 (Sugar Yeast)',
  'BG-Lure',
  'Gravid Water'
);

CREATE TABLE fieldseeker.trapdata (
  objectid BIGSERIAL NOT NULL,
  traptype fieldseeker.trapdata_mosquitotraptype_enum,
  trapactivitytype fieldseeker.trapdata_notinuitrapactivitytype_enum,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  comments VARCHAR(250),
  idbytech VARCHAR(25),
  sortbytech VARCHAR(25),
  processed fieldseeker.trapdata_notinuit_f_enum,
  sitecond fieldseeker.trapdata_mosquitositecondition_enum,
  locationname VARCHAR(25),
  recordstatus SMALLINT,
  reviewed fieldseeker.trapdata_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  trapcondition fieldseeker.trapdata_mosquitotrapcondition_enum,
  trapnights SMALLINT,
  zone VARCHAR(25),
  zone2 VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  srid UUID,
  fieldtech VARCHAR(25),
  gatewaysync SMALLINT,
  loc_id UUID,
  voltage DOUBLE PRECISION,
  winddir fieldseeker.trapdata_trapdata_winddir_c1a31e05_d0b9_4b22_8800_be127bb3f166_enum,
  windspeed DOUBLE PRECISION,
  avetemp DOUBLE PRECISION,
  raingauge DOUBLE PRECISION,
  lr SMALLINT,
  field INTEGER,
  vectorsurvtrapdataid VARCHAR(50),
  vectorsurvtraplocationid VARCHAR(50),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  lure fieldseeker.trapdata_trapdata_lure_25fe542f_077f_4254_8681_76e8f436354b_enum,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.trapdata.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.trapdata.traptype IS 'Trap Type';
COMMENT ON COLUMN fieldseeker.trapdata.trapactivitytype IS 'Trap Activity Type';
COMMENT ON COLUMN fieldseeker.trapdata.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.trapdata.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.trapdata.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.trapdata.idbytech IS 'Tech Identifying Species in Lab';
COMMENT ON COLUMN fieldseeker.trapdata.sortbytech IS 'Tech Sorting Trap Results in Lab';
COMMENT ON COLUMN fieldseeker.trapdata.processed IS 'Processed';
COMMENT ON COLUMN fieldseeker.trapdata.sitecond IS 'Site Conditions';
COMMENT ON COLUMN fieldseeker.trapdata.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.trapdata.recordstatus IS 'RecordStatus';
COMMENT ON COLUMN fieldseeker.trapdata.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.trapdata.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.trapdata.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.trapdata.trapcondition IS 'Trap Condition';
COMMENT ON COLUMN fieldseeker.trapdata.trapnights IS 'Trap Nights';
COMMENT ON COLUMN fieldseeker.trapdata.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.trapdata.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.trapdata.fieldtech IS 'Field Tech';
COMMENT ON COLUMN fieldseeker.trapdata.gatewaysync IS 'Gateway Sync';
COMMENT ON COLUMN fieldseeker.trapdata.voltage IS 'Voltage';
COMMENT ON COLUMN fieldseeker.trapdata.lr IS 'Landing Rate';

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.TrapLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.traplocation_traplocation_priority_680fb011063b41d59f39271c959b857f_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None',
  'Project',
  'Fixed',
  'Response '
);

CREATE TYPE fieldseeker.traplocation_traplocation_usetype_5e0eff9231fb404c98cc53c1d49a2193_enum AS ENUM (
  'Fixed Trapping',
  'Response Trapping',
  'Service Request',
  'Project Trap'
);

CREATE TYPE fieldseeker.traplocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.traplocation_traplocation_accessdesc_154cbd10_4524_4e3a_8ca0_f099ec86556a_enum AS ENUM (
  'homeowner preference',
  'no longer needed'
);

CREATE TYPE fieldseeker.traplocation_traplocation_habitat_5c349680f5ff40b1aeca88c17993e8f3_enum AS ENUM (
  'Trap'
);

CREATE TABLE fieldseeker.traplocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.traplocation_traplocation_habitat_5c349680f5ff40b1aeca88c17993e8f3_enum,
  priority fieldseeker.traplocation_traplocation_priority_680fb011063b41d59f39271c959b857f_enum,
  usetype fieldseeker.traplocation_traplocation_usetype_5e0eff9231fb404c98cc53c1d49a2193_enum,
  active fieldseeker.traplocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc fieldseeker.traplocation_traplocation_accessdesc_154cbd10_4524_4e3a_8ca0_f099ec86556a_enum,
  comments VARCHAR(250),
  externalid VARCHAR(50),
  nextactiondatescheduled TIMESTAMP,
  zone2 VARCHAR(25),
  locationnumber INTEGER,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  gatewaysync SMALLINT,
  route INTEGER,
  set_dow INTEGER,
  route_order INTEGER,
  vectorsurvsiteid VARCHAR(50),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  h3r7 VARCHAR(255),
  h3r8 VARCHAR(255),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.traplocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.traplocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.traplocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.traplocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.traplocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.traplocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.traplocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.traplocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.traplocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.traplocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.traplocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.traplocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.traplocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.traplocation.gatewaysync IS 'Gateway Sync';
COMMENT ON COLUMN fieldseeker.traplocation.route IS 'Route';
COMMENT ON COLUMN fieldseeker.traplocation.set_dow IS 'Set Day of Week';
COMMENT ON COLUMN fieldseeker.traplocation.route_order IS 'Route order';

-- Field active has default value: 1

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.Treatment
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.treatment_treatment_sitecond_5a15bf36fa124280b961f31cd1a9b571_enum AS ENUM (
  'Dry',
  'Flowing',
  'Maintained',
  'Unmaintained',
  'High Organic',
  'Fish Present',
  'Stagnant'
);

CREATE TYPE fieldseeker.treatment_mosquitoactivity_enum AS ENUM (
  'Routine inspection',
  'Pre-treatment',
  'Maintenance',
  'ULV',
  'BARRIER',
  'LOGIN',
  'TREATSD',
  'SD',
  'SITEVISIT',
  'ONLINE',
  'SYNC',
  'CREATESR',
  'LC',
  'ACCEPTSR',
  'POINT',
  'DOWNLOAD',
  'COMPLETESR',
  'POLYGON',
  'TRAP',
  'SAMPLE',
  'QA',
  'PTA',
  'FIELDSCOUTING',
  'OFFLINE',
  'LINE',
  'TRAPLOCATION',
  'SAMPLELOCATION',
  'LCLOCATION'
);

CREATE TYPE fieldseeker.treatment_mosquitoproductmeasureunit_enum AS ENUM (
  'briquet',
  'dry oz',
  'each',
  'fl oz',
  'gal',
  'lb',
  'packet',
  'pouch'
);

CREATE TYPE fieldseeker.treatment_treatment_method_d558ca3ccf43440c8160758253967621_enum AS ENUM (
  'Argo',
  'ATV',
  'Backpack',
  'Drone',
  'Manual',
  'Truck',
  'ULV',
  'WALS',
  'Administrative Action'
);

CREATE TYPE fieldseeker.treatment_treatment_equiptype_45694d79_ff21_42cc_be4f_a0d1def4fba0_enum AS ENUM (
  'Backpack #1',
  'A1  Mist Sprayer (T-3) ',
  'Spreader #2',
  'Guardian #73 ',
  'ULV #74 (Grizzly)',
  'Clark ULV Sprayer #71',
  'Clark ULV Sprayer #72',
  'Spray bottle'
);

CREATE TYPE fieldseeker.treatment_notinuiwinddirection_enum AS ENUM (
  'N',
  'NE',
  'E',
  'SE',
  'S',
  'SW',
  'W',
  'NW'
);

CREATE TYPE fieldseeker.treatment_treatment_sitecond_f812e1f64dcb4dc9a75da9d00abe6169_enum AS ENUM (
  'Dry',
  'Flowing',
  'Maintained',
  'Unmaintained',
  'High Organic',
  'Fish Present'
);

CREATE TYPE fieldseeker.treatment_mosquitoproductareaunit_enum AS ENUM (
  'acre',
  'sq ft'
);

CREATE TYPE fieldseeker.treatment_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.treatment_treatment_habitat_0afee7eb_f9ea_4707_8483_cccfe60f0d16_enum AS ENUM (
  'orchard',
  'row_crops',
  'vine_crops',
  'ag_grass_or_grain',
  'pasture',
  'irrigation_standpipe',
  'ditch',
  'pond',
  'sump',
  'drain',
  'dairy_lagoon',
  'wastewater_treatment',
  'trough',
  'depression',
  'gutter',
  'rain_gutter',
  'culvert',
  'Utility',
  'catch_basin',
  'stream_or_creek',
  'slough',
  'river',
  'marsh_or_wetlands',
  'containers',
  'watering_bowl',
  'plant_saucer',
  'yard_drain',
  'plant_axil',
  'treehole',
  'foutain_or_water_feature',
  'bird_bath',
  'misc_water_accumulation',
  'tarp_or_cover',
  'swimming_pool',
  'aboveground_pool',
  'kid_pool',
  'hot_tub',
  'applicance',
  'flooded_structure',
  'low_point'
);

CREATE TABLE fieldseeker.treatment (
  objectid BIGSERIAL NOT NULL,
  activity fieldseeker.treatment_mosquitoactivity_enum,
  treatarea DOUBLE PRECISION,
  areaunit fieldseeker.treatment_mosquitoproductareaunit_enum,
  product VARCHAR(25),
  qty DOUBLE PRECISION,
  qtyunit fieldseeker.treatment_mosquitoproductmeasureunit_enum,
  method fieldseeker.treatment_treatment_method_d558ca3ccf43440c8160758253967621_enum,
  equiptype fieldseeker.treatment_treatment_equiptype_45694d79_ff21_42cc_be4f_a0d1def4fba0_enum,
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  winddir fieldseeker.treatment_notinuiwinddirection_enum,
  raingauge DOUBLE PRECISION,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  insp_id UUID,
  reviewed fieldseeker.treatment_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  warningoverride fieldseeker.treatment_notinuit_f_enum,
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  treatacres DOUBLE PRECISION,
  tirecount SMALLINT,
  cbcount SMALLINT,
  containercount SMALLINT,
  globalid UUID,
  treatmentlength DOUBLE PRECISION,
  treatmenthours DOUBLE PRECISION,
  treatmentlengthunits VARCHAR(5),
  linelocid UUID,
  pointlocid UUID,
  polygonlocid UUID,
  srid UUID,
  sdid UUID,
  barrierrouteid UUID,
  ulvrouteid UUID,
  fieldtech VARCHAR(25),
  ptaid UUID,
  flowrate DOUBLE PRECISION,
  habitat fieldseeker.treatment_treatment_habitat_0afee7eb_f9ea_4707_8483_cccfe60f0d16_enum,
  treathectares DOUBLE PRECISION,
  invloc VARCHAR(25),
  temp_sitecond fieldseeker.treatment_treatment_sitecond_f812e1f64dcb4dc9a75da9d00abe6169_enum,
  sitecond fieldseeker.treatment_treatment_sitecond_5a15bf36fa124280b961f31cd1a9b571_enum,
  totalcostprodcut DOUBLE PRECISION,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  targetspecies VARCHAR(250),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.treatment.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.treatment.activity IS 'Activity';
COMMENT ON COLUMN fieldseeker.treatment.treatarea IS 'Area Treated';
COMMENT ON COLUMN fieldseeker.treatment.areaunit IS 'Area Unit';
COMMENT ON COLUMN fieldseeker.treatment.product IS 'Product';
COMMENT ON COLUMN fieldseeker.treatment.qty IS 'Quantity';
COMMENT ON COLUMN fieldseeker.treatment.qtyunit IS 'Quantity Unit';
COMMENT ON COLUMN fieldseeker.treatment.method IS 'Method';
COMMENT ON COLUMN fieldseeker.treatment.equiptype IS 'Equipment Type';
COMMENT ON COLUMN fieldseeker.treatment.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.treatment.avetemp IS 'Average Temperature';
COMMENT ON COLUMN fieldseeker.treatment.windspeed IS 'Wind Speed';
COMMENT ON COLUMN fieldseeker.treatment.winddir IS 'Wind Direction';
COMMENT ON COLUMN fieldseeker.treatment.raingauge IS 'Rain Gauge';
COMMENT ON COLUMN fieldseeker.treatment.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.treatment.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.treatment.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.treatment.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.treatment.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.treatment.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.treatment.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.treatment.warningoverride IS 'Warning Override';
COMMENT ON COLUMN fieldseeker.treatment.recordstatus IS 'RecordStatus';
COMMENT ON COLUMN fieldseeker.treatment.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.treatment.treatacres IS 'Treated Acres';
COMMENT ON COLUMN fieldseeker.treatment.tirecount IS 'Tire Count';
COMMENT ON COLUMN fieldseeker.treatment.cbcount IS 'Catch Basin Count';
COMMENT ON COLUMN fieldseeker.treatment.containercount IS 'Container Count';
COMMENT ON COLUMN fieldseeker.treatment.treatmentlength IS 'Treatment Length';
COMMENT ON COLUMN fieldseeker.treatment.treatmenthours IS 'Treatment Hours';
COMMENT ON COLUMN fieldseeker.treatment.treatmentlengthunits IS 'Treatment Length Units';
COMMENT ON COLUMN fieldseeker.treatment.fieldtech IS 'Field Tech';
COMMENT ON COLUMN fieldseeker.treatment.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.treatment.treathectares IS 'Treat Hectares';
COMMENT ON COLUMN fieldseeker.treatment.invloc IS 'Inventory Location';
COMMENT ON COLUMN fieldseeker.treatment.temp_sitecond IS 'temp_Conditions';
COMMENT ON COLUMN fieldseeker.treatment.sitecond IS 'Conditions';
COMMENT ON COLUMN fieldseeker.treatment.totalcostprodcut IS 'TotalCostProduct';
COMMENT ON COLUMN fieldseeker.treatment.targetspecies IS 'Target Species';

-- Prepared statement for conditional insert with versioning
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
-- Table definition for fieldseeker.TreatmentArea
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.treatmentarea (
  objectid BIGSERIAL NOT NULL,
  treat_id UUID,
  session_id UUID,
  treatdate TIMESTAMP,
  comments VARCHAR(250),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  notified SMALLINT,
  type VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.treatmentarea.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.treatmentarea.treatdate IS 'Treatment Date';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_treatmentarea_versioned(bigint, uuid, uuid, timestamp, varchar, uuid, varchar, timestamp, varchar, timestamp, smallint, varchar, timestamp, varchar, timestamp, varchar, double precision, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.treatmentarea
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.treatmentarea
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.treatmentarea (
  objectid, treat_id, session_id, treatdate, comments, globalid, created_user, created_date, last_edited_user, last_edited_date, notified, type, creationdate, creator, editdate, editor, shape__area, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.treat_id IS NOT DISTINCT FROM $2 AND
    lv.session_id IS NOT DISTINCT FROM $3 AND
    lv.treatdate IS NOT DISTINCT FROM $4 AND
    lv.comments IS NOT DISTINCT FROM $5 AND
    lv.globalid IS NOT DISTINCT FROM $6 AND
    lv.created_user IS NOT DISTINCT FROM $7 AND
    lv.created_date IS NOT DISTINCT FROM $8 AND
    lv.last_edited_user IS NOT DISTINCT FROM $9 AND
    lv.last_edited_date IS NOT DISTINCT FROM $10 AND
    lv.notified IS NOT DISTINCT FROM $11 AND
    lv.type IS NOT DISTINCT FROM $12 AND
    lv.creationdate IS NOT DISTINCT FROM $13 AND
    lv.creator IS NOT DISTINCT FROM $14 AND
    lv.editdate IS NOT DISTINCT FROM $15 AND
    lv.editor IS NOT DISTINCT FROM $16 AND
    lv.shape__area IS NOT DISTINCT FROM $17 AND
    lv.shape__length IS NOT DISTINCT FROM $18
  )
RETURNING *;

-- Example usage: EXECUTE insert_treatmentarea_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.Zones
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.zones (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(50),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  active INTEGER,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.zones.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.zones.name IS 'Name';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_zones_versioned(bigint, varchar, uuid, varchar, timestamp, varchar, timestamp, integer, timestamp, varchar, timestamp, varchar, double precision, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.zones
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.zones
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.zones (
  objectid, name, globalid, created_user, created_date, last_edited_user, last_edited_date, active, creationdate, creator, editdate, editor, shape__area, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.name IS NOT DISTINCT FROM $2 AND
    lv.globalid IS NOT DISTINCT FROM $3 AND
    lv.created_user IS NOT DISTINCT FROM $4 AND
    lv.created_date IS NOT DISTINCT FROM $5 AND
    lv.last_edited_user IS NOT DISTINCT FROM $6 AND
    lv.last_edited_date IS NOT DISTINCT FROM $7 AND
    lv.active IS NOT DISTINCT FROM $8 AND
    lv.creationdate IS NOT DISTINCT FROM $9 AND
    lv.creator IS NOT DISTINCT FROM $10 AND
    lv.editdate IS NOT DISTINCT FROM $11 AND
    lv.editor IS NOT DISTINCT FROM $12 AND
    lv.shape__area IS NOT DISTINCT FROM $13 AND
    lv.shape__length IS NOT DISTINCT FROM $14
  )
RETURNING *;

-- Example usage: EXECUTE insert_zones_versioned(id, value1, value2, ...);
-- Table definition for fieldseeker.Zones2
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.zones2 (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(50),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.zones2.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.zones2.name IS 'Name';

-- Prepared statement for conditional insert with versioning
-- Only inserts a new version if data has changed
PREPARE insert_zones2_versioned(bigint, varchar, uuid, varchar, timestamp, varchar, timestamp, timestamp, varchar, timestamp, varchar, double precision, double precision) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.zones2
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.zones2
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.zones2 (
  objectid, name, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, shape__area, shape__length,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.name IS NOT DISTINCT FROM $2 AND
    lv.globalid IS NOT DISTINCT FROM $3 AND
    lv.created_user IS NOT DISTINCT FROM $4 AND
    lv.created_date IS NOT DISTINCT FROM $5 AND
    lv.last_edited_user IS NOT DISTINCT FROM $6 AND
    lv.last_edited_date IS NOT DISTINCT FROM $7 AND
    lv.creationdate IS NOT DISTINCT FROM $8 AND
    lv.creator IS NOT DISTINCT FROM $9 AND
    lv.editdate IS NOT DISTINCT FROM $10 AND
    lv.editor IS NOT DISTINCT FROM $11 AND
    lv.shape__area IS NOT DISTINCT FROM $12 AND
    lv.shape__length IS NOT DISTINCT FROM $13
  )
RETURNING *;

-- Example usage: EXECUTE insert_zones2_versioned(id, value1, value2, ...);
-- +goose Down
DROP SCHEMA fieldseeker CASCADE;
