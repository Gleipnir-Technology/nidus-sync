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

-- See insert/insert_inspectionsample_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
