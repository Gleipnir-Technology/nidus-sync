-- Table definition for fieldseeker.TrapData
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.trapdata_mosquitotraptype_enum AS ENUM (
  'GRVD',
  'BGSENT',
  'CO2'
);

CREATE TYPE fieldseeker.trapdata_notinuitrapactivitytype_enum AS ENUM (
  'S',
  'R'
);

CREATE TYPE fieldseeker.trapdata_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.trapdata_mosquitositecondition_enum AS ENUM (
  'Dry',
  'Clean',
  'Full',
  'Low'
);

CREATE TYPE fieldseeker.trapdata_mosquitotrapcondition_enum AS ENUM (
  'Damaged',
  'Missing',
  'Fan Off',
  'Fan Slow'
);

CREATE TYPE fieldseeker.trapdata_trapdata_winddir_c1a31e05_d0b9_4b22_8800_be127bb3f166_enum AS ENUM (
  'E',
  'N',
  'NE',
  'NW',
  'S',
  'SE',
  'SW',
  'W'
);

CREATE TYPE fieldseeker.trapdata_trapdata_lure_25fe542f_077f_4254_8681_76e8f436354b_enum AS ENUM (
  'CO2 (Dry Ice)',
  'CO2 (Sugar Yeast)',
  'BG-Lure',
  'Gravid Water'
);

CREATE TABLE fieldseeker.trapdata (
  objectid BIGSERIAL NOT NULL,
  traptype fieldseeker.trapdata_mosquitotraptype_enum,
  trapactivitytype fieldseeker.trapdata_notinuitrapactivitytype_enum,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  comments VARCHAR(250),
  idbytech VARCHAR(25),
  sortbytech VARCHAR(25),
  processed fieldseeker.trapdata_notinuit_f_enum,
  sitecond fieldseeker.trapdata_mosquitositecondition_enum,
  locationname VARCHAR(25),
  recordstatus SMALLINT,
  reviewed fieldseeker.trapdata_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  trapcondition fieldseeker.trapdata_mosquitotrapcondition_enum,
  trapnights SMALLINT,
  zone VARCHAR(25),
  zone2 VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  srid UUID,
  fieldtech VARCHAR(25),
  gatewaysync SMALLINT,
  loc_id UUID,
  voltage DOUBLE PRECISION,
  winddir fieldseeker.trapdata_trapdata_winddir_c1a31e05_d0b9_4b22_8800_be127bb3f166_enum,
  windspeed DOUBLE PRECISION,
  avetemp DOUBLE PRECISION,
  raingauge DOUBLE PRECISION,
  lr SMALLINT,
  field INTEGER,
  vectorsurvtrapdataid VARCHAR(50),
  vectorsurvtraplocationid VARCHAR(50),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  lure fieldseeker.trapdata_trapdata_lure_25fe542f_077f_4254_8681_76e8f436354b_enum,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.trapdata.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.trapdata.traptype IS 'Trap Type';
COMMENT ON COLUMN fieldseeker.trapdata.trapactivitytype IS 'Trap Activity Type';
COMMENT ON COLUMN fieldseeker.trapdata.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.trapdata.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.trapdata.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.trapdata.idbytech IS 'Tech Identifying Species in Lab';
COMMENT ON COLUMN fieldseeker.trapdata.sortbytech IS 'Tech Sorting Trap Results in Lab';
COMMENT ON COLUMN fieldseeker.trapdata.processed IS 'Processed';
COMMENT ON COLUMN fieldseeker.trapdata.sitecond IS 'Site Conditions';
COMMENT ON COLUMN fieldseeker.trapdata.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.trapdata.recordstatus IS 'RecordStatus';
COMMENT ON COLUMN fieldseeker.trapdata.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.trapdata.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.trapdata.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.trapdata.trapcondition IS 'Trap Condition';
COMMENT ON COLUMN fieldseeker.trapdata.trapnights IS 'Trap Nights';
COMMENT ON COLUMN fieldseeker.trapdata.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.trapdata.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.trapdata.fieldtech IS 'Field Tech';
COMMENT ON COLUMN fieldseeker.trapdata.gatewaysync IS 'Gateway Sync';
COMMENT ON COLUMN fieldseeker.trapdata.voltage IS 'Voltage';
COMMENT ON COLUMN fieldseeker.trapdata.lr IS 'Landing Rate';

-- See insert/insert_trapdata_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
