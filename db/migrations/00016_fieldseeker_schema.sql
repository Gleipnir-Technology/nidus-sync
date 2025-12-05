-- +goose Up
-- Table definition for fieldseeker.containerrelate
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.containerrelate (
  objectid BIGSERIAL NOT NULL,
  
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  inspsampleid UUID,
  mosquitoinspid UUID,
  treatmentid UUID,
  containertype VARCHAR(250),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.containerrelate.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.containerrelate.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.containerrelate.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.containerrelate.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.containerrelate.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.containerrelate.inspsampleid IS 'Original attribute from ArcGIS API is INSPSAMPLEID';
COMMENT ON COLUMN fieldseeker.containerrelate.mosquitoinspid IS 'Original attribute from ArcGIS API is MOSQUITOINSPID';
COMMENT ON COLUMN fieldseeker.containerrelate.treatmentid IS 'Original attribute from ArcGIS API is TREATMENTID';
COMMENT ON COLUMN fieldseeker.containerrelate.containertype IS 'Original attribute from ArcGIS API is CONTAINERTYPE';
COMMENT ON COLUMN fieldseeker.containerrelate.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.containerrelate.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.containerrelate.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.containerrelate.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.fieldscoutinglog
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.fieldscoutinglog (
  objectid BIGSERIAL NOT NULL,
  
  status SMALLINT,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.fieldscoutinglog.status IS 'Original attribute from ArcGIS API is STATUS';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.fieldscoutinglog.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.habitatrelate
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.habitatrelate (
  objectid BIGSERIAL NOT NULL,
  
  foreign_id UUID,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  habitattype VARCHAR(250),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.habitatrelate.foreign_id IS 'Original attribute from ArcGIS API is FOREIGN_ID';
COMMENT ON COLUMN fieldseeker.habitatrelate.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.habitatrelate.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.habitatrelate.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.habitatrelate.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.habitatrelate.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.habitatrelate.habitattype IS 'Original attribute from ArcGIS API is HABITATTYPE';
COMMENT ON COLUMN fieldseeker.habitatrelate.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.habitatrelate.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.habitatrelate.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.habitatrelate.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.inspectionsample
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.inspectionsample (
  objectid BIGSERIAL NOT NULL,
  
  insp_id UUID,
  sampleid VARCHAR(25),
  processed SMALLINT,
  idbytech VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.inspectionsample.insp_id IS 'Original attribute from ArcGIS API is INSP_ID';
COMMENT ON COLUMN fieldseeker.inspectionsample.sampleid IS 'Original attribute from ArcGIS API is SAMPLEID';
COMMENT ON COLUMN fieldseeker.inspectionsample.processed IS 'Original attribute from ArcGIS API is PROCESSED';
COMMENT ON COLUMN fieldseeker.inspectionsample.idbytech IS 'Original attribute from ArcGIS API is IDBYTECH';
COMMENT ON COLUMN fieldseeker.inspectionsample.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.inspectionsample.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.inspectionsample.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.inspectionsample.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.inspectionsample.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.inspectionsample.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.inspectionsample.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.inspectionsample.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.inspectionsample.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.inspectionsampledetail
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.inspectionsampledetail (
  objectid BIGSERIAL NOT NULL,
  
  inspsample_id UUID,
  fieldspecies VARCHAR(25),
  flarvcount SMALLINT,
  fpupcount SMALLINT,
  feggcount SMALLINT,
  flstages VARCHAR(25),
  fdomstage VARCHAR(25),
  fadultact VARCHAR(25),
  labspecies VARCHAR(50),
  llarvcount SMALLINT,
  lpupcount SMALLINT,
  leggcount SMALLINT,
  ldomstage VARCHAR(25),
  comments VARCHAR(250),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  processed SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.inspectionsampledetail.inspsample_id IS 'Original attribute from ArcGIS API is INSPSAMPLE_ID';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.fieldspecies IS 'Original attribute from ArcGIS API is FIELDSPECIES';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.flarvcount IS 'Original attribute from ArcGIS API is FLARVCOUNT';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.fpupcount IS 'Original attribute from ArcGIS API is FPUPCOUNT';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.feggcount IS 'Original attribute from ArcGIS API is FEGGCOUNT';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.flstages IS 'Original attribute from ArcGIS API is FLSTAGES';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.fdomstage IS 'Original attribute from ArcGIS API is FDOMSTAGE';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.fadultact IS 'Original attribute from ArcGIS API is FADULTACT';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.labspecies IS 'Original attribute from ArcGIS API is LABSPECIES';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.llarvcount IS 'Original attribute from ArcGIS API is LLARVCOUNT';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.lpupcount IS 'Original attribute from ArcGIS API is LPUPCOUNT';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.leggcount IS 'Original attribute from ArcGIS API is LEGGCOUNT';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.ldomstage IS 'Original attribute from ArcGIS API is LDOMSTAGE';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.processed IS 'Original attribute from ArcGIS API is PROCESSED';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.inspectionsampledetail.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.linelocation
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.linelocation (
  objectid BIGSERIAL NOT NULL,
  
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat VARCHAR(25),
  priority VARCHAR(25),
  usetype VARCHAR(25),
  active SMALLINT,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology VARCHAR(10),
  externalid VARCHAR(50),
  acres DOUBLE PRECISION,
  nextactiondatescheduled TIMESTAMP,
  larvinspectinterval SMALLINT,
  length_ft DOUBLE PRECISION,
  width_ft DOUBLE PRECISION,
  zone2 VARCHAR(25),
  locationnumber INTEGER,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  lastinspectdate TIMESTAMP,
  lastinspectbreeding VARCHAR(25),
  lastinspectavglarvae DOUBLE PRECISION,
  lastinspectavgpupae DOUBLE PRECISION,
  lastinspectlstages VARCHAR(25),
  lastinspectactiontaken VARCHAR(50),
  lastinspectfieldspecies VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  hectares DOUBLE PRECISION,
  lastinspectactivity VARCHAR(25),
  lasttreatactivity VARCHAR(25),
  length_meters DOUBLE PRECISION,
  width_meters DOUBLE PRECISION,
  lastinspectconditions VARCHAR(250),
  waterorigin VARCHAR(50),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  shape__length DOUBLE PRECISION,
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(LINESTRING, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.linelocation.name IS 'Original attribute from ArcGIS API is NAME';
COMMENT ON COLUMN fieldseeker.linelocation.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.linelocation.habitat IS 'Original attribute from ArcGIS API is HABITAT';
COMMENT ON COLUMN fieldseeker.linelocation.priority IS 'Original attribute from ArcGIS API is PRIORITY';
COMMENT ON COLUMN fieldseeker.linelocation.usetype IS 'Original attribute from ArcGIS API is USETYPE';
COMMENT ON COLUMN fieldseeker.linelocation.active IS 'Original attribute from ArcGIS API is ACTIVE';
COMMENT ON COLUMN fieldseeker.linelocation.description IS 'Original attribute from ArcGIS API is DESCRIPTION';
COMMENT ON COLUMN fieldseeker.linelocation.accessdesc IS 'Original attribute from ArcGIS API is ACCESSDESC';
COMMENT ON COLUMN fieldseeker.linelocation.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.linelocation.symbology IS 'Original attribute from ArcGIS API is SYMBOLOGY';
COMMENT ON COLUMN fieldseeker.linelocation.externalid IS 'Original attribute from ArcGIS API is EXTERNALID';
COMMENT ON COLUMN fieldseeker.linelocation.acres IS 'Original attribute from ArcGIS API is ACRES';
COMMENT ON COLUMN fieldseeker.linelocation.nextactiondatescheduled IS 'Original attribute from ArcGIS API is NEXTACTIONDATESCHEDULED';
COMMENT ON COLUMN fieldseeker.linelocation.larvinspectinterval IS 'Original attribute from ArcGIS API is LARVINSPECTINTERVAL';
COMMENT ON COLUMN fieldseeker.linelocation.length_ft IS 'Original attribute from ArcGIS API is LENGTH_FT';
COMMENT ON COLUMN fieldseeker.linelocation.width_ft IS 'Original attribute from ArcGIS API is WIDTH_FT';
COMMENT ON COLUMN fieldseeker.linelocation.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.linelocation.locationnumber IS 'Original attribute from ArcGIS API is LOCATIONNUMBER';
COMMENT ON COLUMN fieldseeker.linelocation.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.linelocation.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.linelocation.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.linelocation.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.linelocation.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectdate IS 'Original attribute from ArcGIS API is LASTINSPECTDATE';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectbreeding IS 'Original attribute from ArcGIS API is LASTINSPECTBREEDING';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectavglarvae IS 'Original attribute from ArcGIS API is LASTINSPECTAVGLARVAE';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectavgpupae IS 'Original attribute from ArcGIS API is LASTINSPECTAVGPUPAE';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectlstages IS 'Original attribute from ArcGIS API is LASTINSPECTLSTAGES';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectactiontaken IS 'Original attribute from ArcGIS API is LASTINSPECTACTIONTAKEN';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectfieldspecies IS 'Original attribute from ArcGIS API is LASTINSPECTFIELDSPECIES';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatdate IS 'Original attribute from ArcGIS API is LASTTREATDATE';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatproduct IS 'Original attribute from ArcGIS API is LASTTREATPRODUCT';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatqty IS 'Original attribute from ArcGIS API is LASTTREATQTY';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatqtyunit IS 'Original attribute from ArcGIS API is LASTTREATQTYUNIT';
COMMENT ON COLUMN fieldseeker.linelocation.hectares IS 'Original attribute from ArcGIS API is HECTARES';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectactivity IS 'Original attribute from ArcGIS API is LASTINSPECTACTIVITY';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatactivity IS 'Original attribute from ArcGIS API is LASTTREATACTIVITY';
COMMENT ON COLUMN fieldseeker.linelocation.length_meters IS 'Original attribute from ArcGIS API is LENGTH_METERS';
COMMENT ON COLUMN fieldseeker.linelocation.width_meters IS 'Original attribute from ArcGIS API is WIDTH_METERS';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectconditions IS 'Original attribute from ArcGIS API is LASTINSPECTCONDITIONS';
COMMENT ON COLUMN fieldseeker.linelocation.waterorigin IS 'Original attribute from ArcGIS API is WATERORIGIN';
COMMENT ON COLUMN fieldseeker.linelocation.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.linelocation.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.linelocation.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.linelocation.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.linelocation.jurisdiction IS 'Original attribute from ArcGIS API is JURISDICTION';
COMMENT ON COLUMN fieldseeker.linelocation.shape__length IS 'Original attribute from ArcGIS API is Shape__Length';

-- Table definition for fieldseeker.locationtracking
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.locationtracking (
  objectid BIGSERIAL NOT NULL,
  
  accuracy DOUBLE PRECISION,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  globalid UUID,
  fieldtech VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.locationtracking.accuracy IS 'Original attribute from ArcGIS API is Accuracy';
COMMENT ON COLUMN fieldseeker.locationtracking.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.locationtracking.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.locationtracking.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.locationtracking.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.locationtracking.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.locationtracking.fieldtech IS 'Original attribute from ArcGIS API is FIELDTECH';
COMMENT ON COLUMN fieldseeker.locationtracking.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.locationtracking.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.locationtracking.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.locationtracking.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.mosquitoinspection
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.mosquitoinspection (
  objectid BIGSERIAL NOT NULL,
  
  numdips SMALLINT,
  activity VARCHAR(25),
  breeding VARCHAR(25),
  totlarvae SMALLINT,
  totpupae SMALLINT,
  eggs SMALLINT,
  posdips SMALLINT,
  adultact VARCHAR(25),
  lstages VARCHAR(25),
  domstage VARCHAR(25),
  actiontaken VARCHAR(50),
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  raingauge DOUBLE PRECISION,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  winddir VARCHAR(3),
  avglarvae DOUBLE PRECISION,
  avgpupae DOUBLE PRECISION,
  reviewed SMALLINT,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  personalcontact SMALLINT,
  tirecount SMALLINT,
  cbcount SMALLINT,
  containercount SMALLINT,
  fieldspecies VARCHAR(25),
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
  larvaepresent SMALLINT,
  pupaepresent SMALLINT,
  sdid UUID,
  sitecond VARCHAR(250),
  positivecontainercount SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  visualmonitoring SMALLINT,
  vmcomments VARCHAR(250),
  adminaction VARCHAR(256),
  ptaid UUID,
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.mosquitoinspection.numdips IS 'Original attribute from ArcGIS API is NUMDIPS';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.activity IS 'Original attribute from ArcGIS API is ACTIVITY';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.breeding IS 'Original attribute from ArcGIS API is BREEDING';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.totlarvae IS 'Original attribute from ArcGIS API is TOTLARVAE';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.totpupae IS 'Original attribute from ArcGIS API is TOTPUPAE';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.eggs IS 'Original attribute from ArcGIS API is EGGS';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.posdips IS 'Original attribute from ArcGIS API is POSDIPS';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.adultact IS 'Original attribute from ArcGIS API is ADULTACT';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.lstages IS 'Original attribute from ArcGIS API is LSTAGES';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.domstage IS 'Original attribute from ArcGIS API is DOMSTAGE';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.actiontaken IS 'Original attribute from ArcGIS API is ACTIONTAKEN';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avetemp IS 'Original attribute from ArcGIS API is AVETEMP';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.windspeed IS 'Original attribute from ArcGIS API is WINDSPEED';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.raingauge IS 'Original attribute from ArcGIS API is RAINGAUGE';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.startdatetime IS 'Original attribute from ArcGIS API is STARTDATETIME';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.enddatetime IS 'Original attribute from ArcGIS API is ENDDATETIME';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.winddir IS 'Original attribute from ArcGIS API is WINDDIR';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avglarvae IS 'Original attribute from ArcGIS API is AVGLARVAE';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.avgpupae IS 'Original attribute from ArcGIS API is AVGPUPAE';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.reviewed IS 'Original attribute from ArcGIS API is REVIEWED';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.reviewedby IS 'Original attribute from ArcGIS API is REVIEWEDBY';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.revieweddate IS 'Original attribute from ArcGIS API is REVIEWEDDATE';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.locationname IS 'Original attribute from ArcGIS API is LOCATIONNAME';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.recordstatus IS 'Original attribute from ArcGIS API is RECORDSTATUS';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.personalcontact IS 'Original attribute from ArcGIS API is PERSONALCONTACT';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.tirecount IS 'Original attribute from ArcGIS API is TIRECOUNT';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.cbcount IS 'Original attribute from ArcGIS API is CBCOUNT';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.containercount IS 'Original attribute from ArcGIS API is CONTAINERCOUNT';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.fieldspecies IS 'Original attribute from ArcGIS API is FIELDSPECIES';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.linelocid IS 'Original attribute from ArcGIS API is LINELOCID';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.pointlocid IS 'Original attribute from ArcGIS API is POINTLOCID';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.polygonlocid IS 'Original attribute from ArcGIS API is POLYGONLOCID';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.srid IS 'Original attribute from ArcGIS API is SRID';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.fieldtech IS 'Original attribute from ArcGIS API is FIELDTECH';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.larvaepresent IS 'Original attribute from ArcGIS API is LARVAEPRESENT';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.pupaepresent IS 'Original attribute from ArcGIS API is PUPAEPRESENT';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.sdid IS 'Original attribute from ArcGIS API is SDID';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.sitecond IS 'Original attribute from ArcGIS API is SITECOND';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.positivecontainercount IS 'Original attribute from ArcGIS API is POSITIVECONTAINERCOUNT';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.jurisdiction IS 'Original attribute from ArcGIS API is JURISDICTION';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.visualmonitoring IS 'Original attribute from ArcGIS API is VISUALMONITORING';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.vmcomments IS 'Original attribute from ArcGIS API is VMCOMMENTS';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.adminaction IS 'Original attribute from ArcGIS API is adminAction';
COMMENT ON COLUMN fieldseeker.mosquitoinspection.ptaid IS 'Original attribute from ArcGIS API is PTAID';

-- Table definition for fieldseeker.pointlocation
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.pointlocation (
  objectid BIGSERIAL NOT NULL,
  
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat VARCHAR(25),
  priority VARCHAR(25),
  usetype VARCHAR(25),
  active SMALLINT,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology VARCHAR(10),
  externalid VARCHAR(50),
  nextactiondatescheduled TIMESTAMP,
  larvinspectinterval SMALLINT,
  zone2 VARCHAR(25),
  locationnumber INTEGER,
  globalid UUID,
  stype VARCHAR(3),
  lastinspectdate TIMESTAMP,
  lastinspectbreeding VARCHAR(25),
  lastinspectavglarvae DOUBLE PRECISION,
  lastinspectavgpupae DOUBLE PRECISION,
  lastinspectlstages VARCHAR(25),
  lastinspectactiontaken VARCHAR(50),
  lastinspectfieldspecies VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  lastinspectactivity VARCHAR(25),
  lasttreatactivity VARCHAR(25),
  lastinspectconditions VARCHAR(250),
  waterorigin VARCHAR(50),
  x DOUBLE PRECISION,
  y DOUBLE PRECISION,
  assignedtech VARCHAR(256),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  deactivate_reason VARCHAR(256),
  scalarpriority INTEGER,
  sourcestatus VARCHAR(255),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.pointlocation.name IS 'Original attribute from ArcGIS API is NAME';
COMMENT ON COLUMN fieldseeker.pointlocation.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.pointlocation.habitat IS 'Original attribute from ArcGIS API is HABITAT';
COMMENT ON COLUMN fieldseeker.pointlocation.priority IS 'Original attribute from ArcGIS API is PRIORITY';
COMMENT ON COLUMN fieldseeker.pointlocation.usetype IS 'Original attribute from ArcGIS API is USETYPE';
COMMENT ON COLUMN fieldseeker.pointlocation.active IS 'Original attribute from ArcGIS API is ACTIVE';
COMMENT ON COLUMN fieldseeker.pointlocation.description IS 'Original attribute from ArcGIS API is DESCRIPTION';
COMMENT ON COLUMN fieldseeker.pointlocation.accessdesc IS 'Original attribute from ArcGIS API is ACCESSDESC';
COMMENT ON COLUMN fieldseeker.pointlocation.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.pointlocation.symbology IS 'Original attribute from ArcGIS API is SYMBOLOGY';
COMMENT ON COLUMN fieldseeker.pointlocation.externalid IS 'Original attribute from ArcGIS API is EXTERNALID';
COMMENT ON COLUMN fieldseeker.pointlocation.nextactiondatescheduled IS 'Original attribute from ArcGIS API is NEXTACTIONDATESCHEDULED';
COMMENT ON COLUMN fieldseeker.pointlocation.larvinspectinterval IS 'Original attribute from ArcGIS API is LARVINSPECTINTERVAL';
COMMENT ON COLUMN fieldseeker.pointlocation.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.pointlocation.locationnumber IS 'Original attribute from ArcGIS API is LOCATIONNUMBER';
COMMENT ON COLUMN fieldseeker.pointlocation.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.pointlocation.stype IS 'Original attribute from ArcGIS API is STYPE';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectdate IS 'Original attribute from ArcGIS API is LASTINSPECTDATE';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectbreeding IS 'Original attribute from ArcGIS API is LASTINSPECTBREEDING';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectavglarvae IS 'Original attribute from ArcGIS API is LASTINSPECTAVGLARVAE';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectavgpupae IS 'Original attribute from ArcGIS API is LASTINSPECTAVGPUPAE';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectlstages IS 'Original attribute from ArcGIS API is LASTINSPECTLSTAGES';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectactiontaken IS 'Original attribute from ArcGIS API is LASTINSPECTACTIONTAKEN';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectfieldspecies IS 'Original attribute from ArcGIS API is LASTINSPECTFIELDSPECIES';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatdate IS 'Original attribute from ArcGIS API is LASTTREATDATE';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatproduct IS 'Original attribute from ArcGIS API is LASTTREATPRODUCT';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatqty IS 'Original attribute from ArcGIS API is LASTTREATQTY';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatqtyunit IS 'Original attribute from ArcGIS API is LASTTREATQTYUNIT';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectactivity IS 'Original attribute from ArcGIS API is LASTINSPECTACTIVITY';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatactivity IS 'Original attribute from ArcGIS API is LASTTREATACTIVITY';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectconditions IS 'Original attribute from ArcGIS API is LASTINSPECTCONDITIONS';
COMMENT ON COLUMN fieldseeker.pointlocation.waterorigin IS 'Original attribute from ArcGIS API is WATERORIGIN';
COMMENT ON COLUMN fieldseeker.pointlocation.x IS 'Original attribute from ArcGIS API is X';
COMMENT ON COLUMN fieldseeker.pointlocation.y IS 'Original attribute from ArcGIS API is Y';
COMMENT ON COLUMN fieldseeker.pointlocation.assignedtech IS 'Original attribute from ArcGIS API is ASSIGNEDTECH';
COMMENT ON COLUMN fieldseeker.pointlocation.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.pointlocation.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.pointlocation.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.pointlocation.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.pointlocation.jurisdiction IS 'Original attribute from ArcGIS API is JURISDICTION';
COMMENT ON COLUMN fieldseeker.pointlocation.deactivate_reason IS 'Original attribute from ArcGIS API is deactivate_reason';
COMMENT ON COLUMN fieldseeker.pointlocation.scalarpriority IS 'Original attribute from ArcGIS API is scalarPriority';
COMMENT ON COLUMN fieldseeker.pointlocation.sourcestatus IS 'Original attribute from ArcGIS API is sourceStatus';

-- Table definition for fieldseeker.polygonlocation
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.polygonlocation (
  objectid BIGSERIAL NOT NULL,
  
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat VARCHAR(25),
  priority VARCHAR(25),
  usetype VARCHAR(25),
  active SMALLINT,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology VARCHAR(10),
  externalid VARCHAR(50),
  acres DOUBLE PRECISION,
  nextactiondatescheduled TIMESTAMP,
  larvinspectinterval SMALLINT,
  zone2 VARCHAR(25),
  locationnumber INTEGER,
  globalid UUID,
  lastinspectdate TIMESTAMP,
  lastinspectbreeding VARCHAR(25),
  lastinspectavglarvae DOUBLE PRECISION,
  lastinspectavgpupae DOUBLE PRECISION,
  lastinspectlstages VARCHAR(25),
  lastinspectactiontaken VARCHAR(50),
  lastinspectfieldspecies VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  hectares DOUBLE PRECISION,
  lastinspectactivity VARCHAR(25),
  lasttreatactivity VARCHAR(25),
  lastinspectconditions VARCHAR(250),
  waterorigin VARCHAR(50),
  filter VARCHAR(255),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POLYGON, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.polygonlocation.name IS 'Original attribute from ArcGIS API is NAME';
COMMENT ON COLUMN fieldseeker.polygonlocation.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.polygonlocation.habitat IS 'Original attribute from ArcGIS API is HABITAT';
COMMENT ON COLUMN fieldseeker.polygonlocation.priority IS 'Original attribute from ArcGIS API is PRIORITY';
COMMENT ON COLUMN fieldseeker.polygonlocation.usetype IS 'Original attribute from ArcGIS API is USETYPE';
COMMENT ON COLUMN fieldseeker.polygonlocation.active IS 'Original attribute from ArcGIS API is ACTIVE';
COMMENT ON COLUMN fieldseeker.polygonlocation.description IS 'Original attribute from ArcGIS API is DESCRIPTION';
COMMENT ON COLUMN fieldseeker.polygonlocation.accessdesc IS 'Original attribute from ArcGIS API is ACCESSDESC';
COMMENT ON COLUMN fieldseeker.polygonlocation.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.polygonlocation.symbology IS 'Original attribute from ArcGIS API is SYMBOLOGY';
COMMENT ON COLUMN fieldseeker.polygonlocation.externalid IS 'Original attribute from ArcGIS API is EXTERNALID';
COMMENT ON COLUMN fieldseeker.polygonlocation.acres IS 'Original attribute from ArcGIS API is ACRES';
COMMENT ON COLUMN fieldseeker.polygonlocation.nextactiondatescheduled IS 'Original attribute from ArcGIS API is NEXTACTIONDATESCHEDULED';
COMMENT ON COLUMN fieldseeker.polygonlocation.larvinspectinterval IS 'Original attribute from ArcGIS API is LARVINSPECTINTERVAL';
COMMENT ON COLUMN fieldseeker.polygonlocation.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.polygonlocation.locationnumber IS 'Original attribute from ArcGIS API is LOCATIONNUMBER';
COMMENT ON COLUMN fieldseeker.polygonlocation.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectdate IS 'Original attribute from ArcGIS API is LASTINSPECTDATE';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectbreeding IS 'Original attribute from ArcGIS API is LASTINSPECTBREEDING';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectavglarvae IS 'Original attribute from ArcGIS API is LASTINSPECTAVGLARVAE';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectavgpupae IS 'Original attribute from ArcGIS API is LASTINSPECTAVGPUPAE';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectlstages IS 'Original attribute from ArcGIS API is LASTINSPECTLSTAGES';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectactiontaken IS 'Original attribute from ArcGIS API is LASTINSPECTACTIONTAKEN';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectfieldspecies IS 'Original attribute from ArcGIS API is LASTINSPECTFIELDSPECIES';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatdate IS 'Original attribute from ArcGIS API is LASTTREATDATE';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatproduct IS 'Original attribute from ArcGIS API is LASTTREATPRODUCT';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatqty IS 'Original attribute from ArcGIS API is LASTTREATQTY';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatqtyunit IS 'Original attribute from ArcGIS API is LASTTREATQTYUNIT';
COMMENT ON COLUMN fieldseeker.polygonlocation.hectares IS 'Original attribute from ArcGIS API is HECTARES';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectactivity IS 'Original attribute from ArcGIS API is LASTINSPECTACTIVITY';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatactivity IS 'Original attribute from ArcGIS API is LASTTREATACTIVITY';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectconditions IS 'Original attribute from ArcGIS API is LASTINSPECTCONDITIONS';
COMMENT ON COLUMN fieldseeker.polygonlocation.waterorigin IS 'Original attribute from ArcGIS API is WATERORIGIN';
COMMENT ON COLUMN fieldseeker.polygonlocation.filter IS 'Original attribute from ArcGIS API is Filter';
COMMENT ON COLUMN fieldseeker.polygonlocation.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.polygonlocation.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.polygonlocation.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.polygonlocation.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.polygonlocation.jurisdiction IS 'Original attribute from ArcGIS API is JURISDICTION';
COMMENT ON COLUMN fieldseeker.polygonlocation.shape__area IS 'Original attribute from ArcGIS API is Shape__Area';
COMMENT ON COLUMN fieldseeker.polygonlocation.shape__length IS 'Original attribute from ArcGIS API is Shape__Length';

-- Table definition for fieldseeker.pool
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.pool (
  objectid BIGSERIAL NOT NULL,
  
  trapdata_id UUID,
  datesent TIMESTAMP,
  survtech VARCHAR(25),
  datetested TIMESTAMP,
  testtech VARCHAR(25),
  comments VARCHAR(250),
  sampleid VARCHAR(50),
  processed SMALLINT,
  lab_id UUID,
  testmethod VARCHAR(100),
  diseasetested VARCHAR(100),
  diseasepos VARCHAR(100),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  lab VARCHAR(25),
  poolyear SMALLINT,
  gatewaysync SMALLINT,
  vectorsurvcollectionid VARCHAR(50),
  vectorsurvpoolid VARCHAR(50),
  vectorsurvtrapdataid VARCHAR(50),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.pool.trapdata_id IS 'Original attribute from ArcGIS API is TRAPDATA_ID';
COMMENT ON COLUMN fieldseeker.pool.datesent IS 'Original attribute from ArcGIS API is DATESENT';
COMMENT ON COLUMN fieldseeker.pool.survtech IS 'Original attribute from ArcGIS API is SURVTECH';
COMMENT ON COLUMN fieldseeker.pool.datetested IS 'Original attribute from ArcGIS API is DATETESTED';
COMMENT ON COLUMN fieldseeker.pool.testtech IS 'Original attribute from ArcGIS API is TESTTECH';
COMMENT ON COLUMN fieldseeker.pool.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.pool.sampleid IS 'Original attribute from ArcGIS API is SAMPLEID';
COMMENT ON COLUMN fieldseeker.pool.processed IS 'Original attribute from ArcGIS API is PROCESSED';
COMMENT ON COLUMN fieldseeker.pool.lab_id IS 'Original attribute from ArcGIS API is LAB_ID';
COMMENT ON COLUMN fieldseeker.pool.testmethod IS 'Original attribute from ArcGIS API is TESTMETHOD';
COMMENT ON COLUMN fieldseeker.pool.diseasetested IS 'Original attribute from ArcGIS API is DISEASETESTED';
COMMENT ON COLUMN fieldseeker.pool.diseasepos IS 'Original attribute from ArcGIS API is DISEASEPOS';
COMMENT ON COLUMN fieldseeker.pool.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.pool.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.pool.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.pool.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.pool.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.pool.lab IS 'Original attribute from ArcGIS API is LAB';
COMMENT ON COLUMN fieldseeker.pool.poolyear IS 'Original attribute from ArcGIS API is POOLYEAR';
COMMENT ON COLUMN fieldseeker.pool.gatewaysync IS 'Original attribute from ArcGIS API is GATEWAYSYNC';
COMMENT ON COLUMN fieldseeker.pool.vectorsurvcollectionid IS 'Original attribute from ArcGIS API is VECTORSURVCOLLECTIONID';
COMMENT ON COLUMN fieldseeker.pool.vectorsurvpoolid IS 'Original attribute from ArcGIS API is VECTORSURVPOOLID';
COMMENT ON COLUMN fieldseeker.pool.vectorsurvtrapdataid IS 'Original attribute from ArcGIS API is VECTORSURVTRAPDATAID';
COMMENT ON COLUMN fieldseeker.pool.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.pool.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.pool.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.pool.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.pooldetail
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
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
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.pooldetail.trapdata_id IS 'Original attribute from ArcGIS API is TRAPDATA_ID';
COMMENT ON COLUMN fieldseeker.pooldetail.pool_id IS 'Original attribute from ArcGIS API is POOL_ID';
COMMENT ON COLUMN fieldseeker.pooldetail.species IS 'Original attribute from ArcGIS API is SPECIES';
COMMENT ON COLUMN fieldseeker.pooldetail.females IS 'Original attribute from ArcGIS API is FEMALES';
COMMENT ON COLUMN fieldseeker.pooldetail.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.pooldetail.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.pooldetail.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.pooldetail.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.pooldetail.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.pooldetail.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.pooldetail.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.pooldetail.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.pooldetail.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.proposedtreatmentarea
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.proposedtreatmentarea (
  objectid BIGSERIAL NOT NULL,
  
  method VARCHAR(25),
  comments VARCHAR(250),
  zone VARCHAR(25),
  reviewed SMALLINT,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  zone2 VARCHAR(25),
  completeddate TIMESTAMP,
  completedby VARCHAR(25),
  completed SMALLINT,
  issprayroute SMALLINT,
  name VARCHAR(25),
  acres DOUBLE PRECISION,
  globalid UUID,
  exported SMALLINT,
  targetproduct VARCHAR(25),
  targetapprate DOUBLE PRECISION,
  hectares DOUBLE PRECISION,
  lasttreatactivity VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  priority VARCHAR(25),
  duedate TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  targetspecies VARCHAR(250),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POLYGON, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.method IS 'Original attribute from ArcGIS API is METHOD';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.reviewed IS 'Original attribute from ArcGIS API is REVIEWED';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.reviewedby IS 'Original attribute from ArcGIS API is REVIEWEDBY';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.revieweddate IS 'Original attribute from ArcGIS API is REVIEWEDDATE';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completeddate IS 'Original attribute from ArcGIS API is COMPLETEDDATE';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completedby IS 'Original attribute from ArcGIS API is COMPLETEDBY';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completed IS 'Original attribute from ArcGIS API is COMPLETED';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.issprayroute IS 'Original attribute from ArcGIS API is ISSPRAYROUTE';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.name IS 'Original attribute from ArcGIS API is NAME';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.acres IS 'Original attribute from ArcGIS API is ACRES';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.exported IS 'Original attribute from ArcGIS API is EXPORTED';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetproduct IS 'Original attribute from ArcGIS API is TARGETPRODUCT';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetapprate IS 'Original attribute from ArcGIS API is TARGETAPPRATE';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.hectares IS 'Original attribute from ArcGIS API is HECTARES';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatactivity IS 'Original attribute from ArcGIS API is LASTTREATACTIVITY';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatdate IS 'Original attribute from ArcGIS API is LASTTREATDATE';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatproduct IS 'Original attribute from ArcGIS API is LASTTREATPRODUCT';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatqty IS 'Original attribute from ArcGIS API is LASTTREATQTY';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatqtyunit IS 'Original attribute from ArcGIS API is LASTTREATQTYUNIT';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.priority IS 'Original attribute from ArcGIS API is PRIORITY';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.duedate IS 'Original attribute from ArcGIS API is DUEDATE';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetspecies IS 'Original attribute from ArcGIS API is TARGETSPECIES';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.shape__area IS 'Original attribute from ArcGIS API is Shape__Area';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.shape__length IS 'Original attribute from ArcGIS API is Shape__Length';

-- Table definition for fieldseeker.qamosquitoinspection
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.qamosquitoinspection (
  objectid BIGSERIAL NOT NULL,
  
  posdips SMALLINT,
  actiontaken VARCHAR(250),
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  raingauge DOUBLE PRECISION,
  globalid UUID,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  winddir VARCHAR(3),
  reviewed SMALLINT,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  lr SMALLINT,
  negdips SMALLINT,
  totalacres DOUBLE PRECISION,
  acresbreeding DOUBLE PRECISION,
  fish SMALLINT,
  sitetype VARCHAR(250),
  breedingpotential VARCHAR(25),
  movingwater SMALLINT,
  nowaterever SMALLINT,
  mosquitohabitat VARCHAR(250),
  habvalue1 SMALLINT,
  habvalue1percent SMALLINT,
  habvalue2 SMALLINT,
  habvalue2percent SMALLINT,
  potential SMALLINT,
  larvaepresent SMALLINT,
  larvaeinsidetreatedarea SMALLINT,
  larvaeoutsidetreatedarea SMALLINT,
  larvaereason VARCHAR(250),
  aquaticorganisms VARCHAR(500),
  vegetation VARCHAR(500),
  sourcereduction VARCHAR(250),
  waterpresent SMALLINT,
  watermovement1 VARCHAR(25),
  watermovement1percent SMALLINT,
  watermovement2 VARCHAR(25),
  watermovement2percent SMALLINT,
  soilconditions VARCHAR(250),
  waterduration VARCHAR(25),
  watersource VARCHAR(250),
  waterconditions VARCHAR(250),
  adultactivity SMALLINT,
  linelocid UUID,
  pointlocid UUID,
  polygonlocid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  fieldtech VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.qamosquitoinspection.posdips IS 'Original attribute from ArcGIS API is POSDIPS';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.actiontaken IS 'Original attribute from ArcGIS API is ACTIONTAKEN';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.avetemp IS 'Original attribute from ArcGIS API is AVETEMP';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.windspeed IS 'Original attribute from ArcGIS API is WINDSPEED';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.raingauge IS 'Original attribute from ArcGIS API is RAINGAUGE';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.startdatetime IS 'Original attribute from ArcGIS API is STARTDATETIME';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.enddatetime IS 'Original attribute from ArcGIS API is ENDDATETIME';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.winddir IS 'Original attribute from ArcGIS API is WINDDIR';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.reviewed IS 'Original attribute from ArcGIS API is REVIEWED';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.reviewedby IS 'Original attribute from ArcGIS API is REVIEWEDBY';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.revieweddate IS 'Original attribute from ArcGIS API is REVIEWEDDATE';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.locationname IS 'Original attribute from ArcGIS API is LOCATIONNAME';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.recordstatus IS 'Original attribute from ArcGIS API is RECORDSTATUS';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.lr IS 'Original attribute from ArcGIS API is LR';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.negdips IS 'Original attribute from ArcGIS API is NEGDIPS';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.totalacres IS 'Original attribute from ArcGIS API is TOTALACRES';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.acresbreeding IS 'Original attribute from ArcGIS API is ACRESBREEDING';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.fish IS 'Original attribute from ArcGIS API is FISH';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.sitetype IS 'Original attribute from ArcGIS API is SITETYPE';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.breedingpotential IS 'Original attribute from ArcGIS API is BREEDINGPOTENTIAL';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.movingwater IS 'Original attribute from ArcGIS API is MOVINGWATER';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.nowaterever IS 'Original attribute from ArcGIS API is NOWATEREVER';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.mosquitohabitat IS 'Original attribute from ArcGIS API is MOSQUITOHABITAT';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue1 IS 'Original attribute from ArcGIS API is HABVALUE1';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue1percent IS 'Original attribute from ArcGIS API is HABVALUE1PERCENT';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue2 IS 'Original attribute from ArcGIS API is HABVALUE2';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue2percent IS 'Original attribute from ArcGIS API is HABVALUE2PERCENT';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.potential IS 'Original attribute from ArcGIS API is POTENTIAL';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaepresent IS 'Original attribute from ArcGIS API is LARVAEPRESENT';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaeinsidetreatedarea IS 'Original attribute from ArcGIS API is LARVAEINSIDETREATEDAREA';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaeoutsidetreatedarea IS 'Original attribute from ArcGIS API is LARVAEOUTSIDETREATEDAREA';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaereason IS 'Original attribute from ArcGIS API is LARVAEREASON';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.aquaticorganisms IS 'Original attribute from ArcGIS API is AQUATICORGANISMS';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.vegetation IS 'Original attribute from ArcGIS API is VEGETATION';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.sourcereduction IS 'Original attribute from ArcGIS API is SOURCEREDUCTION';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterpresent IS 'Original attribute from ArcGIS API is WATERPRESENT';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement1 IS 'Original attribute from ArcGIS API is WATERMOVEMENT1';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement1percent IS 'Original attribute from ArcGIS API is WATERMOVEMENT1PERCENT';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement2 IS 'Original attribute from ArcGIS API is WATERMOVEMENT2';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement2percent IS 'Original attribute from ArcGIS API is WATERMOVEMENT2PERCENT';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.soilconditions IS 'Original attribute from ArcGIS API is SOILCONDITIONS';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterduration IS 'Original attribute from ArcGIS API is WATERDURATION';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watersource IS 'Original attribute from ArcGIS API is WATERSOURCE';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterconditions IS 'Original attribute from ArcGIS API is WATERCONDITIONS';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.adultactivity IS 'Original attribute from ArcGIS API is ADULTACTIVITY';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.linelocid IS 'Original attribute from ArcGIS API is LINELOCID';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.pointlocid IS 'Original attribute from ArcGIS API is POINTLOCID';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.polygonlocid IS 'Original attribute from ArcGIS API is POLYGONLOCID';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.fieldtech IS 'Original attribute from ArcGIS API is FIELDTECH';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.rodentlocation
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.rodentlocation (
  objectid BIGSERIAL NOT NULL,
  
  locationname VARCHAR(25),
  zone VARCHAR(25),
  zone2 VARCHAR(25),
  habitat VARCHAR(25),
  priority VARCHAR(25),
  usetype VARCHAR(25),
  active SMALLINT,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology VARCHAR(10),
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
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.rodentlocation.locationname IS 'Original attribute from ArcGIS API is LOCATIONNAME';
COMMENT ON COLUMN fieldseeker.rodentlocation.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.rodentlocation.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.rodentlocation.habitat IS 'Original attribute from ArcGIS API is HABITAT';
COMMENT ON COLUMN fieldseeker.rodentlocation.priority IS 'Original attribute from ArcGIS API is PRIORITY';
COMMENT ON COLUMN fieldseeker.rodentlocation.usetype IS 'Original attribute from ArcGIS API is USETYPE';
COMMENT ON COLUMN fieldseeker.rodentlocation.active IS 'Original attribute from ArcGIS API is ACTIVE';
COMMENT ON COLUMN fieldseeker.rodentlocation.description IS 'Original attribute from ArcGIS API is DESCRIPTION';
COMMENT ON COLUMN fieldseeker.rodentlocation.accessdesc IS 'Original attribute from ArcGIS API is ACCESSDESC';
COMMENT ON COLUMN fieldseeker.rodentlocation.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.rodentlocation.symbology IS 'Original attribute from ArcGIS API is SYMBOLOGY';
COMMENT ON COLUMN fieldseeker.rodentlocation.externalid IS 'Original attribute from ArcGIS API is EXTERNALID';
COMMENT ON COLUMN fieldseeker.rodentlocation.nextactiondatescheduled IS 'Original attribute from ArcGIS API is NEXTACTIONDATESCHEDULED';
COMMENT ON COLUMN fieldseeker.rodentlocation.locationnumber IS 'Original attribute from ArcGIS API is LOCATIONNUMBER';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectdate IS 'Original attribute from ArcGIS API is LASTINSPECTDATE';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectspecies IS 'Original attribute from ArcGIS API is LASTINSPECTSPECIES';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectaction IS 'Original attribute from ArcGIS API is LASTINSPECTACTION';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectconditions IS 'Original attribute from ArcGIS API is LASTINSPECTCONDITIONS';
COMMENT ON COLUMN fieldseeker.rodentlocation.lastinspectrodentevidence IS 'Original attribute from ArcGIS API is LASTINSPECTRODENTEVIDENCE';
COMMENT ON COLUMN fieldseeker.rodentlocation.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.rodentlocation.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.rodentlocation.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.rodentlocation.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.rodentlocation.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.rodentlocation.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.rodentlocation.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.rodentlocation.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.rodentlocation.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.rodentlocation.jurisdiction IS 'Original attribute from ArcGIS API is JURISDICTION';

-- Table definition for fieldseeker.samplecollection
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.samplecollection (
  objectid BIGSERIAL NOT NULL,
  
  loc_id UUID,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  sitecond VARCHAR(25),
  sampleid VARCHAR(25),
  survtech VARCHAR(25),
  datesent TIMESTAMP,
  datetested TIMESTAMP,
  testtech VARCHAR(25),
  comments VARCHAR(250),
  processed SMALLINT,
  sampletype VARCHAR(25),
  samplecond VARCHAR(25),
  species VARCHAR(25),
  sex VARCHAR(1),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  winddir VARCHAR(3),
  raingauge DOUBLE PRECISION,
  activity VARCHAR(25),
  testmethod VARCHAR(100),
  diseasetested VARCHAR(100),
  diseasepos VARCHAR(100),
  reviewed SMALLINT,
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
  lab VARCHAR(25),
  fieldtech VARCHAR(25),
  flockid UUID,
  samplecount SMALLINT,
  chickenid UUID,
  gatewaysync SMALLINT,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.samplecollection.loc_id IS 'Original attribute from ArcGIS API is LOC_ID';
COMMENT ON COLUMN fieldseeker.samplecollection.startdatetime IS 'Original attribute from ArcGIS API is STARTDATETIME';
COMMENT ON COLUMN fieldseeker.samplecollection.enddatetime IS 'Original attribute from ArcGIS API is ENDDATETIME';
COMMENT ON COLUMN fieldseeker.samplecollection.sitecond IS 'Original attribute from ArcGIS API is SITECOND';
COMMENT ON COLUMN fieldseeker.samplecollection.sampleid IS 'Original attribute from ArcGIS API is SAMPLEID';
COMMENT ON COLUMN fieldseeker.samplecollection.survtech IS 'Original attribute from ArcGIS API is SURVTECH';
COMMENT ON COLUMN fieldseeker.samplecollection.datesent IS 'Original attribute from ArcGIS API is DATESENT';
COMMENT ON COLUMN fieldseeker.samplecollection.datetested IS 'Original attribute from ArcGIS API is DATETESTED';
COMMENT ON COLUMN fieldseeker.samplecollection.testtech IS 'Original attribute from ArcGIS API is TESTTECH';
COMMENT ON COLUMN fieldseeker.samplecollection.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.samplecollection.processed IS 'Original attribute from ArcGIS API is PROCESSED';
COMMENT ON COLUMN fieldseeker.samplecollection.sampletype IS 'Original attribute from ArcGIS API is SAMPLETYPE';
COMMENT ON COLUMN fieldseeker.samplecollection.samplecond IS 'Original attribute from ArcGIS API is SAMPLECOND';
COMMENT ON COLUMN fieldseeker.samplecollection.species IS 'Original attribute from ArcGIS API is SPECIES';
COMMENT ON COLUMN fieldseeker.samplecollection.sex IS 'Original attribute from ArcGIS API is SEX';
COMMENT ON COLUMN fieldseeker.samplecollection.avetemp IS 'Original attribute from ArcGIS API is AVETEMP';
COMMENT ON COLUMN fieldseeker.samplecollection.windspeed IS 'Original attribute from ArcGIS API is WINDSPEED';
COMMENT ON COLUMN fieldseeker.samplecollection.winddir IS 'Original attribute from ArcGIS API is WINDDIR';
COMMENT ON COLUMN fieldseeker.samplecollection.raingauge IS 'Original attribute from ArcGIS API is RAINGAUGE';
COMMENT ON COLUMN fieldseeker.samplecollection.activity IS 'Original attribute from ArcGIS API is ACTIVITY';
COMMENT ON COLUMN fieldseeker.samplecollection.testmethod IS 'Original attribute from ArcGIS API is TESTMETHOD';
COMMENT ON COLUMN fieldseeker.samplecollection.diseasetested IS 'Original attribute from ArcGIS API is DISEASETESTED';
COMMENT ON COLUMN fieldseeker.samplecollection.diseasepos IS 'Original attribute from ArcGIS API is DISEASEPOS';
COMMENT ON COLUMN fieldseeker.samplecollection.reviewed IS 'Original attribute from ArcGIS API is REVIEWED';
COMMENT ON COLUMN fieldseeker.samplecollection.reviewedby IS 'Original attribute from ArcGIS API is REVIEWEDBY';
COMMENT ON COLUMN fieldseeker.samplecollection.revieweddate IS 'Original attribute from ArcGIS API is REVIEWEDDATE';
COMMENT ON COLUMN fieldseeker.samplecollection.locationname IS 'Original attribute from ArcGIS API is LOCATIONNAME';
COMMENT ON COLUMN fieldseeker.samplecollection.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.samplecollection.recordstatus IS 'Original attribute from ArcGIS API is RECORDSTATUS';
COMMENT ON COLUMN fieldseeker.samplecollection.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.samplecollection.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.samplecollection.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.samplecollection.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.samplecollection.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.samplecollection.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.samplecollection.lab IS 'Original attribute from ArcGIS API is LAB';
COMMENT ON COLUMN fieldseeker.samplecollection.fieldtech IS 'Original attribute from ArcGIS API is FIELDTECH';
COMMENT ON COLUMN fieldseeker.samplecollection.flockid IS 'Original attribute from ArcGIS API is FLOCKID';
COMMENT ON COLUMN fieldseeker.samplecollection.samplecount IS 'Original attribute from ArcGIS API is SAMPLECOUNT';
COMMENT ON COLUMN fieldseeker.samplecollection.chickenid IS 'Original attribute from ArcGIS API is CHICKENID';
COMMENT ON COLUMN fieldseeker.samplecollection.gatewaysync IS 'Original attribute from ArcGIS API is GATEWAYSYNC';
COMMENT ON COLUMN fieldseeker.samplecollection.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.samplecollection.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.samplecollection.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.samplecollection.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.samplelocation
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.samplelocation (
  objectid BIGSERIAL NOT NULL,
  
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat VARCHAR(25),
  priority VARCHAR(25),
  usetype VARCHAR(25),
  active SMALLINT,
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
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.samplelocation.name IS 'Original attribute from ArcGIS API is NAME';
COMMENT ON COLUMN fieldseeker.samplelocation.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.samplelocation.habitat IS 'Original attribute from ArcGIS API is HABITAT';
COMMENT ON COLUMN fieldseeker.samplelocation.priority IS 'Original attribute from ArcGIS API is PRIORITY';
COMMENT ON COLUMN fieldseeker.samplelocation.usetype IS 'Original attribute from ArcGIS API is USETYPE';
COMMENT ON COLUMN fieldseeker.samplelocation.active IS 'Original attribute from ArcGIS API is ACTIVE';
COMMENT ON COLUMN fieldseeker.samplelocation.description IS 'Original attribute from ArcGIS API is DESCRIPTION';
COMMENT ON COLUMN fieldseeker.samplelocation.accessdesc IS 'Original attribute from ArcGIS API is ACCESSDESC';
COMMENT ON COLUMN fieldseeker.samplelocation.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.samplelocation.externalid IS 'Original attribute from ArcGIS API is EXTERNALID';
COMMENT ON COLUMN fieldseeker.samplelocation.nextactiondatescheduled IS 'Original attribute from ArcGIS API is NEXTACTIONDATESCHEDULED';
COMMENT ON COLUMN fieldseeker.samplelocation.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.samplelocation.locationnumber IS 'Original attribute from ArcGIS API is LOCATIONNUMBER';
COMMENT ON COLUMN fieldseeker.samplelocation.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.samplelocation.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.samplelocation.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.samplelocation.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.samplelocation.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.samplelocation.gatewaysync IS 'Original attribute from ArcGIS API is GATEWAYSYNC';
COMMENT ON COLUMN fieldseeker.samplelocation.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.samplelocation.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.samplelocation.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.samplelocation.editor IS 'Original attribute from ArcGIS API is Editor';

-- Table definition for fieldseeker.servicerequest
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.servicerequest (
  objectid BIGSERIAL NOT NULL,
  
  recdatetime TIMESTAMP,
  source VARCHAR(25),
  entrytech VARCHAR(25),
  priority VARCHAR(25),
  supervisor VARCHAR(25),
  assignedtech VARCHAR(25),
  status VARCHAR(25),
  clranon SMALLINT,
  clrfname VARCHAR(25),
  clrphone1 VARCHAR(25),
  clrphone2 VARCHAR(25),
  clremail VARCHAR(50),
  clrcompany VARCHAR(25),
  clraddr1 VARCHAR(50),
  clraddr2 VARCHAR(50),
  clrcity VARCHAR(25),
  clrstate VARCHAR(25),
  clrzip VARCHAR(25),
  clrother VARCHAR(25),
  clrcontpref VARCHAR(25),
  reqcompany VARCHAR(25),
  reqaddr1 VARCHAR(50),
  reqaddr2 VARCHAR(50),
  reqcity VARCHAR(25),
  reqstate VARCHAR(25),
  reqzip VARCHAR(25),
  reqcrossst VARCHAR(25),
  reqsubdiv VARCHAR(25),
  reqmapgrid VARCHAR(25),
  reqpermission SMALLINT,
  reqtarget VARCHAR(25),
  reqdescr VARCHAR(1000),
  reqnotesfortech VARCHAR(250),
  reqnotesforcust VARCHAR(250),
  reqfldnotes VARCHAR(250),
  reqprogramactions VARCHAR(250),
  datetimeclosed TIMESTAMP,
  techclosed VARCHAR(25),
  sr_number INTEGER,
  reviewed SMALLINT,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  accepted SMALLINT,
  accepteddate TIMESTAMP,
  rejectedby VARCHAR(25),
  rejecteddate TIMESTAMP,
  rejectedreason VARCHAR(25),
  duedate TIMESTAMP,
  acceptedby VARCHAR(25),
  comments VARCHAR(2500),
  estcompletedate TIMESTAMP,
  nextaction VARCHAR(25),
  recordstatus SMALLINT,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  firstresponsedate TIMESTAMP,
  responsedaycount SMALLINT,
  allowed VARCHAR(25),
  xvalue VARCHAR(25),
  yvalue VARCHAR(25),
  validx VARCHAR(25),
  validy VARCHAR(25),
  externalid VARCHAR(25),
  externalerror VARCHAR(500),
  pointlocid UUID,
  notified SMALLINT,
  notifieddate TIMESTAMP,
  scheduled SMALLINT,
  scheduleddate TIMESTAMP,
  dog INTEGER,
  schedule_period VARCHAR(20),
  schedule_notes VARCHAR(256),
  spanish INTEGER,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  issuesreported VARCHAR(250),
  jurisdiction VARCHAR(25),
  notificationtimestamp VARCHAR(250),
  zone VARCHAR(50),
  zone2 VARCHAR(50),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.servicerequest.recdatetime IS 'Original attribute from ArcGIS API is RECDATETIME';
COMMENT ON COLUMN fieldseeker.servicerequest.source IS 'Original attribute from ArcGIS API is SOURCE';
COMMENT ON COLUMN fieldseeker.servicerequest.entrytech IS 'Original attribute from ArcGIS API is ENTRYTECH';
COMMENT ON COLUMN fieldseeker.servicerequest.priority IS 'Original attribute from ArcGIS API is PRIORITY';
COMMENT ON COLUMN fieldseeker.servicerequest.supervisor IS 'Original attribute from ArcGIS API is SUPERVISOR';
COMMENT ON COLUMN fieldseeker.servicerequest.assignedtech IS 'Original attribute from ArcGIS API is ASSIGNEDTECH';
COMMENT ON COLUMN fieldseeker.servicerequest.status IS 'Original attribute from ArcGIS API is STATUS';
COMMENT ON COLUMN fieldseeker.servicerequest.clranon IS 'Original attribute from ArcGIS API is CLRANON';
COMMENT ON COLUMN fieldseeker.servicerequest.clrfname IS 'Original attribute from ArcGIS API is CLRFNAME';
COMMENT ON COLUMN fieldseeker.servicerequest.clrphone1 IS 'Original attribute from ArcGIS API is CLRPHONE1';
COMMENT ON COLUMN fieldseeker.servicerequest.clrphone2 IS 'Original attribute from ArcGIS API is CLRPHONE2';
COMMENT ON COLUMN fieldseeker.servicerequest.clremail IS 'Original attribute from ArcGIS API is CLREMAIL';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcompany IS 'Original attribute from ArcGIS API is CLRCOMPANY';
COMMENT ON COLUMN fieldseeker.servicerequest.clraddr1 IS 'Original attribute from ArcGIS API is CLRADDR1';
COMMENT ON COLUMN fieldseeker.servicerequest.clraddr2 IS 'Original attribute from ArcGIS API is CLRADDR2';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcity IS 'Original attribute from ArcGIS API is CLRCITY';
COMMENT ON COLUMN fieldseeker.servicerequest.clrstate IS 'Original attribute from ArcGIS API is CLRSTATE';
COMMENT ON COLUMN fieldseeker.servicerequest.clrzip IS 'Original attribute from ArcGIS API is CLRZIP';
COMMENT ON COLUMN fieldseeker.servicerequest.clrother IS 'Original attribute from ArcGIS API is CLROTHER';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcontpref IS 'Original attribute from ArcGIS API is CLRCONTPREF';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcompany IS 'Original attribute from ArcGIS API is REQCOMPANY';
COMMENT ON COLUMN fieldseeker.servicerequest.reqaddr1 IS 'Original attribute from ArcGIS API is REQADDR1';
COMMENT ON COLUMN fieldseeker.servicerequest.reqaddr2 IS 'Original attribute from ArcGIS API is REQADDR2';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcity IS 'Original attribute from ArcGIS API is REQCITY';
COMMENT ON COLUMN fieldseeker.servicerequest.reqstate IS 'Original attribute from ArcGIS API is REQSTATE';
COMMENT ON COLUMN fieldseeker.servicerequest.reqzip IS 'Original attribute from ArcGIS API is REQZIP';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcrossst IS 'Original attribute from ArcGIS API is REQCROSSST';
COMMENT ON COLUMN fieldseeker.servicerequest.reqsubdiv IS 'Original attribute from ArcGIS API is REQSUBDIV';
COMMENT ON COLUMN fieldseeker.servicerequest.reqmapgrid IS 'Original attribute from ArcGIS API is REQMAPGRID';
COMMENT ON COLUMN fieldseeker.servicerequest.reqpermission IS 'Original attribute from ArcGIS API is REQPERMISSION';
COMMENT ON COLUMN fieldseeker.servicerequest.reqtarget IS 'Original attribute from ArcGIS API is REQTARGET';
COMMENT ON COLUMN fieldseeker.servicerequest.reqdescr IS 'Original attribute from ArcGIS API is REQDESCR';
COMMENT ON COLUMN fieldseeker.servicerequest.reqnotesfortech IS 'Original attribute from ArcGIS API is REQNOTESFORTECH';
COMMENT ON COLUMN fieldseeker.servicerequest.reqnotesforcust IS 'Original attribute from ArcGIS API is REQNOTESFORCUST';
COMMENT ON COLUMN fieldseeker.servicerequest.reqfldnotes IS 'Original attribute from ArcGIS API is REQFLDNOTES';
COMMENT ON COLUMN fieldseeker.servicerequest.reqprogramactions IS 'Original attribute from ArcGIS API is REQPROGRAMACTIONS';
COMMENT ON COLUMN fieldseeker.servicerequest.datetimeclosed IS 'Original attribute from ArcGIS API is DATETIMECLOSED';
COMMENT ON COLUMN fieldseeker.servicerequest.techclosed IS 'Original attribute from ArcGIS API is TECHCLOSED';
COMMENT ON COLUMN fieldseeker.servicerequest.sr_number IS 'Original attribute from ArcGIS API is SR_NUMBER';
COMMENT ON COLUMN fieldseeker.servicerequest.reviewed IS 'Original attribute from ArcGIS API is REVIEWED';
COMMENT ON COLUMN fieldseeker.servicerequest.reviewedby IS 'Original attribute from ArcGIS API is REVIEWEDBY';
COMMENT ON COLUMN fieldseeker.servicerequest.revieweddate IS 'Original attribute from ArcGIS API is REVIEWEDDATE';
COMMENT ON COLUMN fieldseeker.servicerequest.accepted IS 'Original attribute from ArcGIS API is ACCEPTED';
COMMENT ON COLUMN fieldseeker.servicerequest.accepteddate IS 'Original attribute from ArcGIS API is ACCEPTEDDATE';
COMMENT ON COLUMN fieldseeker.servicerequest.rejectedby IS 'Original attribute from ArcGIS API is REJECTEDBY';
COMMENT ON COLUMN fieldseeker.servicerequest.rejecteddate IS 'Original attribute from ArcGIS API is REJECTEDDATE';
COMMENT ON COLUMN fieldseeker.servicerequest.rejectedreason IS 'Original attribute from ArcGIS API is REJECTEDREASON';
COMMENT ON COLUMN fieldseeker.servicerequest.duedate IS 'Original attribute from ArcGIS API is DUEDATE';
COMMENT ON COLUMN fieldseeker.servicerequest.acceptedby IS 'Original attribute from ArcGIS API is ACCEPTEDBY';
COMMENT ON COLUMN fieldseeker.servicerequest.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.servicerequest.estcompletedate IS 'Original attribute from ArcGIS API is ESTCOMPLETEDATE';
COMMENT ON COLUMN fieldseeker.servicerequest.nextaction IS 'Original attribute from ArcGIS API is NEXTACTION';
COMMENT ON COLUMN fieldseeker.servicerequest.recordstatus IS 'Original attribute from ArcGIS API is RECORDSTATUS';
COMMENT ON COLUMN fieldseeker.servicerequest.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.servicerequest.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.servicerequest.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.servicerequest.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.servicerequest.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.servicerequest.firstresponsedate IS 'Original attribute from ArcGIS API is FIRSTRESPONSEDATE';
COMMENT ON COLUMN fieldseeker.servicerequest.responsedaycount IS 'Original attribute from ArcGIS API is RESPONSEDAYCOUNT';
COMMENT ON COLUMN fieldseeker.servicerequest.allowed IS 'Original attribute from ArcGIS API is ALLOWED';
COMMENT ON COLUMN fieldseeker.servicerequest.xvalue IS 'Original attribute from ArcGIS API is XVALUE';
COMMENT ON COLUMN fieldseeker.servicerequest.yvalue IS 'Original attribute from ArcGIS API is YVALUE';
COMMENT ON COLUMN fieldseeker.servicerequest.validx IS 'Original attribute from ArcGIS API is VALIDX';
COMMENT ON COLUMN fieldseeker.servicerequest.validy IS 'Original attribute from ArcGIS API is VALIDY';
COMMENT ON COLUMN fieldseeker.servicerequest.externalid IS 'Original attribute from ArcGIS API is EXTERNALID';
COMMENT ON COLUMN fieldseeker.servicerequest.externalerror IS 'Original attribute from ArcGIS API is EXTERNALERROR';
COMMENT ON COLUMN fieldseeker.servicerequest.pointlocid IS 'Original attribute from ArcGIS API is POINTLOCID';
COMMENT ON COLUMN fieldseeker.servicerequest.notified IS 'Original attribute from ArcGIS API is NOTIFIED';
COMMENT ON COLUMN fieldseeker.servicerequest.notifieddate IS 'Original attribute from ArcGIS API is NOTIFIEDDATE';
COMMENT ON COLUMN fieldseeker.servicerequest.scheduled IS 'Original attribute from ArcGIS API is SCHEDULED';
COMMENT ON COLUMN fieldseeker.servicerequest.scheduleddate IS 'Original attribute from ArcGIS API is SCHEDULEDDATE';
COMMENT ON COLUMN fieldseeker.servicerequest.dog IS 'Original attribute from ArcGIS API is DOG';
COMMENT ON COLUMN fieldseeker.servicerequest.schedule_period IS 'Original attribute from ArcGIS API is schedule_period';
COMMENT ON COLUMN fieldseeker.servicerequest.schedule_notes IS 'Original attribute from ArcGIS API is schedule_notes';
COMMENT ON COLUMN fieldseeker.servicerequest.spanish IS 'Original attribute from ArcGIS API is Spanish';
COMMENT ON COLUMN fieldseeker.servicerequest.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.servicerequest.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.servicerequest.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.servicerequest.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.servicerequest.issuesreported IS 'Original attribute from ArcGIS API is ISSUESREPORTED';
COMMENT ON COLUMN fieldseeker.servicerequest.jurisdiction IS 'Original attribute from ArcGIS API is JURISDICTION';
COMMENT ON COLUMN fieldseeker.servicerequest.notificationtimestamp IS 'Original attribute from ArcGIS API is NOTIFICATIONTIMESTAMP';
COMMENT ON COLUMN fieldseeker.servicerequest.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.servicerequest.zone2 IS 'Original attribute from ArcGIS API is ZONE2';

-- Table definition for fieldseeker.speciesabundance
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

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
  processed SMALLINT,
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
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.speciesabundance.trapdata_id IS 'Original attribute from ArcGIS API is TRAPDATA_ID';
COMMENT ON COLUMN fieldseeker.speciesabundance.species IS 'Original attribute from ArcGIS API is SPECIES';
COMMENT ON COLUMN fieldseeker.speciesabundance.males IS 'Original attribute from ArcGIS API is MALES';
COMMENT ON COLUMN fieldseeker.speciesabundance.unknown IS 'Original attribute from ArcGIS API is UNKNOWN';
COMMENT ON COLUMN fieldseeker.speciesabundance.bloodedfem IS 'Original attribute from ArcGIS API is BLOODEDFEM';
COMMENT ON COLUMN fieldseeker.speciesabundance.gravidfem IS 'Original attribute from ArcGIS API is GRAVIDFEM';
COMMENT ON COLUMN fieldseeker.speciesabundance.larvae IS 'Original attribute from ArcGIS API is LARVAE';
COMMENT ON COLUMN fieldseeker.speciesabundance.poolstogen IS 'Original attribute from ArcGIS API is POOLSTOGEN';
COMMENT ON COLUMN fieldseeker.speciesabundance.processed IS 'Original attribute from ArcGIS API is PROCESSED';
COMMENT ON COLUMN fieldseeker.speciesabundance.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.speciesabundance.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.speciesabundance.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.speciesabundance.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.speciesabundance.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.speciesabundance.pupae IS 'Original attribute from ArcGIS API is PUPAE';
COMMENT ON COLUMN fieldseeker.speciesabundance.eggs IS 'Original attribute from ArcGIS API is EGGS';
COMMENT ON COLUMN fieldseeker.speciesabundance.females IS 'Original attribute from ArcGIS API is FEMALES';
COMMENT ON COLUMN fieldseeker.speciesabundance.total IS 'Original attribute from ArcGIS API is TOTAL';
COMMENT ON COLUMN fieldseeker.speciesabundance.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.speciesabundance.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.speciesabundance.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.speciesabundance.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.speciesabundance.yearweek IS 'Original attribute from ArcGIS API is yearWeek';
COMMENT ON COLUMN fieldseeker.speciesabundance.globalzscore IS 'Original attribute from ArcGIS API is globalZScore';
COMMENT ON COLUMN fieldseeker.speciesabundance.r7score IS 'Original attribute from ArcGIS API is r7Score';
COMMENT ON COLUMN fieldseeker.speciesabundance.r8score IS 'Original attribute from ArcGIS API is r8Score';
COMMENT ON COLUMN fieldseeker.speciesabundance.h3r7 IS 'Original attribute from ArcGIS API is h3r7';
COMMENT ON COLUMN fieldseeker.speciesabundance.h3r8 IS 'Original attribute from ArcGIS API is h3r8';

-- Table definition for fieldseeker.stormdrain
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.stormdrain (
  objectid BIGSERIAL NOT NULL,
  
  nexttreatmentdate TIMESTAMP,
  lasttreatdate TIMESTAMP,
  lastaction VARCHAR(25),
  symbology VARCHAR(25),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  laststatus VARCHAR(25),
  zone VARCHAR(25),
  zone2 VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  type VARCHAR(25),
  jurisdiction VARCHAR(25),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.stormdrain.nexttreatmentdate IS 'Original attribute from ArcGIS API is NextTreatmentDate';
COMMENT ON COLUMN fieldseeker.stormdrain.lasttreatdate IS 'Original attribute from ArcGIS API is LastTreatDate';
COMMENT ON COLUMN fieldseeker.stormdrain.lastaction IS 'Original attribute from ArcGIS API is LastAction';
COMMENT ON COLUMN fieldseeker.stormdrain.symbology IS 'Original attribute from ArcGIS API is Symbology';
COMMENT ON COLUMN fieldseeker.stormdrain.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.stormdrain.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.stormdrain.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.stormdrain.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.stormdrain.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.stormdrain.laststatus IS 'Original attribute from ArcGIS API is LastStatus';
COMMENT ON COLUMN fieldseeker.stormdrain.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.stormdrain.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.stormdrain.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.stormdrain.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.stormdrain.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.stormdrain.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.stormdrain.type IS 'Original attribute from ArcGIS API is TYPE';
COMMENT ON COLUMN fieldseeker.stormdrain.jurisdiction IS 'Original attribute from ArcGIS API is JURISDICTION';

-- Table definition for fieldseeker.timecard
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.timecard (
  objectid BIGSERIAL NOT NULL,
  
  activity VARCHAR(25),
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  comments VARCHAR(250),
  externalid VARCHAR(25),
  equiptype VARCHAR(25),
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
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.timecard.activity IS 'Original attribute from ArcGIS API is ACTIVITY';
COMMENT ON COLUMN fieldseeker.timecard.startdatetime IS 'Original attribute from ArcGIS API is STARTDATETIME';
COMMENT ON COLUMN fieldseeker.timecard.enddatetime IS 'Original attribute from ArcGIS API is ENDDATETIME';
COMMENT ON COLUMN fieldseeker.timecard.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.timecard.externalid IS 'Original attribute from ArcGIS API is EXTERNALID';
COMMENT ON COLUMN fieldseeker.timecard.equiptype IS 'Original attribute from ArcGIS API is EQUIPTYPE';
COMMENT ON COLUMN fieldseeker.timecard.locationname IS 'Original attribute from ArcGIS API is LOCATIONNAME';
COMMENT ON COLUMN fieldseeker.timecard.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.timecard.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.timecard.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.timecard.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.timecard.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.timecard.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.timecard.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.timecard.linelocid IS 'Original attribute from ArcGIS API is LINELOCID';
COMMENT ON COLUMN fieldseeker.timecard.pointlocid IS 'Original attribute from ArcGIS API is POINTLOCID';
COMMENT ON COLUMN fieldseeker.timecard.polygonlocid IS 'Original attribute from ArcGIS API is POLYGONLOCID';
COMMENT ON COLUMN fieldseeker.timecard.lclocid IS 'Original attribute from ArcGIS API is LCLOCID';
COMMENT ON COLUMN fieldseeker.timecard.samplelocid IS 'Original attribute from ArcGIS API is SAMPLELOCID';
COMMENT ON COLUMN fieldseeker.timecard.srid IS 'Original attribute from ArcGIS API is SRID';
COMMENT ON COLUMN fieldseeker.timecard.traplocid IS 'Original attribute from ArcGIS API is TRAPLOCID';
COMMENT ON COLUMN fieldseeker.timecard.fieldtech IS 'Original attribute from ArcGIS API is FIELDTECH';
COMMENT ON COLUMN fieldseeker.timecard.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.timecard.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.timecard.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.timecard.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.timecard.rodentlocid IS 'Original attribute from ArcGIS API is RODENTLOCID';

-- Table definition for fieldseeker.trapdata
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.trapdata (
  objectid BIGSERIAL NOT NULL,
  
  traptype VARCHAR(25),
  trapactivitytype VARCHAR(1),
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  comments VARCHAR(250),
  idbytech VARCHAR(25),
  sortbytech VARCHAR(25),
  processed SMALLINT,
  sitecond VARCHAR(25),
  locationname VARCHAR(25),
  recordstatus SMALLINT,
  reviewed SMALLINT,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  trapcondition VARCHAR(25),
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
  winddir VARCHAR(3),
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
  lure VARCHAR(256),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.trapdata.traptype IS 'Original attribute from ArcGIS API is TRAPTYPE';
COMMENT ON COLUMN fieldseeker.trapdata.trapactivitytype IS 'Original attribute from ArcGIS API is TRAPACTIVITYTYPE';
COMMENT ON COLUMN fieldseeker.trapdata.startdatetime IS 'Original attribute from ArcGIS API is STARTDATETIME';
COMMENT ON COLUMN fieldseeker.trapdata.enddatetime IS 'Original attribute from ArcGIS API is ENDDATETIME';
COMMENT ON COLUMN fieldseeker.trapdata.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.trapdata.idbytech IS 'Original attribute from ArcGIS API is IDBYTECH';
COMMENT ON COLUMN fieldseeker.trapdata.sortbytech IS 'Original attribute from ArcGIS API is SORTBYTECH';
COMMENT ON COLUMN fieldseeker.trapdata.processed IS 'Original attribute from ArcGIS API is PROCESSED';
COMMENT ON COLUMN fieldseeker.trapdata.sitecond IS 'Original attribute from ArcGIS API is SITECOND';
COMMENT ON COLUMN fieldseeker.trapdata.locationname IS 'Original attribute from ArcGIS API is LOCATIONNAME';
COMMENT ON COLUMN fieldseeker.trapdata.recordstatus IS 'Original attribute from ArcGIS API is RECORDSTATUS';
COMMENT ON COLUMN fieldseeker.trapdata.reviewed IS 'Original attribute from ArcGIS API is REVIEWED';
COMMENT ON COLUMN fieldseeker.trapdata.reviewedby IS 'Original attribute from ArcGIS API is REVIEWEDBY';
COMMENT ON COLUMN fieldseeker.trapdata.revieweddate IS 'Original attribute from ArcGIS API is REVIEWEDDATE';
COMMENT ON COLUMN fieldseeker.trapdata.trapcondition IS 'Original attribute from ArcGIS API is TRAPCONDITION';
COMMENT ON COLUMN fieldseeker.trapdata.trapnights IS 'Original attribute from ArcGIS API is TRAPNIGHTS';
COMMENT ON COLUMN fieldseeker.trapdata.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.trapdata.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.trapdata.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.trapdata.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.trapdata.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.trapdata.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.trapdata.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.trapdata.srid IS 'Original attribute from ArcGIS API is SRID';
COMMENT ON COLUMN fieldseeker.trapdata.fieldtech IS 'Original attribute from ArcGIS API is FIELDTECH';
COMMENT ON COLUMN fieldseeker.trapdata.gatewaysync IS 'Original attribute from ArcGIS API is GATEWAYSYNC';
COMMENT ON COLUMN fieldseeker.trapdata.loc_id IS 'Original attribute from ArcGIS API is LOC_ID';
COMMENT ON COLUMN fieldseeker.trapdata.voltage IS 'Original attribute from ArcGIS API is VOLTAGE';
COMMENT ON COLUMN fieldseeker.trapdata.winddir IS 'Original attribute from ArcGIS API is WINDDIR';
COMMENT ON COLUMN fieldseeker.trapdata.windspeed IS 'Original attribute from ArcGIS API is WINDSPEED';
COMMENT ON COLUMN fieldseeker.trapdata.avetemp IS 'Original attribute from ArcGIS API is AVETEMP';
COMMENT ON COLUMN fieldseeker.trapdata.raingauge IS 'Original attribute from ArcGIS API is RAINGAUGE';
COMMENT ON COLUMN fieldseeker.trapdata.lr IS 'Original attribute from ArcGIS API is LR';
COMMENT ON COLUMN fieldseeker.trapdata.field IS 'Original attribute from ArcGIS API is Field';
COMMENT ON COLUMN fieldseeker.trapdata.vectorsurvtrapdataid IS 'Original attribute from ArcGIS API is VECTORSURVTRAPDATAID';
COMMENT ON COLUMN fieldseeker.trapdata.vectorsurvtraplocationid IS 'Original attribute from ArcGIS API is VECTORSURVTRAPLOCATIONID';
COMMENT ON COLUMN fieldseeker.trapdata.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.trapdata.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.trapdata.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.trapdata.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.trapdata.lure IS 'Original attribute from ArcGIS API is Lure';

-- Table definition for fieldseeker.traplocation
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.traplocation (
  objectid BIGSERIAL NOT NULL,
  
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat VARCHAR(25),
  priority VARCHAR(25),
  usetype VARCHAR(25),
  active SMALLINT,
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
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.traplocation.name IS 'Original attribute from ArcGIS API is NAME';
COMMENT ON COLUMN fieldseeker.traplocation.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.traplocation.habitat IS 'Original attribute from ArcGIS API is HABITAT';
COMMENT ON COLUMN fieldseeker.traplocation.priority IS 'Original attribute from ArcGIS API is PRIORITY';
COMMENT ON COLUMN fieldseeker.traplocation.usetype IS 'Original attribute from ArcGIS API is USETYPE';
COMMENT ON COLUMN fieldseeker.traplocation.active IS 'Original attribute from ArcGIS API is ACTIVE';
COMMENT ON COLUMN fieldseeker.traplocation.description IS 'Original attribute from ArcGIS API is DESCRIPTION';
COMMENT ON COLUMN fieldseeker.traplocation.accessdesc IS 'Original attribute from ArcGIS API is ACCESSDESC';
COMMENT ON COLUMN fieldseeker.traplocation.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.traplocation.externalid IS 'Original attribute from ArcGIS API is EXTERNALID';
COMMENT ON COLUMN fieldseeker.traplocation.nextactiondatescheduled IS 'Original attribute from ArcGIS API is NEXTACTIONDATESCHEDULED';
COMMENT ON COLUMN fieldseeker.traplocation.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.traplocation.locationnumber IS 'Original attribute from ArcGIS API is LOCATIONNUMBER';
COMMENT ON COLUMN fieldseeker.traplocation.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.traplocation.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.traplocation.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.traplocation.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.traplocation.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.traplocation.gatewaysync IS 'Original attribute from ArcGIS API is GATEWAYSYNC';
COMMENT ON COLUMN fieldseeker.traplocation.route IS 'Original attribute from ArcGIS API is route';
COMMENT ON COLUMN fieldseeker.traplocation.set_dow IS 'Original attribute from ArcGIS API is set_dow';
COMMENT ON COLUMN fieldseeker.traplocation.route_order IS 'Original attribute from ArcGIS API is route_order';
COMMENT ON COLUMN fieldseeker.traplocation.vectorsurvsiteid IS 'Original attribute from ArcGIS API is VECTORSURVSITEID';
COMMENT ON COLUMN fieldseeker.traplocation.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.traplocation.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.traplocation.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.traplocation.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.traplocation.h3r7 IS 'Original attribute from ArcGIS API is h3r7';
COMMENT ON COLUMN fieldseeker.traplocation.h3r8 IS 'Original attribute from ArcGIS API is h3r8';

-- Table definition for fieldseeker.treatment
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.treatment (
  objectid BIGSERIAL NOT NULL,
  
  activity VARCHAR(25),
  treatarea DOUBLE PRECISION,
  areaunit VARCHAR(10),
  product VARCHAR(25),
  qty DOUBLE PRECISION,
  qtyunit VARCHAR(10),
  method VARCHAR(25),
  equiptype VARCHAR(25),
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  winddir VARCHAR(3),
  raingauge DOUBLE PRECISION,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  insp_id UUID,
  reviewed SMALLINT,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  warningoverride SMALLINT,
  recordstatus SMALLINT,
  zone2 VARCHAR(25),
  treatacres DOUBLE PRECISION,
  tirecount SMALLINT,
  cbcount SMALLINT,
  containercount SMALLINT,
  globalid UUID,
  treatmentlength DOUBLE PRECISION,
  treatmenthours DOUBLE PRECISION,
  treatmentlengthunits VARCHAR(5),
  linelocid UUID,
  pointlocid UUID,
  polygonlocid UUID,
  srid UUID,
  sdid UUID,
  barrierrouteid UUID,
  ulvrouteid UUID,
  fieldtech VARCHAR(25),
  ptaid UUID,
  flowrate DOUBLE PRECISION,
  habitat VARCHAR(250),
  treathectares DOUBLE PRECISION,
  invloc VARCHAR(25),
  temp_sitecond VARCHAR(250),
  sitecond VARCHAR(250),
  totalcostprodcut DOUBLE PRECISION,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  targetspecies VARCHAR(250),
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POINT, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.treatment.activity IS 'Original attribute from ArcGIS API is ACTIVITY';
COMMENT ON COLUMN fieldseeker.treatment.treatarea IS 'Original attribute from ArcGIS API is TREATAREA';
COMMENT ON COLUMN fieldseeker.treatment.areaunit IS 'Original attribute from ArcGIS API is AREAUNIT';
COMMENT ON COLUMN fieldseeker.treatment.product IS 'Original attribute from ArcGIS API is PRODUCT';
COMMENT ON COLUMN fieldseeker.treatment.qty IS 'Original attribute from ArcGIS API is QTY';
COMMENT ON COLUMN fieldseeker.treatment.qtyunit IS 'Original attribute from ArcGIS API is QTYUNIT';
COMMENT ON COLUMN fieldseeker.treatment.method IS 'Original attribute from ArcGIS API is METHOD';
COMMENT ON COLUMN fieldseeker.treatment.equiptype IS 'Original attribute from ArcGIS API is EQUIPTYPE';
COMMENT ON COLUMN fieldseeker.treatment.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.treatment.avetemp IS 'Original attribute from ArcGIS API is AVETEMP';
COMMENT ON COLUMN fieldseeker.treatment.windspeed IS 'Original attribute from ArcGIS API is WINDSPEED';
COMMENT ON COLUMN fieldseeker.treatment.winddir IS 'Original attribute from ArcGIS API is WINDDIR';
COMMENT ON COLUMN fieldseeker.treatment.raingauge IS 'Original attribute from ArcGIS API is RAINGAUGE';
COMMENT ON COLUMN fieldseeker.treatment.startdatetime IS 'Original attribute from ArcGIS API is STARTDATETIME';
COMMENT ON COLUMN fieldseeker.treatment.enddatetime IS 'Original attribute from ArcGIS API is ENDDATETIME';
COMMENT ON COLUMN fieldseeker.treatment.insp_id IS 'Original attribute from ArcGIS API is INSP_ID';
COMMENT ON COLUMN fieldseeker.treatment.reviewed IS 'Original attribute from ArcGIS API is REVIEWED';
COMMENT ON COLUMN fieldseeker.treatment.reviewedby IS 'Original attribute from ArcGIS API is REVIEWEDBY';
COMMENT ON COLUMN fieldseeker.treatment.revieweddate IS 'Original attribute from ArcGIS API is REVIEWEDDATE';
COMMENT ON COLUMN fieldseeker.treatment.locationname IS 'Original attribute from ArcGIS API is LOCATIONNAME';
COMMENT ON COLUMN fieldseeker.treatment.zone IS 'Original attribute from ArcGIS API is ZONE';
COMMENT ON COLUMN fieldseeker.treatment.warningoverride IS 'Original attribute from ArcGIS API is WARNINGOVERRIDE';
COMMENT ON COLUMN fieldseeker.treatment.recordstatus IS 'Original attribute from ArcGIS API is RECORDSTATUS';
COMMENT ON COLUMN fieldseeker.treatment.zone2 IS 'Original attribute from ArcGIS API is ZONE2';
COMMENT ON COLUMN fieldseeker.treatment.treatacres IS 'Original attribute from ArcGIS API is TREATACRES';
COMMENT ON COLUMN fieldseeker.treatment.tirecount IS 'Original attribute from ArcGIS API is TIRECOUNT';
COMMENT ON COLUMN fieldseeker.treatment.cbcount IS 'Original attribute from ArcGIS API is CBCOUNT';
COMMENT ON COLUMN fieldseeker.treatment.containercount IS 'Original attribute from ArcGIS API is CONTAINERCOUNT';
COMMENT ON COLUMN fieldseeker.treatment.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.treatment.treatmentlength IS 'Original attribute from ArcGIS API is TREATMENTLENGTH';
COMMENT ON COLUMN fieldseeker.treatment.treatmenthours IS 'Original attribute from ArcGIS API is TREATMENTHOURS';
COMMENT ON COLUMN fieldseeker.treatment.treatmentlengthunits IS 'Original attribute from ArcGIS API is TREATMENTLENGTHUNITS';
COMMENT ON COLUMN fieldseeker.treatment.linelocid IS 'Original attribute from ArcGIS API is LINELOCID';
COMMENT ON COLUMN fieldseeker.treatment.pointlocid IS 'Original attribute from ArcGIS API is POINTLOCID';
COMMENT ON COLUMN fieldseeker.treatment.polygonlocid IS 'Original attribute from ArcGIS API is POLYGONLOCID';
COMMENT ON COLUMN fieldseeker.treatment.srid IS 'Original attribute from ArcGIS API is SRID';
COMMENT ON COLUMN fieldseeker.treatment.sdid IS 'Original attribute from ArcGIS API is SDID';
COMMENT ON COLUMN fieldseeker.treatment.barrierrouteid IS 'Original attribute from ArcGIS API is BARRIERROUTEID';
COMMENT ON COLUMN fieldseeker.treatment.ulvrouteid IS 'Original attribute from ArcGIS API is ULVROUTEID';
COMMENT ON COLUMN fieldseeker.treatment.fieldtech IS 'Original attribute from ArcGIS API is FIELDTECH';
COMMENT ON COLUMN fieldseeker.treatment.ptaid IS 'Original attribute from ArcGIS API is PTAID';
COMMENT ON COLUMN fieldseeker.treatment.flowrate IS 'Original attribute from ArcGIS API is FLOWRATE';
COMMENT ON COLUMN fieldseeker.treatment.habitat IS 'Original attribute from ArcGIS API is HABITAT';
COMMENT ON COLUMN fieldseeker.treatment.treathectares IS 'Original attribute from ArcGIS API is TREATHECTARES';
COMMENT ON COLUMN fieldseeker.treatment.invloc IS 'Original attribute from ArcGIS API is INVLOC';
COMMENT ON COLUMN fieldseeker.treatment.temp_sitecond IS 'Original attribute from ArcGIS API is temp_SITECOND';
COMMENT ON COLUMN fieldseeker.treatment.sitecond IS 'Original attribute from ArcGIS API is SITECOND';
COMMENT ON COLUMN fieldseeker.treatment.totalcostprodcut IS 'Original attribute from ArcGIS API is TotalCostProdcut';
COMMENT ON COLUMN fieldseeker.treatment.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.treatment.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.treatment.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.treatment.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.treatment.targetspecies IS 'Original attribute from ArcGIS API is TARGETSPECIES';

-- Table definition for fieldseeker.treatmentarea
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.treatmentarea (
  objectid BIGSERIAL NOT NULL,
  
  treat_id UUID,
  session_id UUID,
  treatdate TIMESTAMP,
  comments VARCHAR(250),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  notified SMALLINT,
  type VARCHAR(25),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POLYGON, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.treatmentarea.treat_id IS 'Original attribute from ArcGIS API is TREAT_ID';
COMMENT ON COLUMN fieldseeker.treatmentarea.session_id IS 'Original attribute from ArcGIS API is SESSION_ID';
COMMENT ON COLUMN fieldseeker.treatmentarea.treatdate IS 'Original attribute from ArcGIS API is TREATDATE';
COMMENT ON COLUMN fieldseeker.treatmentarea.comments IS 'Original attribute from ArcGIS API is COMMENTS';
COMMENT ON COLUMN fieldseeker.treatmentarea.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.treatmentarea.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.treatmentarea.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.treatmentarea.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.treatmentarea.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.treatmentarea.notified IS 'Original attribute from ArcGIS API is Notified';
COMMENT ON COLUMN fieldseeker.treatmentarea.type IS 'Original attribute from ArcGIS API is Type';
COMMENT ON COLUMN fieldseeker.treatmentarea.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.treatmentarea.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.treatmentarea.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.treatmentarea.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.treatmentarea.shape__area IS 'Original attribute from ArcGIS API is Shape__Area';
COMMENT ON COLUMN fieldseeker.treatmentarea.shape__length IS 'Original attribute from ArcGIS API is Shape__Length';

-- Table definition for fieldseeker.zones
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.zones (
  objectid BIGSERIAL NOT NULL,
  
  name VARCHAR(50),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  active INTEGER,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POLYGON, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.zones.name IS 'Original attribute from ArcGIS API is NAME';
COMMENT ON COLUMN fieldseeker.zones.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.zones.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.zones.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.zones.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.zones.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.zones.active IS 'Original attribute from ArcGIS API is ACTIVE';
COMMENT ON COLUMN fieldseeker.zones.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.zones.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.zones.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.zones.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.zones.shape__area IS 'Original attribute from ArcGIS API is Shape__Area';
COMMENT ON COLUMN fieldseeker.zones.shape__length IS 'Original attribute from ArcGIS API is Shape__Length';

-- Table definition for fieldseeker.zones2
-- Includes versioning for tracking changes

-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TABLE fieldseeker.zones2 (
  objectid BIGSERIAL NOT NULL,
  
  name VARCHAR(50),
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  geometry JSONB NOT NULL,
  geospatial GEOMETRY(POLYGON, 3857),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);


COMMENT ON COLUMN fieldseeker.zones2.name IS 'Original attribute from ArcGIS API is NAME';
COMMENT ON COLUMN fieldseeker.zones2.globalid IS 'Original attribute from ArcGIS API is GlobalID';
COMMENT ON COLUMN fieldseeker.zones2.created_user IS 'Original attribute from ArcGIS API is created_user';
COMMENT ON COLUMN fieldseeker.zones2.created_date IS 'Original attribute from ArcGIS API is created_date';
COMMENT ON COLUMN fieldseeker.zones2.last_edited_user IS 'Original attribute from ArcGIS API is last_edited_user';
COMMENT ON COLUMN fieldseeker.zones2.last_edited_date IS 'Original attribute from ArcGIS API is last_edited_date';
COMMENT ON COLUMN fieldseeker.zones2.creationdate IS 'Original attribute from ArcGIS API is CreationDate';
COMMENT ON COLUMN fieldseeker.zones2.creator IS 'Original attribute from ArcGIS API is Creator';
COMMENT ON COLUMN fieldseeker.zones2.editdate IS 'Original attribute from ArcGIS API is EditDate';
COMMENT ON COLUMN fieldseeker.zones2.editor IS 'Original attribute from ArcGIS API is Editor';
COMMENT ON COLUMN fieldseeker.zones2.shape__area IS 'Original attribute from ArcGIS API is Shape__Area';
COMMENT ON COLUMN fieldseeker.zones2.shape__length IS 'Original attribute from ArcGIS API is Shape__Length';
-- +goose Down
DROP SCHEMA fieldseeker CASCADE;


