
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_treatmentarea(
	p_objectid bigint,
	
	p_treat_id uuid,
	p_session_id uuid,
	p_treatdate timestamp,
	p_comments varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_notified smallint,
	p_type varchar,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_shape__area double precision,
	p_shape__length double precision,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.treatmentarea lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.treat_id IS NOT DISTINCT FROM p_treat_id 
		AND lv.session_id IS NOT DISTINCT FROM p_session_id 
		AND lv.treatdate IS NOT DISTINCT FROM p_treatdate 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.notified IS NOT DISTINCT FROM p_notified 
		AND lv.type IS NOT DISTINCT FROM p_type 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.shape__area IS NOT DISTINCT FROM p_shape__area 
		AND lv.shape__length IS NOT DISTINCT FROM p_shape__length 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.treatmentarea 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.treatmentarea
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.treatmentarea (
		objectid,
		
		treat_id,
		session_id,
		treatdate,
		comments,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		notified,
		type,
		creationdate,
		creator,
		editdate,
		editor,
		shape__area,
		shape__length,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_treat_id,
		p_session_id,
		p_treatdate,
		p_comments,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_notified,
		p_type,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_shape__area,
		p_shape__length,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
