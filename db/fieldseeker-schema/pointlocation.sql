-- Table definition for fieldseeker.PointLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.pointlocation_pointlocation_deactivate_reason_dd303085_b33c_4894_8c47_fa847dd9d7c5_enum AS ENUM (
  'Source Removed',
  'Pool Maintained',
  'Source Screened',
  'Crop Change',
  'Low or No Mosquito Activity',
  'Consistent Fish Presence'
);

CREATE TYPE fieldseeker.pointlocation_pointlocation_habitat_b4d8135a_4979_49c8_8bb3_67ec7230e661_enum AS ENUM (
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
  'marsh_or_wetlands',
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

CREATE TYPE fieldseeker.pointlocation_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TYPE fieldseeker.pointlocation_pointlocation_usetype_58d62d18ef4f47fc8cb9874df867f89e_enum AS ENUM (
  'residential',
  'commercial',
  'agricultural',
  'industrial',
  'mixed_use',
  'public_domain',
  'natural',
  'municipal'
);

CREATE TYPE fieldseeker.pointlocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.pointlocation_locationsymbology_enum AS ENUM (
  'ACTION',
  'INACTIVE',
  'NONE'
);

CREATE TYPE fieldseeker.pointlocation_pointlocation_waterorigin_197b22bf_f3eb_4dad_8899_986460f6ea97_enum AS ENUM (
  'flood_irrigation',
  'furrow_irrigation',
  'drip_irrigation',
  'sprinkler_irrigation',
  'wastewater_irrigation',
  'irrigation_runoff',
  'stormwater_or_municipal_runoff',
  'industrial_runoff',
  'rainwater_accumulation',
  'leak',
  'seepage',
  'stored_water',
  'wastewater_system',
  'permanent_natural_water',
  'temporary_natural_water',
  'recreational_or_orenamental_water',
  'water_conveyance'
);

CREATE TYPE fieldseeker.pointlocation_pointlocation_assignedtech_9393a162_2474_429d_85be_daa44e4c091f_enum AS ENUM (
  'Bryan Feguson',
  'Rick Alverez',
  'Alysia Davis',
  'Bryan Ruiz',
  'Kory Wilson',
  'Adrian Sifuentes',
  'Marco Martinez',
  'Carlos Rodriguez',
  'Landon McGill',
  'Ted McGill',
  'Mario Sanchez',
  'Jorge Perez',
  'Arturo Garcia-Trejo',
  'Lisa Salgado',
  'Lawrence Guzman',
  'Tricia Snowden',
  'Ryan Spratt',
  'Andrea Troupin',
  'Mark Nakata',
  'Pablo Ortega',
  'Benjamin Sperry',
  'Fatima Hidalgo',
  'Zackery Barragan',
  'Yajaira Godinez',
  'Jake Maldonado',
  'Rafael Ramirez',
  'Carlos Palacios',
  'Aaron Fredrick',
  'Josh Malone',
  'Alec Caposella',
  'Laura Ramos'
);

CREATE TABLE fieldseeker.pointlocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.pointlocation_pointlocation_habitat_b4d8135a_4979_49c8_8bb3_67ec7230e661_enum,
  priority fieldseeker.pointlocation_locationpriority_enum,
  usetype fieldseeker.pointlocation_pointlocation_usetype_58d62d18ef4f47fc8cb9874df867f89e_enum,
  active fieldseeker.pointlocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology fieldseeker.pointlocation_locationsymbology_enum,
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
  waterorigin fieldseeker.pointlocation_pointlocation_waterorigin_197b22bf_f3eb_4dad_8899_986460f6ea97_enum,
  x DOUBLE PRECISION,
  y DOUBLE PRECISION,
  assignedtech fieldseeker.pointlocation_pointlocation_assignedtech_9393a162_2474_429d_85be_daa44e4c091f_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  deactivate_reason fieldseeker.pointlocation_pointlocation_deactivate_reason_dd303085_b33c_4894_8c47_fa847dd9d7c5_enum,
  scalarpriority INTEGER,
  sourcestatus VARCHAR(255),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.pointlocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.pointlocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.pointlocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.pointlocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.pointlocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.pointlocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.pointlocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.pointlocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.pointlocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.pointlocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.pointlocation.symbology IS 'Symbology';
COMMENT ON COLUMN fieldseeker.pointlocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.pointlocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.pointlocation.larvinspectinterval IS 'Larval Inspection Interval';
COMMENT ON COLUMN fieldseeker.pointlocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.pointlocation.stype IS 'SourceType';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectdate IS 'Last Inspection Date';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectbreeding IS 'Last Inspection Breeding';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectavglarvae IS 'Last Inspection Average Larvae';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectavgpupae IS 'Last Inspection Average Pupae';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectlstages IS 'Last Inspection Larval Stages';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectactiontaken IS 'Last Inspection Action';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectfieldspecies IS 'Last Inspection Field Species';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatdate IS 'Last Treatment Date';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatproduct IS 'Last Treatment Product';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatqty IS 'Last Treatment Quantity';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatqtyunit IS 'Last Treatment Quantity Unit';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectactivity IS 'Last Inspection Activity';
COMMENT ON COLUMN fieldseeker.pointlocation.lasttreatactivity IS 'Last Treatment Activity';
COMMENT ON COLUMN fieldseeker.pointlocation.lastinspectconditions IS 'Last Inspection Conditions';
COMMENT ON COLUMN fieldseeker.pointlocation.waterorigin IS 'Water Origin';
COMMENT ON COLUMN fieldseeker.pointlocation.assignedtech IS 'Assigned Tech';
COMMENT ON COLUMN fieldseeker.pointlocation.jurisdiction IS 'Jurisdiction';
COMMENT ON COLUMN fieldseeker.pointlocation.deactivate_reason IS 'Reason for Deactivation';

-- Field active has default value: 1

-- Field symbology has default value: 'NONE'

-- See insert/insert_pointlocation_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
