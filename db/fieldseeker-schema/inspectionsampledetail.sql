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

-- See insert/insert_inspectionsampledetail_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
