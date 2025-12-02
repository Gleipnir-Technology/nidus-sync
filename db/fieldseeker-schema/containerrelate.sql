-- Table definition for fieldseeker.ContainerRelate
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.containerrelate_mosquitocontainertype_enum AS ENUM (
  'Aquarium',
  'Flower pot',
  '5 gallon bucket',
  'Fountain',
  'Bird bath'
);

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
  containertype fieldseeker.containerrelate_mosquitocontainertype_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.containerrelate.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.containerrelate.containertype IS 'Container Type';

-- See insert/insert_containerrelate_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
