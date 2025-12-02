-- Table definition for fieldseeker.SampleCollection
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.samplecollection_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.samplecollection_mosquitosampletype_enum AS ENUM (
  'Blood',
  'Tissue',
  'Specimen',
  'Carcass'
);

CREATE TYPE fieldseeker.samplecollection_mosquitosamplecondition_enum AS ENUM (
  'Usable',
  'Unusable'
);

CREATE TYPE fieldseeker.samplecollection_mosquitosamplespecies_enum AS ENUM (
  'Chicken',
  'Wild bird',
  'Horse',
  'Human'
);

CREATE TYPE fieldseeker.samplecollection_notinuisex_enum AS ENUM (
  'M',
  'F',
  'U'
);

CREATE TYPE fieldseeker.samplecollection_mosquitoactivity_enum AS ENUM (
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

CREATE TYPE fieldseeker.samplecollection_mosquitolabname_enum AS ENUM (
  'Internal Lab',
  'State Lab'
);

CREATE TYPE fieldseeker.samplecollection_mosquitositecondition_enum AS ENUM (
  'Dry',
  'Clean',
  'Full',
  'Low'
);

CREATE TYPE fieldseeker.samplecollection_notinuiwinddirection_enum AS ENUM (
  'N',
  'NE',
  'E',
  'SE',
  'S',
  'SW',
  'W',
  'NW'
);

CREATE TYPE fieldseeker.samplecollection_mosquitotestmethod_enum AS ENUM (
  'RAMP',
  'VecTest',
  'ELISA',
  'RT-PCR'
);

CREATE TYPE fieldseeker.samplecollection_mosquitodisease_enum AS ENUM (
  'EEE',
  'WNV',
  'Dengue',
  'Zika'
);

CREATE TABLE fieldseeker.samplecollection (
  objectid BIGSERIAL NOT NULL,
  loc_id UUID,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  sitecond fieldseeker.samplecollection_mosquitositecondition_enum,
  sampleid VARCHAR(25),
  survtech VARCHAR(25),
  datesent TIMESTAMP,
  datetested TIMESTAMP,
  testtech VARCHAR(25),
  comments VARCHAR(250),
  processed fieldseeker.samplecollection_notinuit_f_enum,
  sampletype fieldseeker.samplecollection_mosquitosampletype_enum,
  samplecond fieldseeker.samplecollection_mosquitosamplecondition_enum,
  species fieldseeker.samplecollection_mosquitosamplespecies_enum,
  sex fieldseeker.samplecollection_notinuisex_enum,
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  winddir fieldseeker.samplecollection_notinuiwinddirection_enum,
  raingauge DOUBLE PRECISION,
  activity fieldseeker.samplecollection_mosquitoactivity_enum,
  testmethod fieldseeker.samplecollection_mosquitotestmethod_enum,
  diseasetested fieldseeker.samplecollection_mosquitodisease_enum,
  diseasepos fieldseeker.samplecollection_mosquitodisease_enum,
  reviewed fieldseeker.samplecollection_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  lab fieldseeker.samplecollection_mosquitolabname_enum,
  fieldtech VARCHAR(25),
  flockid UUID,
  samplecount SMALLINT,
  chickenid UUID,
  gatewaysync SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.samplecollection.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.samplecollection.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.samplecollection.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.samplecollection.sitecond IS 'Conditions';
COMMENT ON COLUMN fieldseeker.samplecollection.sampleid IS 'Sample ID';
COMMENT ON COLUMN fieldseeker.samplecollection.survtech IS 'Surveillance Technician';
COMMENT ON COLUMN fieldseeker.samplecollection.datesent IS 'Sent';
COMMENT ON COLUMN fieldseeker.samplecollection.datetested IS 'Tested';
COMMENT ON COLUMN fieldseeker.samplecollection.testtech IS 'Test Technician';
COMMENT ON COLUMN fieldseeker.samplecollection.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.samplecollection.processed IS 'Processed';
COMMENT ON COLUMN fieldseeker.samplecollection.sampletype IS 'Sample Type';
COMMENT ON COLUMN fieldseeker.samplecollection.samplecond IS 'Sample Condition';
COMMENT ON COLUMN fieldseeker.samplecollection.species IS 'Species';
COMMENT ON COLUMN fieldseeker.samplecollection.sex IS 'Sex';
COMMENT ON COLUMN fieldseeker.samplecollection.avetemp IS 'Average Temperature';
COMMENT ON COLUMN fieldseeker.samplecollection.windspeed IS 'Wind Speed';
COMMENT ON COLUMN fieldseeker.samplecollection.winddir IS 'Wind Direction';
COMMENT ON COLUMN fieldseeker.samplecollection.raingauge IS 'Rain Gauge';
COMMENT ON COLUMN fieldseeker.samplecollection.activity IS 'Activity';
COMMENT ON COLUMN fieldseeker.samplecollection.testmethod IS 'Test Method';
COMMENT ON COLUMN fieldseeker.samplecollection.diseasetested IS 'Disease Tested';
COMMENT ON COLUMN fieldseeker.samplecollection.diseasepos IS 'Disease Positive';
COMMENT ON COLUMN fieldseeker.samplecollection.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.samplecollection.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.samplecollection.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.samplecollection.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.samplecollection.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.samplecollection.recordstatus IS 'RecordStatus';
COMMENT ON COLUMN fieldseeker.samplecollection.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.samplecollection.lab IS 'Lab';
COMMENT ON COLUMN fieldseeker.samplecollection.fieldtech IS 'Field Tech';
COMMENT ON COLUMN fieldseeker.samplecollection.samplecount IS 'Sample Count';
COMMENT ON COLUMN fieldseeker.samplecollection.chickenid IS 'ChickenID';
COMMENT ON COLUMN fieldseeker.samplecollection.gatewaysync IS 'Gateway Sync';

-- See insert/insert_samplecollection_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
