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

-- See insert/insert_pooldetail_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
