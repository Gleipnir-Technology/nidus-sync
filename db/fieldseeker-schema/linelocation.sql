-- Table definition for fieldseeker.LineLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.linelocation_linelocation_habitat_fc51bdc4f1954df58206d69ce14182f3_enum AS ENUM (
  'orchard',
  'row_crops',
  'vine_crops',
  'ag_grasses_or_grain',
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
  'utility',
  'catch_basin',
  'stream_or_creek',
  'slough',
  'river',
  'marsh_or_wetland',
  'containers',
  'watering_bowl',
  'plant_saucer',
  'yard_drain',
  'plant_axil',
  'treehole',
  'fountain_or_water_feature',
  'bird_bath',
  'misc_water_accumulation',
  'tarp_or_cover',
  'swimming_pool',
  'aboveground_pool',
  'kid_pool',
  'hot_tub',
  'appliance',
  'tires',
  'flooded_structure',
  'low_point'
);

CREATE TYPE fieldseeker.linelocation_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TYPE fieldseeker.linelocation_linelocation_usetype_2aeca2e60d2f455c86fc34895dc80a02_enum AS ENUM (
  'residential',
  'commercial',
  'industrial',
  'agricultural',
  'mixed_use',
  'public_domain',
  'natural',
  'municipal'
);

CREATE TYPE fieldseeker.linelocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.linelocation_locationsymbology_enum AS ENUM (
  'ACTION',
  'INACTIVE',
  'NONE'
);

CREATE TYPE fieldseeker.linelocation_linelocation_waterorigin_84723d92_306a_46f4_8ef1_69b55a916008_enum AS ENUM (
  'flood_irrigation',
  'furrow_irrigation',
  'drip_irrigation',
  'sprinkler_irrigation',
  'wastewater_irrigation',
  'irrigation_runoff',
  'rainwater_accumulation',
  'leak',
  'seepage',
  'stored_water',
  'wastewater_system',
  'permanent_natural_water',
  'temporary_natural_water',
  'recreational_or_ornamental_water',
  'water_conveyance'
);

CREATE TABLE fieldseeker.linelocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.linelocation_linelocation_habitat_fc51bdc4f1954df58206d69ce14182f3_enum,
  priority fieldseeker.linelocation_locationpriority_enum,
  usetype fieldseeker.linelocation_linelocation_usetype_2aeca2e60d2f455c86fc34895dc80a02_enum,
  active fieldseeker.linelocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology fieldseeker.linelocation_locationsymbology_enum,
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
  waterorigin fieldseeker.linelocation_linelocation_waterorigin_84723d92_306a_46f4_8ef1_69b55a916008_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.linelocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.linelocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.linelocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.linelocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.linelocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.linelocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.linelocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.linelocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.linelocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.linelocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.linelocation.symbology IS 'Symbology';
COMMENT ON COLUMN fieldseeker.linelocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.linelocation.acres IS 'Acres';
COMMENT ON COLUMN fieldseeker.linelocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.linelocation.larvinspectinterval IS 'Larval Inspection Interval';
COMMENT ON COLUMN fieldseeker.linelocation.length_ft IS 'Length';
COMMENT ON COLUMN fieldseeker.linelocation.width_ft IS 'Width';
COMMENT ON COLUMN fieldseeker.linelocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectdate IS 'Last Inspection Date';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectbreeding IS 'Last Inspection Breeding';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectavglarvae IS 'Last Inspection Average Larvae';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectavgpupae IS 'Last Inspection Average Pupae';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectlstages IS 'Last Inspection Larval Stages';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectactiontaken IS 'Last Inspection Action';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectfieldspecies IS 'Last Inspection Field Species';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatdate IS 'Last Treatment Date';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatproduct IS 'Last Treatment Product';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatqty IS 'Last Treatment Quantity';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatqtyunit IS 'Last Treatment Quantity Unit';
COMMENT ON COLUMN fieldseeker.linelocation.hectares IS 'Hectares';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectactivity IS 'Last Inspection Activity';
COMMENT ON COLUMN fieldseeker.linelocation.lasttreatactivity IS 'Last Treatment Activity';
COMMENT ON COLUMN fieldseeker.linelocation.length_meters IS 'Length Meters';
COMMENT ON COLUMN fieldseeker.linelocation.width_meters IS 'Width Meters';
COMMENT ON COLUMN fieldseeker.linelocation.lastinspectconditions IS 'Last Inspection Conditions';
COMMENT ON COLUMN fieldseeker.linelocation.waterorigin IS 'Water Origin';
COMMENT ON COLUMN fieldseeker.linelocation.jurisdiction IS 'Jurisdiction';

-- Field active has default value: 1

-- Field symbology has default value: 'NONE'

-- See insert/insert_linelocation_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
