-- Table definition for fieldseeker.ServiceRequest
-- Includes versioning for tracking changes

CREATE SCHEMA IF NOT EXISTS fieldseeker;

CREATE TYPE fieldseeker.servicerequest_servicerequestnextaction_enum AS ENUM (
  'Night spray',
  'Site visit'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_schedule_period_3f40c046_afd1_4abd_8bf4_389650d29a49_enum AS ENUM (
  'AM',
  'PM'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestsource_enum AS ENUM (
  'Phone',
  'Email',
  'Website',
  'Drop-in',
  '2025_pools'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_supervisor_eba07b90_c885_4fe6_8080_7aa775403b9a_enum AS ENUM (
  'Rick Alverez',
  'Bryan Ferguson',
  'Bryan Ruiz',
  'Andrea Troupin',
  'Conlin Reis'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestregion_enum AS ENUM (
  'FL',
  'ID'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_spanish_aaa3dc66_9f9a_4527_8ecd_c9f76db33879_enum AS ENUM (
  '0',
  '1'
);

CREATE TYPE fieldseeker.servicerequest_servicerequeststatus_enum AS ENUM (
  'Assigned',
  'Closed',
  'FieldRectified',
  'Open',
  'Rejected',
  'Unverified',
  'Accepted'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestcontactpreferences_enum AS ENUM (
  'None',
  'Call',
  'Email',
  'Text'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestrejectedreason_enum AS ENUM (
  'Distance',
  'Workload'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_dog_2b95ec97_1286_4fcd_88f4_f0e31113f696_enum AS ENUM (
  '0',
  '1',
  '2',
  '3'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestissues_enum AS ENUM (
  'Beehive Related',
  'Unsanitary Accumulations',
  'Rooster or Noise',
  'Rats Attracted',
  'Odor',
  'Number of Animals Over Limit',
  'Location',
  'Violation',
  'Inadequate Enclosure',
  'Escaped Animal',
  'Illegal Animal'
);

CREATE TYPE fieldseeker.servicerequest_servicerequestpriority_enum AS ENUM (
  'Low',
  'Medium',
  'High',
  'Follow up Visit',
  'HTC Response',
  'Disease Activity Response'
);

CREATE TYPE fieldseeker.servicerequest_servicerequesttarget_enum AS ENUM (
  'mosquitofish',
  'neglected pool or spa',
  'standing water',
  'mosquito presence',
  'biting mosquitoes',
  'event',
  'fish',
  'mosquito',
  'source',
  'bird'
);

CREATE TYPE fieldseeker.servicerequest_servicerequest_assignedtech_71d0d685_868f_4b7a_87e2_3661a3ee67c5_enum AS ENUM (
  'Alysia Davis',
  'Alejandra Gill',
  'Andrea Troupin',
  'Brenda Rodriguez',
  'Bryan Ferguson',
  'Bryan Ruiz',
  'Conlin Reis',
  'Carlos Rodriguez',
  'Erick Arriga',
  'Landon McGill',
  'Marco Martinez',
  'Mark Nakata',
  'Mario Sanchez',
  'Juan Pablo Ortega',
  'Ryan Spratt',
  'Ted McGill',
  'Benjamin Sperry',
  'Zachery Barragan',
  'Arturo Garcia-Trejo',
  'Jesus Jolano',
  'Yajaira Godinez',
  'Jake Maldonado',
  'Rafael Ramirez',
  'Lisa Salgado',
  'Kory Wilson',
  'Carlos Palacios',
  'Fatima Hidalgo',
  'Aaron Fredrick',
  'Josh Malone',
  'Jorge Perez',
  'Laura Romos'
);

CREATE TYPE fieldseeker.servicerequest_notinuit_f_enum AS ENUM (
  '1',
  '0'
);

CREATE TABLE fieldseeker.servicerequest (
  objectid BIGSERIAL NOT NULL,
  recdatetime TIMESTAMP,
  source fieldseeker.servicerequest_servicerequestsource_enum,
  entrytech VARCHAR(25),
  priority fieldseeker.servicerequest_servicerequestpriority_enum,
  supervisor fieldseeker.servicerequest_servicerequest_supervisor_eba07b90_c885_4fe6_8080_7aa775403b9a_enum,
  assignedtech fieldseeker.servicerequest_servicerequest_assignedtech_71d0d685_868f_4b7a_87e2_3661a3ee67c5_enum,
  status fieldseeker.servicerequest_servicerequeststatus_enum,
  clranon fieldseeker.servicerequest_notinuit_f_enum,
  clrfname VARCHAR(25),
  clrphone1 VARCHAR(25),
  clrphone2 VARCHAR(25),
  clremail VARCHAR(50),
  clrcompany VARCHAR(25),
  clraddr1 VARCHAR(50),
  clraddr2 VARCHAR(50),
  clrcity VARCHAR(25),
  clrstate fieldseeker.servicerequest_servicerequestregion_enum,
  clrzip VARCHAR(25),
  clrother VARCHAR(25),
  clrcontpref fieldseeker.servicerequest_servicerequestcontactpreferences_enum,
  reqcompany VARCHAR(25),
  reqaddr1 VARCHAR(50),
  reqaddr2 VARCHAR(50),
  reqcity VARCHAR(25),
  reqstate fieldseeker.servicerequest_servicerequestregion_enum,
  reqzip VARCHAR(25),
  reqcrossst VARCHAR(25),
  reqsubdiv VARCHAR(25),
  reqmapgrid VARCHAR(25),
  reqpermission fieldseeker.servicerequest_notinuit_f_enum,
  reqtarget fieldseeker.servicerequest_servicerequesttarget_enum,
  reqdescr VARCHAR(1000),
  reqnotesfortech VARCHAR(250),
  reqnotesforcust VARCHAR(250),
  reqfldnotes VARCHAR(250),
  reqprogramactions VARCHAR(250),
  datetimeclosed TIMESTAMP,
  techclosed VARCHAR(25),
  sr_number INTEGER,
  reviewed fieldseeker.servicerequest_notinuit_f_enum,
  reviewedby VARCHAR(25),
  revieweddate TIMESTAMP,
  accepted fieldseeker.servicerequest_notinuit_f_enum,
  accepteddate TIMESTAMP,
  rejectedby VARCHAR(25),
  rejecteddate TIMESTAMP,
  rejectedreason fieldseeker.servicerequest_servicerequestrejectedreason_enum,
  duedate TIMESTAMP,
  acceptedby VARCHAR(25),
  comments VARCHAR(2500),
  estcompletedate TIMESTAMP,
  nextaction fieldseeker.servicerequest_servicerequestnextaction_enum,
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
  dog fieldseeker.servicerequest_servicerequest_dog_2b95ec97_1286_4fcd_88f4_f0e31113f696_enum,
  schedule_period fieldseeker.servicerequest_servicerequest_schedule_period_3f40c046_afd1_4abd_8bf4_389650d29a49_enum,
  schedule_notes VARCHAR(256),
  spanish fieldseeker.servicerequest_servicerequest_spanish_aaa3dc66_9f9a_4527_8ecd_c9f76db33879_enum,
  creationdate TIMESTAMP,
  creator VARCHAR(128),
  editdate TIMESTAMP,
  editor VARCHAR(128),
  issuesreported fieldseeker.servicerequest_servicerequestissues_enum,
  jurisdiction VARCHAR(25),
  notificationtimestamp VARCHAR(250),
  zone VARCHAR(50),
  zone2 VARCHAR(50),
  VERSION INTEGER NOT NULL DEFAULT 1,
  PRIMARY KEY (objectid, VERSION)
);

COMMENT ON COLUMN fieldseeker.servicerequest.VERSION IS 'Tracks version changes to the row. Increases when data is modified.';

COMMENT ON COLUMN fieldseeker.servicerequest.recdatetime IS 'Received';
COMMENT ON COLUMN fieldseeker.servicerequest.source IS 'Source';
COMMENT ON COLUMN fieldseeker.servicerequest.entrytech IS 'Entered By';
COMMENT ON COLUMN fieldseeker.servicerequest.priority IS 'Priority';
COMMENT ON COLUMN fieldseeker.servicerequest.supervisor IS 'Supervisor';
COMMENT ON COLUMN fieldseeker.servicerequest.assignedtech IS 'Assigned To';
COMMENT ON COLUMN fieldseeker.servicerequest.status IS 'Status';
COMMENT ON COLUMN fieldseeker.servicerequest.clranon IS 'Anonymous Caller';
COMMENT ON COLUMN fieldseeker.servicerequest.clrfname IS 'Caller Name';
COMMENT ON COLUMN fieldseeker.servicerequest.clrphone1 IS 'Caller Phone';
COMMENT ON COLUMN fieldseeker.servicerequest.clrphone2 IS 'Caller Alternate Phone';
COMMENT ON COLUMN fieldseeker.servicerequest.clremail IS 'Caller Email';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcompany IS 'Caller Company';
COMMENT ON COLUMN fieldseeker.servicerequest.clraddr1 IS 'Caller Address';
COMMENT ON COLUMN fieldseeker.servicerequest.clraddr2 IS 'Caller Address 2';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcity IS 'Caller City';
COMMENT ON COLUMN fieldseeker.servicerequest.clrstate IS 'Caller State';
COMMENT ON COLUMN fieldseeker.servicerequest.clrzip IS 'Caller ZIP';
COMMENT ON COLUMN fieldseeker.servicerequest.clrother IS 'Caller Other';
COMMENT ON COLUMN fieldseeker.servicerequest.clrcontpref IS 'Caller Contact Preference';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcompany IS 'Request Company';
COMMENT ON COLUMN fieldseeker.servicerequest.reqaddr1 IS 'Request Address';
COMMENT ON COLUMN fieldseeker.servicerequest.reqaddr2 IS 'Request Address 2';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcity IS 'Request City';
COMMENT ON COLUMN fieldseeker.servicerequest.reqstate IS 'Request State';
COMMENT ON COLUMN fieldseeker.servicerequest.reqzip IS 'Request ZIP';
COMMENT ON COLUMN fieldseeker.servicerequest.reqcrossst IS 'Request Cross Street';
COMMENT ON COLUMN fieldseeker.servicerequest.reqsubdiv IS 'Request Subdivision';
COMMENT ON COLUMN fieldseeker.servicerequest.reqmapgrid IS 'Request Map Grid';
COMMENT ON COLUMN fieldseeker.servicerequest.reqpermission IS 'Permission to Enter';
COMMENT ON COLUMN fieldseeker.servicerequest.reqtarget IS 'Request Target';
COMMENT ON COLUMN fieldseeker.servicerequest.reqdescr IS 'Request Description';
COMMENT ON COLUMN fieldseeker.servicerequest.reqnotesfortech IS 'Notes for Field Technician';
COMMENT ON COLUMN fieldseeker.servicerequest.reqnotesforcust IS 'Notes for Customer';
COMMENT ON COLUMN fieldseeker.servicerequest.reqfldnotes IS 'Request Field Notes';
COMMENT ON COLUMN fieldseeker.servicerequest.reqprogramactions IS 'Request Program Actions';
COMMENT ON COLUMN fieldseeker.servicerequest.datetimeclosed IS 'Closed';
COMMENT ON COLUMN fieldseeker.servicerequest.techclosed IS 'Closed By';
COMMENT ON COLUMN fieldseeker.servicerequest.sr_number IS 'SR#';
COMMENT ON COLUMN fieldseeker.servicerequest.reviewed IS 'Reviewed';
COMMENT ON COLUMN fieldseeker.servicerequest.reviewedby IS 'Reviewed By';
COMMENT ON COLUMN fieldseeker.servicerequest.revieweddate IS 'Reviewed Date';
COMMENT ON COLUMN fieldseeker.servicerequest.accepted IS 'Accepted';
COMMENT ON COLUMN fieldseeker.servicerequest.accepteddate IS 'Accepted Date';
COMMENT ON COLUMN fieldseeker.servicerequest.rejectedby IS 'Rejected By';
COMMENT ON COLUMN fieldseeker.servicerequest.rejecteddate IS 'Rejected Date';
COMMENT ON COLUMN fieldseeker.servicerequest.rejectedreason IS 'Rejected Reason';
COMMENT ON COLUMN fieldseeker.servicerequest.duedate IS 'Due Date';
COMMENT ON COLUMN fieldseeker.servicerequest.acceptedby IS 'Accepted By';
COMMENT ON COLUMN fieldseeker.servicerequest.comments IS 'Comments';
COMMENT ON COLUMN fieldseeker.servicerequest.estcompletedate IS 'Estimated Completion Date';
COMMENT ON COLUMN fieldseeker.servicerequest.nextaction IS 'Next Action';
COMMENT ON COLUMN fieldseeker.servicerequest.recordstatus IS 'Record Status';
COMMENT ON COLUMN fieldseeker.servicerequest.firstresponsedate IS 'First Response Date';
COMMENT ON COLUMN fieldseeker.servicerequest.responsedaycount IS 'Response Day Count';
COMMENT ON COLUMN fieldseeker.servicerequest.allowed IS 'Verify Correct Location';
COMMENT ON COLUMN fieldseeker.servicerequest.xvalue IS 'Xvalue';
COMMENT ON COLUMN fieldseeker.servicerequest.yvalue IS 'Yvalue';
COMMENT ON COLUMN fieldseeker.servicerequest.validx IS 'ValidX';
COMMENT ON COLUMN fieldseeker.servicerequest.validy IS 'ValidY';
COMMENT ON COLUMN fieldseeker.servicerequest.externalid IS 'External ID';
COMMENT ON COLUMN fieldseeker.servicerequest.externalerror IS 'External Error';
COMMENT ON COLUMN fieldseeker.servicerequest.notified IS 'Notified';
COMMENT ON COLUMN fieldseeker.servicerequest.notifieddate IS 'Notified Date';
COMMENT ON COLUMN fieldseeker.servicerequest.scheduled IS 'Scheduled';
COMMENT ON COLUMN fieldseeker.servicerequest.scheduleddate IS 'Scheduled Date';
COMMENT ON COLUMN fieldseeker.servicerequest.dog IS 'Dog';
COMMENT ON COLUMN fieldseeker.servicerequest.schedule_period IS 'Schedule Period';
COMMENT ON COLUMN fieldseeker.servicerequest.schedule_notes IS 'Schedule Notes';
COMMENT ON COLUMN fieldseeker.servicerequest.spanish IS 'Prefer speaking Spanish';
COMMENT ON COLUMN fieldseeker.servicerequest.issuesreported IS 'Issues Reported';
COMMENT ON COLUMN fieldseeker.servicerequest.jurisdiction IS 'Jurisdiction';
COMMENT ON COLUMN fieldseeker.servicerequest.notificationtimestamp IS 'Notification Timestamp';
COMMENT ON COLUMN fieldseeker.servicerequest.zone IS 'Zone';
COMMENT ON COLUMN fieldseeker.servicerequest.zone2 IS 'Zone2';

-- Field dog has default value: 0

-- Field spanish has default value: 0

-- See insert/insert_servicerequest_versioned.sql for prepared insert statement

-- Usage notes for versioning:
-- When inserting a new row, VERSION defaults to 1
-- When updating a row, insert a new row with the same ID but incremented VERSION
-- The most recent version of a row has the highest VERSION value
