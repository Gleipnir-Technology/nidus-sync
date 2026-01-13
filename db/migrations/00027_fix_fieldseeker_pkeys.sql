-- +goose Up
ALTER TABLE fieldseeker.containerrelate DROP CONSTRAINT containerrelate_pkey;
ALTER TABLE fieldseeker.containerrelate ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.containerrelate_objectid_seq;
ALTER TABLE fieldseeker.containerrelate ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.fieldscoutinglog DROP CONSTRAINT fieldscoutinglog_pkey;
ALTER TABLE fieldseeker.fieldscoutinglog ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.fieldscoutinglog_objectid_seq;
ALTER TABLE fieldseeker.fieldscoutinglog ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.habitatrelate DROP CONSTRAINT habitatrelate_pkey;
ALTER TABLE fieldseeker.habitatrelate ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.habitatrelate_objectid_seq;
ALTER TABLE fieldseeker.habitatrelate ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.inspectionsample DROP CONSTRAINT inspectionsample_pkey;
ALTER TABLE fieldseeker.inspectionsample ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.inspectionsample_objectid_seq;
ALTER TABLE fieldseeker.inspectionsample ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.inspectionsampledetail DROP CONSTRAINT inspectionsampledetail_pkey;
ALTER TABLE fieldseeker.inspectionsampledetail ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.inspectionsampledetail_objectid_seq;
ALTER TABLE fieldseeker.inspectionsampledetail ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.linelocation DROP CONSTRAINT linelocation_pkey;
ALTER TABLE fieldseeker.linelocation ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.linelocation_objectid_seq;
ALTER TABLE fieldseeker.linelocation ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.locationtracking DROP CONSTRAINT locationtracking_pkey;
ALTER TABLE fieldseeker.locationtracking ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.locationtracking_objectid_seq;
ALTER TABLE fieldseeker.locationtracking ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.mosquitoinspection DROP CONSTRAINT mosquitoinspection_pkey;
ALTER TABLE fieldseeker.mosquitoinspection ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.mosquitoinspection_objectid_seq;
ALTER TABLE fieldseeker.mosquitoinspection ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.pointlocation DROP CONSTRAINT pointlocation_pkey;
ALTER TABLE fieldseeker.pointlocation ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.pointlocation_objectid_seq;
ALTER TABLE fieldseeker.pointlocation ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.polygonlocation DROP CONSTRAINT polygonlocation_pkey;
ALTER TABLE fieldseeker.polygonlocation ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.polygonlocation_objectid_seq;
ALTER TABLE fieldseeker.polygonlocation ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.pool DROP CONSTRAINT pool_pkey;
ALTER TABLE fieldseeker.pool ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.pool_objectid_seq;
ALTER TABLE fieldseeker.pool ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.pooldetail DROP CONSTRAINT pooldetail_pkey;
ALTER TABLE fieldseeker.pooldetail ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.pooldetail_objectid_seq;
ALTER TABLE fieldseeker.pooldetail ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.proposedtreatmentarea DROP CONSTRAINT proposedtreatmentarea_pkey;
ALTER TABLE fieldseeker.proposedtreatmentarea ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.proposedtreatmentarea_objectid_seq;
ALTER TABLE fieldseeker.proposedtreatmentarea ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.qamosquitoinspection DROP CONSTRAINT qamosquitoinspection_pkey;
ALTER TABLE fieldseeker.qamosquitoinspection ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.qamosquitoinspection_objectid_seq;
ALTER TABLE fieldseeker.qamosquitoinspection ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.rodentlocation DROP CONSTRAINT rodentlocation_pkey;
ALTER TABLE fieldseeker.rodentlocation ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.rodentlocation_objectid_seq;
ALTER TABLE fieldseeker.rodentlocation ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.samplecollection DROP CONSTRAINT samplecollection_pkey;
ALTER TABLE fieldseeker.samplecollection ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.samplecollection_objectid_seq;
ALTER TABLE fieldseeker.samplecollection ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.samplelocation DROP CONSTRAINT samplelocation_pkey;
ALTER TABLE fieldseeker.samplelocation ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.samplelocation_objectid_seq;
ALTER TABLE fieldseeker.samplelocation ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.servicerequest DROP CONSTRAINT servicerequest_pkey;
ALTER TABLE fieldseeker.servicerequest ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.servicerequest_objectid_seq;
ALTER TABLE fieldseeker.servicerequest ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.speciesabundance DROP CONSTRAINT speciesabundance_pkey;
ALTER TABLE fieldseeker.speciesabundance ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.speciesabundance_objectid_seq;
ALTER TABLE fieldseeker.speciesabundance ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.stormdrain DROP CONSTRAINT stormdrain_pkey;
ALTER TABLE fieldseeker.stormdrain ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.stormdrain_objectid_seq;
ALTER TABLE fieldseeker.stormdrain ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.timecard DROP CONSTRAINT timecard_pkey;
ALTER TABLE fieldseeker.timecard ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.timecard_objectid_seq;
ALTER TABLE fieldseeker.timecard ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.trapdata DROP CONSTRAINT trapdata_pkey;
ALTER TABLE fieldseeker.trapdata ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.trapdata_objectid_seq;
ALTER TABLE fieldseeker.trapdata ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.traplocation DROP CONSTRAINT traplocation_pkey;
ALTER TABLE fieldseeker.traplocation ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.traplocation_objectid_seq;
ALTER TABLE fieldseeker.traplocation ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.treatment DROP CONSTRAINT treatment_pkey;
ALTER TABLE fieldseeker.treatment ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.treatment_objectid_seq;
ALTER TABLE fieldseeker.treatment ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.treatmentarea DROP CONSTRAINT treatmentarea_pkey;
ALTER TABLE fieldseeker.treatmentarea ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.treatmentarea_objectid_seq;
ALTER TABLE fieldseeker.treatmentarea ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.zones DROP CONSTRAINT zones_pkey;
ALTER TABLE fieldseeker.zones ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.zones_objectid_seq;
ALTER TABLE fieldseeker.zones ADD PRIMARY KEY (globalid, version);

ALTER TABLE fieldseeker.zones2 DROP CONSTRAINT zones2_pkey;
ALTER TABLE fieldseeker.zones2 ALTER COLUMN objectid DROP DEFAULT;
DROP SEQUENCE fieldseeker.zones2_objectid_seq;
ALTER TABLE fieldseeker.zones2 ADD PRIMARY KEY (globalid, version);

-- +goose Down
SELECT "not done" FROM non_existent;
