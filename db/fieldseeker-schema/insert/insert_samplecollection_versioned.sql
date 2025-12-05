
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_samplecollection(
	p_objectid bigint,
	
	p_loc_id uuid,
	p_startdatetime timestamp,
	p_enddatetime timestamp,
	p_sitecond varchar,
	p_sampleid varchar,
	p_survtech varchar,
	p_datesent timestamp,
	p_datetested timestamp,
	p_testtech varchar,
	p_comments varchar,
	p_processed smallint,
	p_sampletype varchar,
	p_samplecond varchar,
	p_species varchar,
	p_sex varchar,
	p_avetemp double precision,
	p_windspeed double precision,
	p_winddir varchar,
	p_raingauge double precision,
	p_activity varchar,
	p_testmethod varchar,
	p_diseasetested varchar,
	p_diseasepos varchar,
	p_reviewed smallint,
	p_reviewedby varchar,
	p_revieweddate timestamp,
	p_locationname varchar,
	p_zone varchar,
	p_recordstatus smallint,
	p_zone2 varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_lab varchar,
	p_fieldtech varchar,
	p_flockid uuid,
	p_samplecount smallint,
	p_chickenid uuid,
	p_gatewaysync smallint,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.samplecollection lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.loc_id IS NOT DISTINCT FROM p_loc_id 
		AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime 
		AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime 
		AND lv.sitecond IS NOT DISTINCT FROM p_sitecond 
		AND lv.sampleid IS NOT DISTINCT FROM p_sampleid 
		AND lv.survtech IS NOT DISTINCT FROM p_survtech 
		AND lv.datesent IS NOT DISTINCT FROM p_datesent 
		AND lv.datetested IS NOT DISTINCT FROM p_datetested 
		AND lv.testtech IS NOT DISTINCT FROM p_testtech 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.processed IS NOT DISTINCT FROM p_processed 
		AND lv.sampletype IS NOT DISTINCT FROM p_sampletype 
		AND lv.samplecond IS NOT DISTINCT FROM p_samplecond 
		AND lv.species IS NOT DISTINCT FROM p_species 
		AND lv.sex IS NOT DISTINCT FROM p_sex 
		AND lv.avetemp IS NOT DISTINCT FROM p_avetemp 
		AND lv.windspeed IS NOT DISTINCT FROM p_windspeed 
		AND lv.winddir IS NOT DISTINCT FROM p_winddir 
		AND lv.raingauge IS NOT DISTINCT FROM p_raingauge 
		AND lv.activity IS NOT DISTINCT FROM p_activity 
		AND lv.testmethod IS NOT DISTINCT FROM p_testmethod 
		AND lv.diseasetested IS NOT DISTINCT FROM p_diseasetested 
		AND lv.diseasepos IS NOT DISTINCT FROM p_diseasepos 
		AND lv.reviewed IS NOT DISTINCT FROM p_reviewed 
		AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby 
		AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate 
		AND lv.locationname IS NOT DISTINCT FROM p_locationname 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.lab IS NOT DISTINCT FROM p_lab 
		AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech 
		AND lv.flockid IS NOT DISTINCT FROM p_flockid 
		AND lv.samplecount IS NOT DISTINCT FROM p_samplecount 
		AND lv.chickenid IS NOT DISTINCT FROM p_chickenid 
		AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.samplecollection 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.samplecollection
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.samplecollection (
		objectid,
		
		loc_id,
		startdatetime,
		enddatetime,
		sitecond,
		sampleid,
		survtech,
		datesent,
		datetested,
		testtech,
		comments,
		processed,
		sampletype,
		samplecond,
		species,
		sex,
		avetemp,
		windspeed,
		winddir,
		raingauge,
		activity,
		testmethod,
		diseasetested,
		diseasepos,
		reviewed,
		reviewedby,
		revieweddate,
		locationname,
		zone,
		recordstatus,
		zone2,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		lab,
		fieldtech,
		flockid,
		samplecount,
		chickenid,
		gatewaysync,
		creationdate,
		creator,
		editdate,
		editor,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_loc_id,
		p_startdatetime,
		p_enddatetime,
		p_sitecond,
		p_sampleid,
		p_survtech,
		p_datesent,
		p_datetested,
		p_testtech,
		p_comments,
		p_processed,
		p_sampletype,
		p_samplecond,
		p_species,
		p_sex,
		p_avetemp,
		p_windspeed,
		p_winddir,
		p_raingauge,
		p_activity,
		p_testmethod,
		p_diseasetested,
		p_diseasepos,
		p_reviewed,
		p_reviewedby,
		p_revieweddate,
		p_locationname,
		p_zone,
		p_recordstatus,
		p_zone2,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_lab,
		p_fieldtech,
		p_flockid,
		p_samplecount,
		p_chickenid,
		p_gatewaysync,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
