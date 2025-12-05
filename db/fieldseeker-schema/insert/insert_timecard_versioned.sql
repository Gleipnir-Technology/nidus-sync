
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_timecard(
	p_objectid bigint,
	
	p_activity varchar,
	p_startdatetime timestamp,
	p_enddatetime timestamp,
	p_comments varchar,
	p_externalid varchar,
	p_equiptype varchar,
	p_locationname varchar,
	p_zone varchar,
	p_zone2 varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_linelocid uuid,
	p_pointlocid uuid,
	p_polygonlocid uuid,
	p_lclocid uuid,
	p_samplelocid uuid,
	p_srid uuid,
	p_traplocid uuid,
	p_fieldtech varchar,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_rodentlocid uuid,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.timecard lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.activity IS NOT DISTINCT FROM p_activity 
		AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime 
		AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.externalid IS NOT DISTINCT FROM p_externalid 
		AND lv.equiptype IS NOT DISTINCT FROM p_equiptype 
		AND lv.locationname IS NOT DISTINCT FROM p_locationname 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.linelocid IS NOT DISTINCT FROM p_linelocid 
		AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid 
		AND lv.polygonlocid IS NOT DISTINCT FROM p_polygonlocid 
		AND lv.lclocid IS NOT DISTINCT FROM p_lclocid 
		AND lv.samplelocid IS NOT DISTINCT FROM p_samplelocid 
		AND lv.srid IS NOT DISTINCT FROM p_srid 
		AND lv.traplocid IS NOT DISTINCT FROM p_traplocid 
		AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.rodentlocid IS NOT DISTINCT FROM p_rodentlocid 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.timecard 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.timecard
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.timecard (
		objectid,
		
		activity,
		startdatetime,
		enddatetime,
		comments,
		externalid,
		equiptype,
		locationname,
		zone,
		zone2,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		linelocid,
		pointlocid,
		polygonlocid,
		lclocid,
		samplelocid,
		srid,
		traplocid,
		fieldtech,
		creationdate,
		creator,
		editdate,
		editor,
		rodentlocid,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_activity,
		p_startdatetime,
		p_enddatetime,
		p_comments,
		p_externalid,
		p_equiptype,
		p_locationname,
		p_zone,
		p_zone2,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_linelocid,
		p_pointlocid,
		p_polygonlocid,
		p_lclocid,
		p_samplelocid,
		p_srid,
		p_traplocid,
		p_fieldtech,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_rodentlocid,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
