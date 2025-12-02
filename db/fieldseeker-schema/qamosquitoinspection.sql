-- Table definition for fieldseeker.QAMosquitoInspection
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.qamosquitoinspection_mosquitoaction_enum AS ENUM (
  'Treatment',
  'Covered container',
  'Cleared debris',
  'Maintenance'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qasitetype_enum AS ENUM (
  'Detention Pond',
  'Ditch',
  'Low Area',
  'Mangrove Edge',
  'Pond',
  'Pond Edge',
  'Swale'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qabreedingpotential_enum AS ENUM (
  'High',
  'Low',
  'Medium',
  'Rare'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qamosquitohabitat_enum AS ENUM (
  'Depressions',
  'Detritus present',
  'Fast',
  'Few predators',
  'Fluctuating levels',
  'H20<6"',
  'Low wave potential',
  'No fish',
  'Shallow edges',
  'Still water edges',
  'Still water whole',
  'Veg. on edges'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qaaquaticorganisms_enum AS ENUM (
  'fish',
  'scuds',
  'snails'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qasoilcondition_enum AS ENUM (
  'Cracked',
  'Dry',
  'Inundated',
  'Saturated',
  'Surface Moist'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qawaterduration_enum AS ENUM (
  '~month',
  '~week',
  '<1 week',
  '<day',
  '<month',
  '>month',
  '>week'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qawaterconditions_enum AS ENUM (
  '"rust" material',
  'Clear',
  'Cloudy/fines',
  'Floating debris',
  'Submerged/decom. debris'
);

CREATE TYPE fieldseeker.qamosquitoinspection_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qalarvaereason_enum AS ENUM (
  'Missed Area',
  'New Site',
  'Not Visited',
  'Rate Low',
  'Treated Recently',
  'Unknown',
  'Wrong Product'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qavegetation_enum AS ENUM (
  'Algae',
  'Cattails',
  'Duckweed',
  'Glasswort',
  'Grass on edge',
  'Mangrove',
  'Mosquito fern',
  'Muskgrass',
  'Myriophyllum',
  'Other',
  'Rotting vegetation',
  'Saltwort',
  'Sedges'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qasourcereduction_enum AS ENUM (
  '1 tractor < day',
  'adjust flood irrigation',
  'adjust turf irrigation',
  'clear outflow',
  'cut ditch',
  'hand grading',
  'laser leveling',
  'multiple loads soil',
  'none'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qawatermovement_enum AS ENUM (
  'Fast',
  'Medium',
  'None',
  'Slow',
  'Very Slow'
);

CREATE TYPE fieldseeker.qamosquitoinspection_qawatersource_enum AS ENUM (
  'Irrigation',
  'Manually Controlled',
  'Percolation',
  'Rain Runoff',
  'Tidal',
  'Water Table'
);

CREATE TABLE fieldseeker.qamosquitoinspection (
  objectid BIGSERIAL NOT NULL,
  posdips SMALLINT,
  actiontaken fieldseeker.qamosquitoinspection_mosquitoaction_enum,
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  raingauge DOUBLE PRECISION,
  globalid UUID,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  winddir VARCHAR(3),
  reviewed fieldseeker.qamosquitoinspection_notinuit_f_enum,
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
  fish fieldseeker.qamosquitoinspection_notinuit_f_enum,
  sitetype fieldseeker.qamosquitoinspection_qasitetype_enum,
  breedingpotential fieldseeker.qamosquitoinspection_qabreedingpotential_enum,
  movingwater fieldseeker.qamosquitoinspection_notinuit_f_enum,
  nowaterever fieldseeker.qamosquitoinspection_notinuit_f_enum,
  mosquitohabitat fieldseeker.qamosquitoinspection_qamosquitohabitat_enum,
  habvalue1 SMALLINT,
  habvalue1percent SMALLINT,
  habvalue2 SMALLINT,
  habvalue2percent SMALLINT,
  potential SMALLINT,
  larvaepresent fieldseeker.qamosquitoinspection_notinuit_f_enum,
  larvaeinsidetreatedarea fieldseeker.qamosquitoinspection_notinuit_f_enum,
  larvaeoutsidetreatedarea fieldseeker.qamosquitoinspection_notinuit_f_enum,
  larvaereason fieldseeker.qamosquitoinspection_qalarvaereason_enum,
  aquaticorganisms fieldseeker.qamosquitoinspection_qaaquaticorganisms_enum,
  vegetation fieldseeker.qamosquitoinspection_qavegetation_enum,
  sourcereduction fieldseeker.qamosquitoinspection_qasourcereduction_enum,
  waterpresent fieldseeker.qamosquitoinspection_notinuit_f_enum,
  watermovement1 fieldseeker.qamosquitoinspection_qawatermovement_enum,
  watermovement1percent SMALLINT,
  watermovement2 fieldseeker.qamosquitoinspection_qawatermovement_enum,
  watermovement2percent SMALLINT,
  soilconditions fieldseeker.qamosquitoinspection_qasoilcondition_enum,
  waterduration fieldseeker.qamosquitoinspection_qawaterduration_enum,
  watersource fieldseeker.qamosquitoinspection_qawatersource_enum,
  waterconditions fieldseeker.qamosquitoinspection_qawaterconditions_enum,
  adultactivity fieldseeker.qamosquitoinspection_notinuit_f_enum,
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
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.qamosquitoinspection.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.qamosquitoinspection.posdips IS 'Positive Dips';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.actiontaken IS 'Action';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.avetemp IS 'Average Temperature';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.windspeed IS 'Wind Speed';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.raingauge IS 'Rain Gauge';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.winddir IS 'Wind Direction';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.lr IS 'Landing Rate';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.negdips IS 'Negative Dips';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.totalacres IS 'Total Acres';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.acresbreeding IS 'Acres Breeding';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.fish IS 'Fish Present?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.sitetype IS 'Site Type';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.breedingpotential IS 'Breeding Potential';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.movingwater IS 'Moving Water';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.nowaterever IS 'No Evidence of Water Ever';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.mosquitohabitat IS 'Mosquito Habitat Indicators';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue1 IS 'Habitat Value';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue1percent IS '%';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue2 IS 'Habitat Value';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.habvalue2percent IS '%';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.potential IS 'Potential';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaepresent IS 'Larvae Present';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaeinsidetreatedarea IS 'Larvae Inside Treated Area?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaeoutsidetreatedarea IS 'Larvae Outside Treated Area?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.larvaereason IS 'Reason Larvae Present';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.aquaticorganisms IS 'Aquatic Organisms';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.vegetation IS 'Vegetation';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.sourcereduction IS 'Source Reduction';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterpresent IS 'Water Present?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement1 IS 'Water Movement';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement1percent IS '%';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement2 IS 'Water Movement';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watermovement2percent IS '%';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.soilconditions IS 'Soil Conditions';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterduration IS 'How Long Water Present?';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.watersource IS 'Water Source';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.waterconditions IS 'Water Conditions';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.adultactivity IS 'Adult Activity';
COMMENT ON COLUMN fieldseeker.qamosquitoinspection.fieldtech IS 'Field Tech';

-- See insert/insert_qamosquitoinspection_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
