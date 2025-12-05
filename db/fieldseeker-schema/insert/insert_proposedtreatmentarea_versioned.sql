
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_proposedtreatmentarea(
	p_objectid bigint,
	
	p_method varchar,
	p_comments varchar,
	p_zone varchar,
	p_reviewed smallint,
	p_reviewedby varchar,
	p_revieweddate timestamp,
	p_zone2 varchar,
	p_completeddate timestamp,
	p_completedby varchar,
	p_completed smallint,
	p_issprayroute smallint,
	p_name varchar,
	p_acres double precision,
	p_globalid uuid,
	p_exported smallint,
	p_targetproduct varchar,
	p_targetapprate double precision,
	p_hectares double precision,
	p_lasttreatactivity varchar,
	p_lasttreatdate timestamp,
	p_lasttreatproduct varchar,
	p_lasttreatqty double precision,
	p_lasttreatqtyunit varchar,
	p_priority varchar,
	p_duedate timestamp,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_targetspecies varchar,
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
		SELECT 1 FROM fieldseeker.proposedtreatmentarea lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.method IS NOT DISTINCT FROM p_method 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.reviewed IS NOT DISTINCT FROM p_reviewed 
		AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby 
		AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.completeddate IS NOT DISTINCT FROM p_completeddate 
		AND lv.completedby IS NOT DISTINCT FROM p_completedby 
		AND lv.completed IS NOT DISTINCT FROM p_completed 
		AND lv.issprayroute IS NOT DISTINCT FROM p_issprayroute 
		AND lv.name IS NOT DISTINCT FROM p_name 
		AND lv.acres IS NOT DISTINCT FROM p_acres 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.exported IS NOT DISTINCT FROM p_exported 
		AND lv.targetproduct IS NOT DISTINCT FROM p_targetproduct 
		AND lv.targetapprate IS NOT DISTINCT FROM p_targetapprate 
		AND lv.hectares IS NOT DISTINCT FROM p_hectares 
		AND lv.lasttreatactivity IS NOT DISTINCT FROM p_lasttreatactivity 
		AND lv.lasttreatdate IS NOT DISTINCT FROM p_lasttreatdate 
		AND lv.lasttreatproduct IS NOT DISTINCT FROM p_lasttreatproduct 
		AND lv.lasttreatqty IS NOT DISTINCT FROM p_lasttreatqty 
		AND lv.lasttreatqtyunit IS NOT DISTINCT FROM p_lasttreatqtyunit 
		AND lv.priority IS NOT DISTINCT FROM p_priority 
		AND lv.duedate IS NOT DISTINCT FROM p_duedate 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.targetspecies IS NOT DISTINCT FROM p_targetspecies 
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
				(SELECT VERSION FROM fieldseeker.proposedtreatmentarea 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.proposedtreatmentarea
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.proposedtreatmentarea (
		objectid,
		
		method,
		comments,
		zone,
		reviewed,
		reviewedby,
		revieweddate,
		zone2,
		completeddate,
		completedby,
		completed,
		issprayroute,
		name,
		acres,
		globalid,
		exported,
		targetproduct,
		targetapprate,
		hectares,
		lasttreatactivity,
		lasttreatdate,
		lasttreatproduct,
		lasttreatqty,
		lasttreatqtyunit,
		priority,
		duedate,
		creationdate,
		creator,
		editdate,
		editor,
		targetspecies,
		shape__area,
		shape__length,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_method,
		p_comments,
		p_zone,
		p_reviewed,
		p_reviewedby,
		p_revieweddate,
		p_zone2,
		p_completeddate,
		p_completedby,
		p_completed,
		p_issprayroute,
		p_name,
		p_acres,
		p_globalid,
		p_exported,
		p_targetproduct,
		p_targetapprate,
		p_hectares,
		p_lasttreatactivity,
		p_lasttreatdate,
		p_lasttreatproduct,
		p_lasttreatqty,
		p_lasttreatqtyunit,
		p_priority,
		p_duedate,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_targetspecies,
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
