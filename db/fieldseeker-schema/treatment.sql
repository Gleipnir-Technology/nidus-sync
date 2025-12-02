-- Table definition for fieldseeker.Treatment
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.treatment_treatment_method_d558ca3ccf43440c8160758253967621_enum AS ENUM (
  'Argo',
  'ATV',
  'Backpack',
  'Drone',
  'Manual',
  'Truck',
  'ULV',
  'WALS',
  'Administrative Action'
);

CREATE TYPE fieldseeker.treatment_treatment_equiptype_45694d79_ff21_42cc_be4f_a0d1def4fba0_enum AS ENUM (
  'Backpack #1',
  'A1  Mist Sprayer (T-3) ',
  'Spreader #2',
  'Guardian #73 ',
  'ULV #74 (Grizzly)',
  'Clark ULV Sprayer #71',
  'Clark ULV Sprayer #72',
  'Spray bottle'
);

CREATE TYPE fieldseeker.treatment_notinuiwinddirection_enum AS ENUM (
  'N',
  'NE',
  'E',
  'SE',
  'S',
  'SW',
  'W',
  'NW'
);

CREATE TYPE fieldseeker.treatment_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.treatment_treatment_habitat_0afee7eb_f9ea_4707_8483_cccfe60f0d16_enum AS ENUM (
  'orchard',
  'row_crops',
  'vine_crops',
  'ag_grass_or_grain',
  'pasture',
  'irrigation_standpipe',
  'ditch',
  'pond',
  'sump',
  'drain',
  'dairy_lagoon',
  'wastewater_treatment',
  'trough',
  'depression',
  'gutter',
  'rain_gutter',
  'culvert',
  'Utility',
  'catch_basin',
  'stream_or_creek',
  'slough',
  'river',
  'marsh_or_wetlands',
  'containers',
  'watering_bowl',
  'plant_saucer',
  'yard_drain',
  'plant_axil',
  'treehole',
  'foutain_or_water_feature',
  'bird_bath',
  'misc_water_accumulation',
  'tarp_or_cover',
  'swimming_pool',
  'aboveground_pool',
  'kid_pool',
  'hot_tub',
  'applicance',
  'flooded_structure',
  'low_point'
);

CREATE TYPE fieldseeker.treatment_treatment_sitecond_f812e1f64dcb4dc9a75da9d00abe6169_enum AS ENUM (
  'Dry',
  'Flowing',
  'Maintained',
  'Unmaintained',
  'High Organic',
  'Fish Present'
);

CREATE TYPE fieldseeker.treatment_treatment_sitecond_5a15bf36fa124280b961f31cd1a9b571_enum AS ENUM (
  'Dry',
  'Flowing',
  'Maintained',
  'Unmaintained',
  'High Organic',
  'Fish Present',
  'Stagnant'
);

