-- Table definition for fieldseeker.ProposedTreatmentArea
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.proposedtreatmentarea_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TYPE fieldseeker.proposedtreatmentarea_locationpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'None'
);

CREATE TYPE fieldseeker.proposedtreatmentarea_mosquitotreatmentmethod_enum AS ENUM (
  'Argo',
  'ATV',
  'Backpack',
  'Drone',
  'Manual',
  'Truck',
  'ULV',
  'Enhanced_Surveillance'
);

CREATE TABLE fieldseeker.proposedtreatmentarea (
  objectid BIGSERIAL NOT NULL,
  method fieldseeker.proposedtreatmentarea_mosquitotreatmentmethod_enum,
  comments VARCHAR(250),
  zone VARCHAR(25),
  reviewed fieldseeker.proposedtreatmentarea_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  zone2 VARCHAR(25),
  completeddate TIMESTAMP,
  completedby VARCHAR(25),
  completed fieldseeker.proposedtreatmentarea_notinuit_f_enum,
  issprayroute fieldseeker.proposedtreatmentarea_notinuit_f_enum,
  name VARCHAR(25),
  acres DOUBLE PRECISION,
  globalid UUID,
  exported fieldseeker.proposedtreatmentarea_notinuit_f_enum,
  targetproduct VARCHAR(25),
  targetapprate DOUBLE PRECISION,
  hectares DOUBLE PRECISION,
  lasttreatactivity VARCHAR(25),
  lasttreatdate TIMESTAMP,
  lasttreatproduct VARCHAR(25),
  lasttreatqty DOUBLE PRECISION,
  lasttreatqtyunit VARCHAR(10),
  priority fieldseeker.proposedtreatmentarea_locationpriority_enum,
  duedate TIMESTAMP,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  targetspecies VARCHAR(250),
  shape__area DOUBLE PRECISION,
  shape__length DOUBLE PRECISION,
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.method IS 'Method';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.zone2 IS 'Zone2';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completeddate IS 'Completed Date';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completedby IS 'Completed By';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.completed IS 'Completed';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.issprayroute IS 'Is Spray Route';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.name IS 'Name';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.acres IS 'Acres';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetproduct IS 'Target Product';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetapprate IS 'Target App Rate';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.hectares IS 'Hectares';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatactivity IS 'Last Treatment Activity';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatdate IS 'Last Treatment Date';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatproduct IS 'Last Treatment Product';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatqty IS 'Last Treatment Quantity';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.lasttreatqtyunit IS 'Last Treatment Quantity Unit';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.duedate IS 'Due Date';
COMMENT ON COLUMN fieldseeker.proposedtreatmentarea.targetspecies IS 'Target Species';

-- See insert/insert_proposedtreatmentarea_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
