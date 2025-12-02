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

-- See insert/insert_fieldscoutinglog_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
