
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_rodentlocation(
	p_objectid bigint,
	
	p_locationname varchar,
	p_zone varchar,
	p_zone2 varchar,
	p_habitat varchar,
	p_priority varchar,
	p_usetype varchar,
	p_active smallint,
	p_description varchar,
	p_accessdesc varchar,
	p_comments varchar,
	p_symbology varchar,
	p_externalid varchar,
	p_nextactiondatescheduled timestamp,
	p_locationnumber integer,
	p_lastinspectdate timestamp,
	p_lastinspectspecies varchar,
	p_lastinspectaction varchar,
	p_lastinspectconditions varchar,
	p_lastinspectrodentevidence varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_jurisdiction varchar,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.rodentlocation lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.locationname IS NOT DISTINCT FROM p_locationname 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.habitat IS NOT DISTINCT FROM p_habitat 
		AND lv.priority IS NOT DISTINCT FROM p_priority 
		AND lv.usetype IS NOT DISTINCT FROM p_usetype 
		AND lv.active IS NOT DISTINCT FROM p_active 
		AND lv.description IS NOT DISTINCT FROM p_description 
		AND lv.accessdesc IS NOT DISTINCT FROM p_accessdesc 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.symbology IS NOT DISTINCT FROM p_symbology 
		AND lv.externalid IS NOT DISTINCT FROM p_externalid 
		AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled 
		AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber 
		AND lv.lastinspectdate IS NOT DISTINCT FROM p_lastinspectdate 
		AND lv.lastinspectspecies IS NOT DISTINCT FROM p_lastinspectspecies 
		AND lv.lastinspectaction IS NOT DISTINCT FROM p_lastinspectaction 
		AND lv.lastinspectconditions IS NOT DISTINCT FROM p_lastinspectconditions 
		AND lv.lastinspectrodentevidence IS NOT DISTINCT FROM p_lastinspectrodentevidence 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.rodentlocation 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.rodentlocation
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.rodentlocation (
		objectid,
		
		locationname,
		zone,
		zone2,
		habitat,
		priority,
		usetype,
		active,
		description,
		accessdesc,
		comments,
		symbology,
		externalid,
		nextactiondatescheduled,
		locationnumber,
		lastinspectdate,
		lastinspectspecies,
		lastinspectaction,
		lastinspectconditions,
		lastinspectrodentevidence,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		creationdate,
		creator,
		editdate,
		editor,
		jurisdiction,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_locationname,
		p_zone,
		p_zone2,
		p_habitat,
		p_priority,
		p_usetype,
		p_active,
		p_description,
		p_accessdesc,
		p_comments,
		p_symbology,
		p_externalid,
		p_nextactiondatescheduled,
		p_locationnumber,
		p_lastinspectdate,
		p_lastinspectspecies,
		p_lastinspectaction,
		p_lastinspectconditions,
		p_lastinspectrodentevidence,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_jurisdiction,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
