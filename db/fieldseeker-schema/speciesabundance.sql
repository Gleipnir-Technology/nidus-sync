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

-- See insert/insert_speciesabundance_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
