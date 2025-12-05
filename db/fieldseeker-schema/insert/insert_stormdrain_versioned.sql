
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_stormdrain(
	p_objectid bigint,
	
	p_nexttreatmentdate timestamp,
	p_lasttreatdate timestamp,
	p_lastaction varchar,
	p_symbology varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_laststatus varchar,
	p_zone varchar,
	p_zone2 varchar,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_type varchar,
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
		SELECT 1 FROM fieldseeker.stormdrain lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.nexttreatmentdate IS NOT DISTINCT FROM p_nexttreatmentdate 
		AND lv.lasttreatdate IS NOT DISTINCT FROM p_lasttreatdate 
		AND lv.lastaction IS NOT DISTINCT FROM p_lastaction 
		AND lv.symbology IS NOT DISTINCT FROM p_symbology 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.laststatus IS NOT DISTINCT FROM p_laststatus 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.type IS NOT DISTINCT FROM p_type 
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
				(SELECT VERSION FROM fieldseeker.stormdrain 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.stormdrain
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.stormdrain (
		objectid,
		
		nexttreatmentdate,
		lasttreatdate,
		lastaction,
		symbology,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		laststatus,
		zone,
		zone2,
		creationdate,
		creator,
		editdate,
		editor,
		type,
		jurisdiction,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_nexttreatmentdate,
		p_lasttreatdate,
		p_lastaction,
		p_symbology,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_laststatus,
		p_zone,
		p_zone2,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_type,
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
