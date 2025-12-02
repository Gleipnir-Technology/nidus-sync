-- Table definition for fieldseeker.TimeCard
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.timecard_timecardequipmenttype_enum AS ENUM (
  'Spreader',
  'ATV',
  'Truck'
);

CREATE TYPE fieldseeker.timecard_timecard_activity_451e67260c084304a35457170dc13366_enum AS ENUM (
  'Routine inspection',
  'Pre-treatment',
  'Maintenance',
  'ULV',
  'BARRIER',
  'LOGIN',
  'TREATSD',
  'SD',
  'SITEVISIT',
  'ONLINE',
  'SYNC',
  'CREATESR',
  'LC',
  'ACCEPTSR',
  'POINT',
  'DOWNLOAD',
  'COMPLETESR',
  'POLYGON',
  'TRAP',
  'SAMPLE',
  'QA',
  'PTA',
  'FIELDSCOUTING',
  'OFFLINE',
  'LINE',
  'TRAPLOCATION',
  'SAMPLELOCATION',
  'LCLOCATION'
);

CREATE TABLE fieldseeker.timecard (
  objectid BIGSERIAL NOT NULL,
  activity fieldseeker.timecard_timecard_activity_451e67260c084304a35457170dc13366_enum,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  comments VARCHAR(250),
  externalid VARCHAR(25),
  equiptype fieldseeker.timecard_timecardequipmenttype_enum,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  zone2 VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  linelocid UUID,
  pointlocid UUID,
  polygonlocid UUID,
  lclocid UUID,
  samplelocid UUID,
  srid UUID,
  traplocid UUID,
  fieldtech VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  rodentlocid UUID,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.timecard.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.timecard.activity IS 'Activity';
COMMENT ON COLUMN fieldseeker.timecard.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.timecard.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.timecard.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.timecard.equiptype IS 'Equipment Type';
COMMENT ON COLUMN fieldseeker.timecard.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.timecard.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.timecard.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.timecard.fieldtech IS 'Field Tech';

-- See insert/insert_timecard_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
