-- Table definition for fieldseeker.MosquitoInspection
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoinspection_adminaction_b74ae1bb_c98b_40f6_8cfa_40e4fd16c270_enum AS ENUM (
  'yes',
  'no'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoactivity_enum AS ENUM (
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

CREATE TYPE fieldseeker.mosquitoinspection_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitofieldspecies_enum AS ENUM (
  'Aedes',
  'Culex'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitobreeding_enum AS ENUM (
  'None',
  'Light',
  'Moderate',
  'Intense'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoadultactivity_enum AS ENUM (
  'None',
  'Light',
  'Moderate',
  'Intense'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoinspection_domstage_b7a6c36bccde49a292020de4812cf5ae_enum AS ENUM (
  '1',
  '2',
  '3',
  '4',
  '5',
  '1-2',
  '2-3',
  '3-4'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoinspection_actiontaken_252243d69b0b44ddbdc229c04ec3a8d5_enum AS ENUM (
  'Treatment',
  'Mechanical or Biological Treatment',
  'Resident Schedule Request',
  'Administrative'
);

CREATE TYPE fieldseeker.mosquitoinspection_notinuiwinddirection_enum AS ENUM (
  'N',
  'NE',
  'E',
  'SE',
  'S',
  'SW',
  'W',
  'NW'
);

CREATE TYPE fieldseeker.mosquitoinspection_mosquitoinspection_sitecond_db7350bc_81e5_401e_858f_cd3e5e5d8a34_enum AS ENUM (
  'Dry',
  'Flowing',
  'Maintained',
  'Unmaintained',
  'High Organic',
  'Unknown',
  'Stagnant',
  'Needs Monitoring',
  'Drying Out',
  'Appears Vacant',
  'Entry Denied',
  'Pool Removed',
  'False Pool'
);

CREATE TABLE fieldseeker.mosquitoinspection (
  objectid BIGSERIAL NOT NULL,
  numdips SMALLINT,
  activity fieldseeker.mosquitoinspection_mosquitoactivity_enum,
  breeding fieldseeker.mosquitoinspection_mosquitobreeding_enum,
  totlarvae SMALLINT,
  totpupae SMALLINT,
  eggs SMALLINT,
  posdips SMALLINT,
  adultact fieldseeker.mosquitoinspection_mosquitoadultactivity_enum,
  lstages VARCHAR(25),
  domstage fieldseeker.mosquitoinspection_mosquitoinspection_domstage_b7a6c36bccde49a292020de4812cf5ae_enum,
  actiontaken fieldseeker.mosquitoinspection_mosquitoinspection_actiontaken_252243d69b0b44ddbdc229c04ec3a8d5_enum,
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  raingauge DOUBLE PRECISION,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  winddir fieldseeker.mosquitoinspection_notinuiwinddirection_enum,
  avglarvae DOUBLE PRECISION,
  avgpupae DOUBLE PRECISION,
  reviewed fieldseeker.mosquitoinspection_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  personalcontact fieldseeker.mosquitoinspection_notinuit_f_enum,
  tirecount SMALLINT,
  cbcount SMALLINT,
  containercount SMALLINT,
  fieldspecies fieldseeker.mosquitoinspection_mosquitofieldspecies_enum,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  linelocid UUID,
  pointlocid UUID,
  polygonlocid UUID,
  srid UUID,
  fieldtech VARCHAR(25),
  larvaepresent fieldseeker.mosquitoinspection_notinuit_f_enum,
  pupaepresent fieldseeker.mosquitoinspection_notinuit_f_enum,
  sdid UUID,
  sitecond fieldseeker.mosquitoinspection_mosquitoinspection_sitecond_db7350bc_81e5_401e_858f_cd3e5e5d8a34_enum,
  positivecontainercount SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  visualmonitoring fieldseeker.mosquitoinspection_notinuit_f_enum,
  vmcomments VARCHAR(250),
  adminaction fieldseeker.mosquitoinspection_mosquitoinspection_adminaction_b74ae1bb_c98b_40f6_8cfa_40e4fd16c270_enum,
  ptaid UUID,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.mosquitoinspection.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.mosquitoinspection.numdips IS '# Dips';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.activity IS 'Activity';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.breeding IS 'Breeding';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.totlarvae IS 'Total Larvae';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.totpupae IS 'Total Pupae';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.eggs IS 'Eggs';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.posdips IS 'Positive Dips';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.adultact IS 'Adult Activity';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.lstages IS 'Larval Stages';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.domstage IS 'Dominant Stage';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.actiontaken IS 'Action';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avetemp IS 'Average Temperature';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.windspeed IS 'Wind Speed';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.raingauge IS 'Rain Gauge';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.winddir IS 'Wind Direction';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avglarvae IS 'Average Larvae';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avgpupae IS 'Average Pupae';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.recordstatus IS 'RecordStatus';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.personalcontact IS 'Personal Contact';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.tirecount IS 'Tire Count';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.cbcount IS 'Catch Basin Count';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.containercount IS 'Container Count';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.fieldspecies IS 'Field Species';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.fieldtech IS 'Field Tech';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.larvaepresent IS 'Larvae Present';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.pupaepresent IS 'Pupae Present';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.sdid IS 'Storm Drain ID';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.sitecond IS 'Conditions';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.positivecontainercount IS 'Positive Container Count';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.jurisdiction IS 'Jurisdiction';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.visualmonitoring IS 'Visual Monitoring';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.vmcomments IS 'VM Comments';

-- Field adminaction has default value: 'no'

-- See insert/insert_mosquitoinspection_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
