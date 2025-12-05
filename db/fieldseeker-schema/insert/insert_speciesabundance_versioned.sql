
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_speciesabundance(
	p_objectid bigint,
	
	p_trapdata_id uuid,
	p_species varchar,
	p_males smallint,
	p_unknown smallint,
	p_bloodedfem smallint,
	p_gravidfem smallint,
	p_larvae smallint,
	p_poolstogen smallint,
	p_processed smallint,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_pupae smallint,
	p_eggs smallint,
	p_females integer,
	p_total integer,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_yearweek integer,
	p_globalzscore double precision,
	p_r7score double precision,
	p_r8score double precision,
	p_h3r7 varchar,
	p_h3r8 varchar,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.speciesabundance lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.trapdata_id IS NOT DISTINCT FROM p_trapdata_id 
		AND lv.species IS NOT DISTINCT FROM p_species 
		AND lv.males IS NOT DISTINCT FROM p_males 
		AND lv.unknown IS NOT DISTINCT FROM p_unknown 
		AND lv.bloodedfem IS NOT DISTINCT FROM p_bloodedfem 
		AND lv.gravidfem IS NOT DISTINCT FROM p_gravidfem 
		AND lv.larvae IS NOT DISTINCT FROM p_larvae 
		AND lv.poolstogen IS NOT DISTINCT FROM p_poolstogen 
		AND lv.processed IS NOT DISTINCT FROM p_processed 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.pupae IS NOT DISTINCT FROM p_pupae 
		AND lv.eggs IS NOT DISTINCT FROM p_eggs 
		AND lv.females IS NOT DISTINCT FROM p_females 
		AND lv.total IS NOT DISTINCT FROM p_total 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.yearweek IS NOT DISTINCT FROM p_yearweek 
		AND lv.globalzscore IS NOT DISTINCT FROM p_globalzscore 
		AND lv.r7score IS NOT DISTINCT FROM p_r7score 
		AND lv.r8score IS NOT DISTINCT FROM p_r8score 
		AND lv.h3r7 IS NOT DISTINCT FROM p_h3r7 
		AND lv.h3r8 IS NOT DISTINCT FROM p_h3r8 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.speciesabundance 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.speciesabundance
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.speciesabundance (
		objectid,
		
		trapdata_id,
		species,
		males,
		unknown,
		bloodedfem,
		gravidfem,
		larvae,
		poolstogen,
		processed,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		pupae,
		eggs,
		females,
		total,
		creationdate,
		creator,
		editdate,
		editor,
		yearweek,
		globalzscore,
		r7score,
		r8score,
		h3r7,
		h3r8,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_trapdata_id,
		p_species,
		p_males,
		p_unknown,
		p_bloodedfem,
		p_gravidfem,
		p_larvae,
		p_poolstogen,
		p_processed,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_pupae,
		p_eggs,
		p_females,
		p_total,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_yearweek,
		p_globalzscore,
		p_r7score,
		p_r8score,
		p_h3r7,
		p_h3r8,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