CREATE TYPE fieldseeker.treatment_mosquitoactivity_enum AS ENUM (
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

CREATE TYPE fieldseeker.treatment_mosquitoproductareaunit_enum AS ENUM (
  'acre',
  'sq ft'
);

CREATE TYPE fieldseeker.treatment_mosquitoproductmeasureunit_enum AS ENUM (
  'briquet',
  'dry oz',
  'each',
  'fl oz',
  'gal',
  'lb',
  'packet',
  'pouch'
);

CREATE TABLE fieldseeker.treatment (
  objectid BIGSERIAL NOT NULL,
  activity fieldseeker.treatment_mosquitoactivity_enum,
  treatarea DOUBLE PRECISION,
  areaunit fieldseeker.treatment_mosquitoproductareaunit_enum,
  product VARCHAR(25),
  qty DOUBLE PRECISION,
  qtyunit fieldseeker.treatment_mosquitoproductmeasureunit_enum,
  method fieldseeker.treatment_treatment_method_d558ca3ccf43440c8160758253967621_enum,
  equiptype fieldseeker.treatment_treatment_equiptype_45694d79_ff21_42cc_be4f_a0d1def4fba0_enum,
  comments VARCHAR(250),
  avetemp DOUBLE PRECISION,
  windspeed DOUBLE PRECISION,
  winddir fieldseeker.treatment_notinuiwinddirection_enum,
  raingauge DOUBLE PRECISION,
  startdatetime TIMESTAMP,
  enddatetime TIMESTAMP,
  insp_id UUID,
  reviewed fieldseeker.treatment_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  locationname VARCHAR(25),
  zone VARCHAR(25),
  warningoverride fieldseeker.treatment_notinuit_f_enum,
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
  habitat fieldseeker.treatment_treatment_habitat_0afee7eb_f9ea_4707_8483_cccfe60f0d16_enum,
  treathectares DOUBLE PRECISION,
  invloc VARCHAR(25),
  temp_sitecond fieldseeker.treatment_treatment_sitecond_f812e1f64dcb4dc9a75da9d00abe6169_enum,
  sitecond fieldseeker.treatment_treatment_sitecond_5a15bf36fa124280b961f31cd1a9b571_enum,
  totalcostprodcut DOUBLE PRECISION,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  targetspecies VARCHAR(250),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.treatment.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.treatment.activity IS 'Activity';
COMMENT ON COLUMN fieldseeker.treatment.treatarea IS 'Area Treated';
COMMENT ON COLUMN fieldseeker.treatment.areaunit IS 'Area Unit';
COMMENT ON COLUMN fieldseeker.treatment.product IS 'Product';
COMMENT ON COLUMN fieldseeker.treatment.qty IS 'Quantity';
COMMENT ON COLUMN fieldseeker.treatment.qtyunit IS 'Quantity Unit';
COMMENT ON COLUMN fieldseeker.treatment.method IS 'Method';
COMMENT ON COLUMN fieldseeker.treatment.equiptype IS 'Equipment Type';
COMMENT ON COLUMN fieldseeker.treatment.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.treatment.avetemp IS 'Average Temperature';
COMMENT ON COLUMN fieldseeker.treatment.windspeed IS 'Wind Speed';
COMMENT ON COLUMN fieldseeker.treatment.winddir IS 'Wind Direction';
COMMENT ON COLUMN fieldseeker.treatment.raingauge IS 'Rain Gauge';
COMMENT ON COLUMN fieldseeker.treatment.startdatetime IS 'Start';
COMMENT ON COLUMN fieldseeker.treatment.enddatetime IS 'Finish';
COMMENT ON COLUMN fieldseeker.treatment.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.treatment.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.treatment.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.treatment.locationname IS 'Location Name';
COMMENT ON COLUMN fieldseeker.treatment.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.treatment.warningoverride IS 'Warning Override';
COMMENT ON COLUMN fieldseeker.treatment.recordstatus IS 'RecordStatus';
COMMENT ON COLUMN fieldseeker.treatment.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.treatment.treatacres IS 'Treated Acres';
COMMENT ON COLUMN fieldseeker.treatment.tirecount IS 'Tire Count';
COMMENT ON COLUMN fieldseeker.treatment.cbcount IS 'Catch Basin Count';
COMMENT ON COLUMN fieldseeker.treatment.containercount IS 'Container Count';
COMMENT ON COLUMN fieldseeker.treatment.treatmentlength IS 'Treatment Length';
COMMENT ON COLUMN fieldseeker.treatment.treatmenthours IS 'Treatment Hours';
COMMENT ON COLUMN fieldseeker.treatment.treatmentlengthunits IS 'Treatment Length Units';
COMMENT ON COLUMN fieldseeker.treatment.fieldtech IS 'Field Tech';
COMMENT ON COLUMN fieldseeker.treatment.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.treatment.treathectares IS 'Treat Hectares';
COMMENT ON COLUMN fieldseeker.treatment.invloc IS 'Inventory Location';
COMMENT ON COLUMN fieldseeker.treatment.temp_sitecond IS 'temp_Conditions';
COMMENT ON COLUMN fieldseeker.treatment.sitecond IS 'Conditions';
COMMENT ON COLUMN fieldseeker.treatment.totalcostprodcut IS 'TotalCostProduct';
COMMENT ON COLUMN fieldseeker.treatment.targetspecies IS 'Target Species';

-- See insert/insert_treatment_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
