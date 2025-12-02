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

-- See insert/insert_pool_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
