-- Table definition for fieldseeker.TrapLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.traplocation_traplocation_habitat_5c349680f5ff40b1aeca88c17993e8f3_enum AS ENUM (
  'Trap'
);

CREATE TYPE fieldseeker.traplocation_traplocation_priority_680fb011063b41d59f39271c959b857f_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None',
  'Project',
  'Fixed',
  'Response '
);

CREATE TYPE fieldseeker.traplocation_traplocation_usetype_5e0eff9231fb404c98cc53c1d49a2193_enum AS ENUM (
  'Fixed Trapping',
  'Response Trapping',
  'Service Request',
  'Project Trap'
);

CREATE TYPE fieldseeker.traplocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.traplocation_traplocation_accessdesc_154cbd10_4524_4e3a_8ca0_f099ec86556a_enum AS ENUM (
  'homeowner preference',
  'no longer needed'
);

CREATE TABLE fieldseeker.traplocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.traplocation_traplocation_habitat_5c349680f5ff40b1aeca88c17993e8f3_enum,
  priority fieldseeker.traplocation_traplocation_priority_680fb011063b41d59f39271c959b857f_enum,
  usetype fieldseeker.traplocation_traplocation_usetype_5e0eff9231fb404c98cc53c1d49a2193_enum,
  active fieldseeker.traplocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc fieldseeker.traplocation_traplocation_accessdesc_154cbd10_4524_4e3a_8ca0_f099ec86556a_enum,
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
  route INTEGER,
  set_dow INTEGER,
  route_order INTEGER,
  vectorsurvsiteid VARCHAR(50),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  h3r7 VARCHAR(255),
  h3r8 VARCHAR(255),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.traplocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.traplocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.traplocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.traplocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.traplocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.traplocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.traplocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.traplocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.traplocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.traplocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.traplocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.traplocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.traplocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.traplocation.gatewaysync IS 'Gateway Sync';
COMMENT ON COLUMN fieldseeker.traplocation.route IS 'Route';
COMMENT ON COLUMN fieldseeker.traplocation.set_dow IS 'Set Day of Week';
COMMENT ON COLUMN fieldseeker.traplocation.route_order IS 'Route order';

-- Field active has default value: 1

-- See insert/insert_traplocation_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
