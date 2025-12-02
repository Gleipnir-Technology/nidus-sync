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

-- See insert/insert_samplelocation_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
