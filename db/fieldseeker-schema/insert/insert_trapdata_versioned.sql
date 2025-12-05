
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_trapdata(
	p_objectid bigint,
	
	p_traptype varchar,
	p_trapactivitytype varchar,
	p_startdatetime timestamp,
	p_enddatetime timestamp,
	p_comments varchar,
	p_idbytech varchar,
	p_sortbytech varchar,
	p_processed smallint,
	p_sitecond varchar,
	p_locationname varchar,
	p_recordstatus smallint,
	p_reviewed smallint,
	p_reviewedby varchar,
	p_revieweddate timestamp,
	p_trapcondition varchar,
	p_trapnights smallint,
	p_zone varchar,
	p_zone2 varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_srid uuid,
	p_fieldtech varchar,
	p_gatewaysync smallint,
	p_loc_id uuid,
	p_voltage double precision,
	p_winddir varchar,
	p_windspeed double precision,
	p_avetemp double precision,
	p_raingauge double precision,
	p_lr smallint,
	p_field integer,
	p_vectorsurvtrapdataid varchar,
	p_vectorsurvtraplocationid varchar,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_lure varchar,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.trapdata lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.traptype IS NOT DISTINCT FROM p_traptype 
		AND lv.trapactivitytype IS NOT DISTINCT FROM p_trapactivitytype 
		AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime 
		AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.idbytech IS NOT DISTINCT FROM p_idbytech 
		AND lv.sortbytech IS NOT DISTINCT FROM p_sortbytech 
		AND lv.processed IS NOT DISTINCT FROM p_processed 
		AND lv.sitecond IS NOT DISTINCT FROM p_sitecond 
		AND lv.locationname IS NOT DISTINCT FROM p_locationname 
		AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus 
		AND lv.reviewed IS NOT DISTINCT FROM p_reviewed 
		AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby 
		AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate 
		AND lv.trapcondition IS NOT DISTINCT FROM p_trapcondition 
		AND lv.trapnights IS NOT DISTINCT FROM p_trapnights 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.srid IS NOT DISTINCT FROM p_srid 
		AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech 
		AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync 
		AND lv.loc_id IS NOT DISTINCT FROM p_loc_id 
		AND lv.voltage IS NOT DISTINCT FROM p_voltage 
		AND lv.winddir IS NOT DISTINCT FROM p_winddir 
		AND lv.windspeed IS NOT DISTINCT FROM p_windspeed 
		AND lv.avetemp IS NOT DISTINCT FROM p_avetemp 
		AND lv.raingauge IS NOT DISTINCT FROM p_raingauge 
		AND lv.lr IS NOT DISTINCT FROM p_lr 
		AND lv.field IS NOT DISTINCT FROM p_field 
		AND lv.vectorsurvtrapdataid IS NOT DISTINCT FROM p_vectorsurvtrapdataid 
		AND lv.vectorsurvtraplocationid IS NOT DISTINCT FROM p_vectorsurvtraplocationid 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.lure IS NOT DISTINCT FROM p_lure 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.trapdata 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.trapdata
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.trapdata (
		objectid,
		
		traptype,
		trapactivitytype,
		startdatetime,
		enddatetime,
		comments,
		idbytech,
		sortbytech,
		processed,
		sitecond,
		locationname,
		recordstatus,
		reviewed,
		reviewedby,
		revieweddate,
		trapcondition,
		trapnights,
		zone,
		zone2,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		srid,
		fieldtech,
		gatewaysync,
		loc_id,
		voltage,
		winddir,
		windspeed,
		avetemp,
		raingauge,
		lr,
		field,
		vectorsurvtrapdataid,
		vectorsurvtraplocationid,
		creationdate,
		creator,
		editdate,
		editor,
		lure,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_traptype,
		p_trapactivitytype,
		p_startdatetime,
		p_enddatetime,
		p_comments,
		p_idbytech,
		p_sortbytech,
		p_processed,
		p_sitecond,
		p_locationname,
		p_recordstatus,
		p_reviewed,
		p_reviewedby,
		p_revieweddate,
		p_trapcondition,
		p_trapnights,
		p_zone,
		p_zone2,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_srid,
		p_fieldtech,
		p_gatewaysync,
		p_loc_id,
		p_voltage,
		p_winddir,
		p_windspeed,
		p_avetemp,
		p_raingauge,
		p_lr,
		p_field,
		p_vectorsurvtrapdataid,
		p_vectorsurvtraplocationid,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_lure,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
