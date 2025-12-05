
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_fieldscoutinglog(
	p_objectid bigint,
	
	p_status smallint,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
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
		SELECT 1 FROM fieldseeker.fieldscoutinglog lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.status IS NOT DISTINCT FROM p_status 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
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
				(SELECT VERSION FROM fieldseeker.fieldscoutinglog 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.fieldscoutinglog
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.fieldscoutinglog (
		objectid,
		
		status,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		creationdate,
		creator,
		editdate,
		editor,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_status,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
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
