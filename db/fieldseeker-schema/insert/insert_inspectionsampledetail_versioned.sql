
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_inspectionsampledetail(
	p_objectid bigint,
	
	p_inspsample_id uuid,
	p_fieldspecies varchar,
	p_flarvcount smallint,
	p_fpupcount smallint,
	p_feggcount smallint,
	p_flstages varchar,
	p_fdomstage varchar,
	p_fadultact varchar,
	p_labspecies varchar,
	p_llarvcount smallint,
	p_lpupcount smallint,
	p_leggcount smallint,
	p_ldomstage varchar,
	p_comments varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_processed smallint,
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
		SELECT 1 FROM fieldseeker.inspectionsampledetail lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.inspsample_id IS NOT DISTINCT FROM p_inspsample_id 
		AND lv.fieldspecies IS NOT DISTINCT FROM p_fieldspecies 
		AND lv.flarvcount IS NOT DISTINCT FROM p_flarvcount 
		AND lv.fpupcount IS NOT DISTINCT FROM p_fpupcount 
		AND lv.feggcount IS NOT DISTINCT FROM p_feggcount 
		AND lv.flstages IS NOT DISTINCT FROM p_flstages 
		AND lv.fdomstage IS NOT DISTINCT FROM p_fdomstage 
		AND lv.fadultact IS NOT DISTINCT FROM p_fadultact 
		AND lv.labspecies IS NOT DISTINCT FROM p_labspecies 
		AND lv.llarvcount IS NOT DISTINCT FROM p_llarvcount 
		AND lv.lpupcount IS NOT DISTINCT FROM p_lpupcount 
		AND lv.leggcount IS NOT DISTINCT FROM p_leggcount 
		AND lv.ldomstage IS NOT DISTINCT FROM p_ldomstage 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.processed IS NOT DISTINCT FROM p_processed 
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
				(SELECT VERSION FROM fieldseeker.inspectionsampledetail 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.inspectionsampledetail
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.inspectionsampledetail (
		objectid,
		
		inspsample_id,
		fieldspecies,
		flarvcount,
		fpupcount,
		feggcount,
		flstages,
		fdomstage,
		fadultact,
		labspecies,
		llarvcount,
		lpupcount,
		leggcount,
		ldomstage,
		comments,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		processed,
		creationdate,
		creator,
		editdate,
		editor,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_inspsample_id,
		p_fieldspecies,
		p_flarvcount,
		p_fpupcount,
		p_feggcount,
		p_flstages,
		p_fdomstage,
		p_fadultact,
		p_labspecies,
		p_llarvcount,
		p_lpupcount,
		p_leggcount,
		p_ldomstage,
		p_comments,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_processed,
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
