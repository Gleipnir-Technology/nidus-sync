
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_traplocation(
	p_objectid bigint,
	
	p_name varchar,
	p_zone varchar,
	p_habitat varchar,
	p_priority varchar,
	p_usetype varchar,
	p_active smallint,
	p_description varchar,
	p_accessdesc varchar,
	p_comments varchar,
	p_externalid varchar,
	p_nextactiondatescheduled timestamp,
	p_zone2 varchar,
	p_locationnumber integer,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_gatewaysync smallint,
	p_route integer,
	p_set_dow integer,
	p_route_order integer,
	p_vectorsurvsiteid varchar,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
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
		SELECT 1 FROM fieldseeker.traplocation lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.name IS NOT DISTINCT FROM p_name 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.habitat IS NOT DISTINCT FROM p_habitat 
		AND lv.priority IS NOT DISTINCT FROM p_priority 
		AND lv.usetype IS NOT DISTINCT FROM p_usetype 
		AND lv.active IS NOT DISTINCT FROM p_active 
		AND lv.description IS NOT DISTINCT FROM p_description 
		AND lv.accessdesc IS NOT DISTINCT FROM p_accessdesc 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.externalid IS NOT DISTINCT FROM p_externalid 
		AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync 
		AND lv.route IS NOT DISTINCT FROM p_route 
		AND lv.set_dow IS NOT DISTINCT FROM p_set_dow 
		AND lv.route_order IS NOT DISTINCT FROM p_route_order 
		AND lv.vectorsurvsiteid IS NOT DISTINCT FROM p_vectorsurvsiteid 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
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
				(SELECT VERSION FROM fieldseeker.traplocation 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.traplocation
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.traplocation (
		objectid,
		
		name,
		zone,
		habitat,
		priority,
		usetype,
		active,
		description,
		accessdesc,
		comments,
		externalid,
		nextactiondatescheduled,
		zone2,
		locationnumber,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		gatewaysync,
		route,
		set_dow,
		route_order,
		vectorsurvsiteid,
		creationdate,
		creator,
		editdate,
		editor,
		h3r7,
		h3r8,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_name,
		p_zone,
		p_habitat,
		p_priority,
		p_usetype,
		p_active,
		p_description,
		p_accessdesc,
		p_comments,
		p_externalid,
		p_nextactiondatescheduled,
		p_zone2,
		p_locationnumber,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_gatewaysync,
		p_route,
		p_set_dow,
		p_route_order,
		p_vectorsurvsiteid,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
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
