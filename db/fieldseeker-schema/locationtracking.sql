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

-- See insert/insert_locationtracking_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
