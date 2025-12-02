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

-- See insert/insert_stormdrain_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
