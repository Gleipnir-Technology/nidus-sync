-- Table definition for fieldseeker.HabitatRelate
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.habitatrelate_habitatrelate_habitattype_2e81cf2f550e400783cf284f3cec3953_enum AS ENUM (
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

CREATE TABLE fieldseeker.habitatrelate (
  objectid BIGSERIAL NOT NULL,
  foreign_id UUID,
  globalid UUID,
  created_user VARCHAR(255),
  created_date TIMESTAMP,
  last_edited_user VARCHAR(255),
  last_edited_date TIMESTAMP,
  habitattype fieldseeker.habitatrelate_habitatrelate_habitattype_2e81cf2f550e400783cf284f3cec3953_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.habitatrelate.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.habitatrelate.habitattype IS 'Habitat Type';

-- See insert/insert_habitatrelate_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
