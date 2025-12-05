
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_linelocation(
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
	p_symbology varchar,
	p_externalid varchar,
	p_acres double precision,
	p_nextactiondatescheduled timestamp,
	p_larvinspectinterval smallint,
	p_length_ft double precision,
	p_width_ft double precision,
	p_zone2 varchar,
	p_locationnumber integer,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_lastinspectdate timestamp,
	p_lastinspectbreeding varchar,
	p_lastinspectavglarvae double precision,
	p_lastinspectavgpupae double precision,
	p_lastinspectlstages varchar,
	p_lastinspectactiontaken varchar,
	p_lastinspectfieldspecies varchar,
	p_lasttreatdate timestamp,
	p_lasttreatproduct varchar,
	p_lasttreatqty double precision,
	p_lasttreatqtyunit varchar,
	p_hectares double precision,
	p_lastinspectactivity varchar,
	p_lasttreatactivity varchar,
	p_length_meters double precision,
	p_width_meters double precision,
	p_lastinspectconditions varchar,
	p_waterorigin varchar,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_jurisdiction varchar,
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
		SELECT 1 FROM fieldseeker.linelocation lv 
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
		AND lv.symbology IS NOT DISTINCT FROM p_symbology 
		AND lv.externalid IS NOT DISTINCT FROM p_externalid 
		AND lv.acres IS NOT DISTINCT FROM p_acres 
		AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled 
		AND lv.larvinspectinterval IS NOT DISTINCT FROM p_larvinspectinterval 
		AND lv.length_ft IS NOT DISTINCT FROM p_length_ft 
		AND lv.width_ft IS NOT DISTINCT FROM p_width_ft 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.lastinspectdate IS NOT DISTINCT FROM p_lastinspectdate 
		AND lv.lastinspectbreeding IS NOT DISTINCT FROM p_lastinspectbreeding 
		AND lv.lastinspectavglarvae IS NOT DISTINCT FROM p_lastinspectavglarvae 
		AND lv.lastinspectavgpupae IS NOT DISTINCT FROM p_lastinspectavgpupae 
		AND lv.lastinspectlstages IS NOT DISTINCT FROM p_lastinspectlstages 
		AND lv.lastinspectactiontaken IS NOT DISTINCT FROM p_lastinspectactiontaken 
		AND lv.lastinspectfieldspecies IS NOT DISTINCT FROM p_lastinspectfieldspecies 
		AND lv.lasttreatdate IS NOT DISTINCT FROM p_lasttreatdate 
		AND lv.lasttreatproduct IS NOT DISTINCT FROM p_lasttreatproduct 
		AND lv.lasttreatqty IS NOT DISTINCT FROM p_lasttreatqty 
		AND lv.lasttreatqtyunit IS NOT DISTINCT FROM p_lasttreatqtyunit 
		AND lv.hectares IS NOT DISTINCT FROM p_hectares 
		AND lv.lastinspectactivity IS NOT DISTINCT FROM p_lastinspectactivity 
		AND lv.lasttreatactivity IS NOT DISTINCT FROM p_lasttreatactivity 
		AND lv.length_meters IS NOT DISTINCT FROM p_length_meters 
		AND lv.width_meters IS NOT DISTINCT FROM p_width_meters 
		AND lv.lastinspectconditions IS NOT DISTINCT FROM p_lastinspectconditions 
		AND lv.waterorigin IS NOT DISTINCT FROM p_waterorigin 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction 
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
				(SELECT VERSION FROM fieldseeker.linelocation 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.linelocation
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.linelocation (
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
		symbology,
		externalid,
		acres,
		nextactiondatescheduled,
		larvinspectinterval,
		length_ft,
		width_ft,
		zone2,
		locationnumber,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		lastinspectdate,
		lastinspectbreeding,
		lastinspectavglarvae,
		lastinspectavgpupae,
		lastinspectlstages,
		lastinspectactiontaken,
		lastinspectfieldspecies,
		lasttreatdate,
		lasttreatproduct,
		lasttreatqty,
		lasttreatqtyunit,
		hectares,
		lastinspectactivity,
		lasttreatactivity,
		length_meters,
		width_meters,
		lastinspectconditions,
		waterorigin,
		creationdate,
		creator,
		editdate,
		editor,
		jurisdiction,
		shape__length,
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
		p_symbology,
		p_externalid,
		p_acres,
		p_nextactiondatescheduled,
		p_larvinspectinterval,
		p_length_ft,
		p_width_ft,
		p_zone2,
		p_locationnumber,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_lastinspectdate,
		p_lastinspectbreeding,
		p_lastinspectavglarvae,
		p_lastinspectavgpupae,
		p_lastinspectlstages,
		p_lastinspectactiontaken,
		p_lastinspectfieldspecies,
		p_lasttreatdate,
		p_lasttreatproduct,
		p_lasttreatqty,
		p_lasttreatqtyunit,
		p_hectares,
		p_lastinspectactivity,
		p_lasttreatactivity,
		p_length_meters,
		p_width_meters,
		p_lastinspectconditions,
		p_waterorigin,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_jurisdiction,
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
