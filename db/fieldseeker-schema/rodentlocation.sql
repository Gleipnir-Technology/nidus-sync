-- Table definition for fieldseeker.RodentLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.rodentlocation_locationusetype_1_enum AS ENUM (
  'Residential',
  'Commercial',
  'Industrial',
  'Agricultural',
  'Mixed use'
);

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

-- Field active has default value: 1

-- Field symbology has default value: 'NONE'

-- See insert/insert_rodentlocation_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
