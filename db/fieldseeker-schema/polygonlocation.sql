-- Table definition for fieldseeker.PolygonLocation
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.polygonlocation_polygonlocation_usetype_e546154cb9544b9aa8e7b13e8e258b27_enum AS ENUM (
  'residential',
  'commercial',
  'industrial',
  'agricultural',
  'mixed_use',
  'public_domain',
  'natural',
  'municipal'
);

CREATE TYPE fieldseeker.polygonlocation_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.polygonlocation_locationsymbology_enum AS ENUM (
  'ACTION',
  'INACTIVE',
  'NONE'
);

CREATE TYPE fieldseeker.polygonlocation_polygonlocation_waterorigin_e9018e92_5f47_4ff9_8a7c_b818d848dc7a_enum AS ENUM (
  'flood_irrigation',
  'furrow_irrigation',
  'drip_irritation',
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

CREATE TYPE fieldseeker.polygonlocation_polygonlocation_habitat_45e9dde79ac84d959df8b65ba7d5dafd_enum AS ENUM (
  'orchard',
  'row_crops',
  'vine_crops',
  'ag_grasses_or_grain',
  'pasture',
  'irrigation_standpipe',
  'ditch',
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
  'low_point',
  'unknown'
);

CREATE TYPE fieldseeker.polygonlocation_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TABLE fieldseeker.polygonlocation (
  objectid BIGSERIAL NOT NULL,
  name VARCHAR(25),
  zone VARCHAR(25),
  habitat fieldseeker.polygonlocation_polygonlocation_habitat_45e9dde79ac84d959df8b65ba7d5dafd_enum,
  priority fieldseeker.polygonlocation_locationpriority_enum,
  usetype fieldseeker.polygonlocation_polygonlocation_usetype_e546154cb9544b9aa8e7b13e8e258b27_enum,
  active fieldseeker.polygonlocation_notinuit_f_enum,
  description VARCHAR(250),
  accessdesc VARCHAR(250),
  comments VARCHAR(250),
  symbology fieldseeker.polygonlocation_locationsymbology_enum,
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
  waterorigin fieldseeker.polygonlocation_polygonlocation_waterorigin_e9018e92_5f47_4ff9_8a7c_b818d848dc7a_enum,
  filter VARCHAR(255),
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  jurisdiction VARCHAR(25),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.polygonlocation.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.polygonlocation.name IS 'Name';
COMMENT ON COLUMN fieldseeker.polygonlocation.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.polygonlocation.habitat IS 'Habitat';
COMMENT ON COLUMN fieldseeker.polygonlocation.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.polygonlocation.usetype IS 'Use Type';
COMMENT ON COLUMN fieldseeker.polygonlocation.active IS 'Active';
COMMENT ON COLUMN fieldseeker.polygonlocation.description IS 'Description';
COMMENT ON COLUMN fieldseeker.polygonlocation.accessdesc IS 'Access Description';
COMMENT ON COLUMN fieldseeker.polygonlocation.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.polygonlocation.symbology IS 'Symbology';
COMMENT ON COLUMN fieldseeker.polygonlocation.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.polygonlocation.acres IS 'Acres';
COMMENT ON COLUMN fieldseeker.polygonlocation.nextactiondatescheduled IS 'Next Scheduled Action';
COMMENT ON COLUMN fieldseeker.polygonlocation.larvinspectinterval IS 'Larval Inspection Interval';
COMMENT ON COLUMN fieldseeker.polygonlocation.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectdate IS 'Last Inspection Date';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectbreeding IS 'Last Inspection Breeding';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectavglarvae IS 'Last Inspection Average Larvae';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectavgpupae IS 'Last Inspection Average Pupae';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectlstages IS 'Last Inspection Larval Stages';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectactiontaken IS 'Last Inspection Action';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectfieldspecies IS 'Last Inspection Field Species';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatdate IS 'Last Treatment Date';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatproduct IS 'Last Treatment Product';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatqty IS 'Last Treatment Quantity';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatqtyunit IS 'Last Treatment Quantity Unit';
COMMENT ON COLUMN fieldseeker.polygonlocation.hectares IS 'Hectares';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectactivity IS 'Last Inspection Activity';
COMMENT ON COLUMN fieldseeker.polygonlocation.lasttreatactivity IS 'Last Treatment Activity';
COMMENT ON COLUMN fieldseeker.polygonlocation.lastinspectconditions IS 'Last Inspection Conditions';
COMMENT ON COLUMN fieldseeker.polygonlocation.waterorigin IS 'Water Origin';
COMMENT ON COLUMN fieldseeker.polygonlocation.jurisdiction IS 'Jurisdiction';

-- Field active has default value: 1

-- Field symbology has default value: 'NONE'

-- See insert/insert_polygonlocation_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
