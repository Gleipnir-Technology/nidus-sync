
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_pool(
	p_objectid bigint,
	
	p_trapdata_id uuid,
	p_datesent timestamp,
	p_survtech varchar,
	p_datetested timestamp,
	p_testtech varchar,
	p_comments varchar,
	p_sampleid varchar,
	p_processed smallint,
	p_lab_id uuid,
	p_testmethod varchar,
	p_diseasetested varchar,
	p_diseasepos varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_lab varchar,
	p_poolyear smallint,
	p_gatewaysync smallint,
	p_vectorsurvcollectionid varchar,
	p_vectorsurvpoolid varchar,
	p_vectorsurvtrapdataid varchar,
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
		SELECT 1 FROM fieldseeker.pool lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.trapdata_id IS NOT DISTINCT FROM p_trapdata_id 
		AND lv.datesent IS NOT DISTINCT FROM p_datesent 
		AND lv.survtech IS NOT DISTINCT FROM p_survtech 
		AND lv.datetested IS NOT DISTINCT FROM p_datetested 
		AND lv.testtech IS NOT DISTINCT FROM p_testtech 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.sampleid IS NOT DISTINCT FROM p_sampleid 
		AND lv.processed IS NOT DISTINCT FROM p_processed 
		AND lv.lab_id IS NOT DISTINCT FROM p_lab_id 
		AND lv.testmethod IS NOT DISTINCT FROM p_testmethod 
		AND lv.diseasetested IS NOT DISTINCT FROM p_diseasetested 
		AND lv.diseasepos IS NOT DISTINCT FROM p_diseasepos 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.lab IS NOT DISTINCT FROM p_lab 
		AND lv.poolyear IS NOT DISTINCT FROM p_poolyear 
		AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync 
		AND lv.vectorsurvcollectionid IS NOT DISTINCT FROM p_vectorsurvcollectionid 
		AND lv.vectorsurvpoolid IS NOT DISTINCT FROM p_vectorsurvpoolid 
		AND lv.vectorsurvtrapdataid IS NOT DISTINCT FROM p_vectorsurvtrapdataid 
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
				(SELECT VERSION FROM fieldseeker.pool 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.pool
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.pool (
		objectid,
		
		trapdata_id,
		datesent,
		survtech,
		datetested,
		testtech,
		comments,
		sampleid,
		processed,
		lab_id,
		testmethod,
		diseasetested,
		diseasepos,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		lab,
		poolyear,
		gatewaysync,
		vectorsurvcollectionid,
		vectorsurvpoolid,
		vectorsurvtrapdataid,
		creationdate,
		creator,
		editdate,
		editor,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_trapdata_id,
		p_datesent,
		p_survtech,
		p_datetested,
		p_testtech,
		p_comments,
		p_sampleid,
		p_processed,
		p_lab_id,
		p_testmethod,
		p_diseasetested,
		p_diseasepos,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_lab,
		p_poolyear,
		p_gatewaysync,
		p_vectorsurvcollectionid,
		p_vectorsurvpoolid,
		p_vectorsurvtrapdataid,
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
