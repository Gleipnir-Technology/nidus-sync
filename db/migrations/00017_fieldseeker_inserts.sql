-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_containerrelate(
	p_objectid bigint,
	p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_inspsampleid uuid,p_mosquitoinspid uuid,p_treatmentid uuid,p_containertype varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.containerrelate lv 
		WHERE lv.objectid = p_objectid
		AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.inspsampleid IS NOT DISTINCT FROM p_inspsampleid AND lv.mosquitoinspid IS NOT DISTINCT FROM p_mosquitoinspid AND lv.treatmentid IS NOT DISTINCT FROM p_treatmentid AND lv.containertype IS NOT DISTINCT FROM p_containertype AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.containerrelate 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.containerrelate
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.containerrelate (
		objectid,
		globalid, created_user, created_date, last_edited_user, last_edited_date, inspsampleid, mosquitoinspid, treatmentid, containertype, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_inspsampleid, p_mosquitoinspid, p_treatmentid, p_containertype, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_fieldscoutinglog(
	p_objectid bigint,
	p_status smallint,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.fieldscoutinglog lv 
		WHERE lv.objectid = p_objectid
		AND lv.status IS NOT DISTINCT FROM p_status AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
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
		status, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_status, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_habitatrelate(
	p_objectid bigint,
	p_foreign_id uuid,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_habitattype varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.habitatrelate lv 
		WHERE lv.objectid = p_objectid
		AND lv.foreign_id IS NOT DISTINCT FROM p_foreign_id AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.habitattype IS NOT DISTINCT FROM p_habitattype AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.habitatrelate 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.habitatrelate
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.habitatrelate (
		objectid,
		foreign_id, globalid, created_user, created_date, last_edited_user, last_edited_date, habitattype, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_foreign_id, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_habitattype, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_inspectionsample(
	p_objectid bigint,
	p_insp_id uuid,p_sampleid varchar,p_processed smallint,p_idbytech varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.inspectionsample lv 
		WHERE lv.objectid = p_objectid
		AND lv.insp_id IS NOT DISTINCT FROM p_insp_id AND lv.sampleid IS NOT DISTINCT FROM p_sampleid AND lv.processed IS NOT DISTINCT FROM p_processed AND lv.idbytech IS NOT DISTINCT FROM p_idbytech AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.inspectionsample 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.inspectionsample
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.inspectionsample (
		objectid,
		insp_id, sampleid, processed, idbytech, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_insp_id, p_sampleid, p_processed, p_idbytech, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_inspectionsampledetail(
	p_objectid bigint,
	p_inspsample_id uuid,p_fieldspecies varchar,p_flarvcount smallint,p_fpupcount smallint,p_feggcount smallint,p_flstages varchar,p_fdomstage varchar,p_fadultact varchar,p_labspecies varchar,p_llarvcount smallint,p_lpupcount smallint,p_leggcount smallint,p_ldomstage varchar,p_comments varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_processed smallint,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.inspectionsampledetail lv 
		WHERE lv.objectid = p_objectid
		AND lv.inspsample_id IS NOT DISTINCT FROM p_inspsample_id AND lv.fieldspecies IS NOT DISTINCT FROM p_fieldspecies AND lv.flarvcount IS NOT DISTINCT FROM p_flarvcount AND lv.fpupcount IS NOT DISTINCT FROM p_fpupcount AND lv.feggcount IS NOT DISTINCT FROM p_feggcount AND lv.flstages IS NOT DISTINCT FROM p_flstages AND lv.fdomstage IS NOT DISTINCT FROM p_fdomstage AND lv.fadultact IS NOT DISTINCT FROM p_fadultact AND lv.labspecies IS NOT DISTINCT FROM p_labspecies AND lv.llarvcount IS NOT DISTINCT FROM p_llarvcount AND lv.lpupcount IS NOT DISTINCT FROM p_lpupcount AND lv.leggcount IS NOT DISTINCT FROM p_leggcount AND lv.ldomstage IS NOT DISTINCT FROM p_ldomstage AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.processed IS NOT DISTINCT FROM p_processed AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
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
		inspsample_id, fieldspecies, flarvcount, fpupcount, feggcount, flstages, fdomstage, fadultact, labspecies, llarvcount, lpupcount, leggcount, ldomstage, comments, globalid, created_user, created_date, last_edited_user, last_edited_date, processed, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_inspsample_id, p_fieldspecies, p_flarvcount, p_fpupcount, p_feggcount, p_flstages, p_fdomstage, p_fadultact, p_labspecies, p_llarvcount, p_lpupcount, p_leggcount, p_ldomstage, p_comments, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_processed, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_linelocation(
	p_objectid bigint,
	p_name varchar,p_zone varchar,p_habitat varchar,p_priority varchar,p_usetype varchar,p_active smallint,p_description varchar,p_accessdesc varchar,p_comments varchar,p_symbology varchar,p_externalid varchar,p_acres double precision,p_nextactiondatescheduled timestamp,p_larvinspectinterval smallint,p_length_ft double precision,p_width_ft double precision,p_zone2 varchar,p_locationnumber integer,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_lastinspectdate timestamp,p_lastinspectbreeding varchar,p_lastinspectavglarvae double precision,p_lastinspectavgpupae double precision,p_lastinspectlstages varchar,p_lastinspectactiontaken varchar,p_lastinspectfieldspecies varchar,p_lasttreatdate timestamp,p_lasttreatproduct varchar,p_lasttreatqty double precision,p_lasttreatqtyunit varchar,p_hectares double precision,p_lastinspectactivity varchar,p_lasttreatactivity varchar,p_length_meters double precision,p_width_meters double precision,p_lastinspectconditions varchar,p_waterorigin varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_jurisdiction varchar,p_shape__length double precision
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.linelocation lv 
		WHERE lv.objectid = p_objectid
		AND lv.name IS NOT DISTINCT FROM p_name AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.habitat IS NOT DISTINCT FROM p_habitat AND lv.priority IS NOT DISTINCT FROM p_priority AND lv.usetype IS NOT DISTINCT FROM p_usetype AND lv.active IS NOT DISTINCT FROM p_active AND lv.description IS NOT DISTINCT FROM p_description AND lv.accessdesc IS NOT DISTINCT FROM p_accessdesc AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.symbology IS NOT DISTINCT FROM p_symbology AND lv.externalid IS NOT DISTINCT FROM p_externalid AND lv.acres IS NOT DISTINCT FROM p_acres AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled AND lv.larvinspectinterval IS NOT DISTINCT FROM p_larvinspectinterval AND lv.length_ft IS NOT DISTINCT FROM p_length_ft AND lv.width_ft IS NOT DISTINCT FROM p_width_ft AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.lastinspectdate IS NOT DISTINCT FROM p_lastinspectdate AND lv.lastinspectbreeding IS NOT DISTINCT FROM p_lastinspectbreeding AND lv.lastinspectavglarvae IS NOT DISTINCT FROM p_lastinspectavglarvae AND lv.lastinspectavgpupae IS NOT DISTINCT FROM p_lastinspectavgpupae AND lv.lastinspectlstages IS NOT DISTINCT FROM p_lastinspectlstages AND lv.lastinspectactiontaken IS NOT DISTINCT FROM p_lastinspectactiontaken AND lv.lastinspectfieldspecies IS NOT DISTINCT FROM p_lastinspectfieldspecies AND lv.lasttreatdate IS NOT DISTINCT FROM p_lasttreatdate AND lv.lasttreatproduct IS NOT DISTINCT FROM p_lasttreatproduct AND lv.lasttreatqty IS NOT DISTINCT FROM p_lasttreatqty AND lv.lasttreatqtyunit IS NOT DISTINCT FROM p_lasttreatqtyunit AND lv.hectares IS NOT DISTINCT FROM p_hectares AND lv.lastinspectactivity IS NOT DISTINCT FROM p_lastinspectactivity AND lv.lasttreatactivity IS NOT DISTINCT FROM p_lasttreatactivity AND lv.length_meters IS NOT DISTINCT FROM p_length_meters AND lv.width_meters IS NOT DISTINCT FROM p_width_meters AND lv.lastinspectconditions IS NOT DISTINCT FROM p_lastinspectconditions AND lv.waterorigin IS NOT DISTINCT FROM p_waterorigin AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction AND lv.shape__length IS NOT DISTINCT FROM p_shape__length 
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
		name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, acres, nextactiondatescheduled, larvinspectinterval, length_ft, width_ft, zone2, locationnumber, globalid, created_user, created_date, last_edited_user, last_edited_date, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, hectares, lastinspectactivity, lasttreatactivity, length_meters, width_meters, lastinspectconditions, waterorigin, creationdate, creator, editdate, editor, jurisdiction, shape__length, 
		VERSION
	) VALUES (
		p_objectid,
		p_name, p_zone, p_habitat, p_priority, p_usetype, p_active, p_description, p_accessdesc, p_comments, p_symbology, p_externalid, p_acres, p_nextactiondatescheduled, p_larvinspectinterval, p_length_ft, p_width_ft, p_zone2, p_locationnumber, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_lastinspectdate, p_lastinspectbreeding, p_lastinspectavglarvae, p_lastinspectavgpupae, p_lastinspectlstages, p_lastinspectactiontaken, p_lastinspectfieldspecies, p_lasttreatdate, p_lasttreatproduct, p_lasttreatqty, p_lasttreatqtyunit, p_hectares, p_lastinspectactivity, p_lasttreatactivity, p_length_meters, p_width_meters, p_lastinspectconditions, p_waterorigin, p_creationdate, p_creator, p_editdate, p_editor, p_jurisdiction, p_shape__length, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_locationtracking(
	p_objectid bigint,
	p_accuracy double precision,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_globalid uuid,p_fieldtech varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.locationtracking lv 
		WHERE lv.objectid = p_objectid
		AND lv.accuracy IS NOT DISTINCT FROM p_accuracy AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.locationtracking 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.locationtracking
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.locationtracking (
		objectid,
		accuracy, created_user, created_date, last_edited_user, last_edited_date, globalid, fieldtech, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_accuracy, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_globalid, p_fieldtech, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_mosquitoinspection(
	p_objectid bigint,
	p_numdips smallint,p_activity varchar,p_breeding varchar,p_totlarvae smallint,p_totpupae smallint,p_eggs smallint,p_posdips smallint,p_adultact varchar,p_lstages varchar,p_domstage varchar,p_actiontaken varchar,p_comments varchar,p_avetemp double precision,p_windspeed double precision,p_raingauge double precision,p_startdatetime timestamp,p_enddatetime timestamp,p_winddir varchar,p_avglarvae double precision,p_avgpupae double precision,p_reviewed smallint,p_reviewedby varchar,p_revieweddate timestamp,p_locationname varchar,p_zone varchar,p_recordstatus smallint,p_zone2 varchar,p_personalcontact smallint,p_tirecount smallint,p_cbcount smallint,p_containercount smallint,p_fieldspecies varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_linelocid uuid,p_pointlocid uuid,p_polygonlocid uuid,p_srid uuid,p_fieldtech varchar,p_larvaepresent smallint,p_pupaepresent smallint,p_sdid uuid,p_sitecond varchar,p_positivecontainercount smallint,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_jurisdiction varchar,p_visualmonitoring smallint,p_vmcomments varchar,p_adminaction varchar,p_ptaid uuid
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.mosquitoinspection lv 
		WHERE lv.objectid = p_objectid
		AND lv.numdips IS NOT DISTINCT FROM p_numdips AND lv.activity IS NOT DISTINCT FROM p_activity AND lv.breeding IS NOT DISTINCT FROM p_breeding AND lv.totlarvae IS NOT DISTINCT FROM p_totlarvae AND lv.totpupae IS NOT DISTINCT FROM p_totpupae AND lv.eggs IS NOT DISTINCT FROM p_eggs AND lv.posdips IS NOT DISTINCT FROM p_posdips AND lv.adultact IS NOT DISTINCT FROM p_adultact AND lv.lstages IS NOT DISTINCT FROM p_lstages AND lv.domstage IS NOT DISTINCT FROM p_domstage AND lv.actiontaken IS NOT DISTINCT FROM p_actiontaken AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.avetemp IS NOT DISTINCT FROM p_avetemp AND lv.windspeed IS NOT DISTINCT FROM p_windspeed AND lv.raingauge IS NOT DISTINCT FROM p_raingauge AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime AND lv.winddir IS NOT DISTINCT FROM p_winddir AND lv.avglarvae IS NOT DISTINCT FROM p_avglarvae AND lv.avgpupae IS NOT DISTINCT FROM p_avgpupae AND lv.reviewed IS NOT DISTINCT FROM p_reviewed AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate AND lv.locationname IS NOT DISTINCT FROM p_locationname AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.personalcontact IS NOT DISTINCT FROM p_personalcontact AND lv.tirecount IS NOT DISTINCT FROM p_tirecount AND lv.cbcount IS NOT DISTINCT FROM p_cbcount AND lv.containercount IS NOT DISTINCT FROM p_containercount AND lv.fieldspecies IS NOT DISTINCT FROM p_fieldspecies AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.linelocid IS NOT DISTINCT FROM p_linelocid AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid AND lv.polygonlocid IS NOT DISTINCT FROM p_polygonlocid AND lv.srid IS NOT DISTINCT FROM p_srid AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech AND lv.larvaepresent IS NOT DISTINCT FROM p_larvaepresent AND lv.pupaepresent IS NOT DISTINCT FROM p_pupaepresent AND lv.sdid IS NOT DISTINCT FROM p_sdid AND lv.sitecond IS NOT DISTINCT FROM p_sitecond AND lv.positivecontainercount IS NOT DISTINCT FROM p_positivecontainercount AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction AND lv.visualmonitoring IS NOT DISTINCT FROM p_visualmonitoring AND lv.vmcomments IS NOT DISTINCT FROM p_vmcomments AND lv.adminaction IS NOT DISTINCT FROM p_adminaction AND lv.ptaid IS NOT DISTINCT FROM p_ptaid 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.mosquitoinspection 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.mosquitoinspection
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.mosquitoinspection (
		objectid,
		numdips, activity, breeding, totlarvae, totpupae, eggs, posdips, adultact, lstages, domstage, actiontaken, comments, avetemp, windspeed, raingauge, startdatetime, enddatetime, winddir, avglarvae, avgpupae, reviewed, reviewedby, revieweddate, locationname, zone, recordstatus, zone2, personalcontact, tirecount, cbcount, containercount, fieldspecies, globalid, created_user, created_date, last_edited_user, last_edited_date, linelocid, pointlocid, polygonlocid, srid, fieldtech, larvaepresent, pupaepresent, sdid, sitecond, positivecontainercount, creationdate, creator, editdate, editor, jurisdiction, visualmonitoring, vmcomments, adminaction, ptaid, 
		VERSION
	) VALUES (
		p_objectid,
		p_numdips, p_activity, p_breeding, p_totlarvae, p_totpupae, p_eggs, p_posdips, p_adultact, p_lstages, p_domstage, p_actiontaken, p_comments, p_avetemp, p_windspeed, p_raingauge, p_startdatetime, p_enddatetime, p_winddir, p_avglarvae, p_avgpupae, p_reviewed, p_reviewedby, p_revieweddate, p_locationname, p_zone, p_recordstatus, p_zone2, p_personalcontact, p_tirecount, p_cbcount, p_containercount, p_fieldspecies, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_linelocid, p_pointlocid, p_polygonlocid, p_srid, p_fieldtech, p_larvaepresent, p_pupaepresent, p_sdid, p_sitecond, p_positivecontainercount, p_creationdate, p_creator, p_editdate, p_editor, p_jurisdiction, p_visualmonitoring, p_vmcomments, p_adminaction, p_ptaid, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_pointlocation(
	p_objectid bigint,
	p_name varchar,p_zone varchar,p_habitat varchar,p_priority varchar,p_usetype varchar,p_active smallint,p_description varchar,p_accessdesc varchar,p_comments varchar,p_symbology varchar,p_externalid varchar,p_nextactiondatescheduled timestamp,p_larvinspectinterval smallint,p_zone2 varchar,p_locationnumber integer,p_globalid uuid,p_stype varchar,p_lastinspectdate timestamp,p_lastinspectbreeding varchar,p_lastinspectavglarvae double precision,p_lastinspectavgpupae double precision,p_lastinspectlstages varchar,p_lastinspectactiontaken varchar,p_lastinspectfieldspecies varchar,p_lasttreatdate timestamp,p_lasttreatproduct varchar,p_lasttreatqty double precision,p_lasttreatqtyunit varchar,p_lastinspectactivity varchar,p_lasttreatactivity varchar,p_lastinspectconditions varchar,p_waterorigin varchar,p_x double precision,p_y double precision,p_assignedtech varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_jurisdiction varchar,p_deactivate_reason varchar,p_scalarpriority integer,p_sourcestatus varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.pointlocation lv 
		WHERE lv.objectid = p_objectid
		AND lv.name IS NOT DISTINCT FROM p_name AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.habitat IS NOT DISTINCT FROM p_habitat AND lv.priority IS NOT DISTINCT FROM p_priority AND lv.usetype IS NOT DISTINCT FROM p_usetype AND lv.active IS NOT DISTINCT FROM p_active AND lv.description IS NOT DISTINCT FROM p_description AND lv.accessdesc IS NOT DISTINCT FROM p_accessdesc AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.symbology IS NOT DISTINCT FROM p_symbology AND lv.externalid IS NOT DISTINCT FROM p_externalid AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled AND lv.larvinspectinterval IS NOT DISTINCT FROM p_larvinspectinterval AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.stype IS NOT DISTINCT FROM p_stype AND lv.lastinspectdate IS NOT DISTINCT FROM p_lastinspectdate AND lv.lastinspectbreeding IS NOT DISTINCT FROM p_lastinspectbreeding AND lv.lastinspectavglarvae IS NOT DISTINCT FROM p_lastinspectavglarvae AND lv.lastinspectavgpupae IS NOT DISTINCT FROM p_lastinspectavgpupae AND lv.lastinspectlstages IS NOT DISTINCT FROM p_lastinspectlstages AND lv.lastinspectactiontaken IS NOT DISTINCT FROM p_lastinspectactiontaken AND lv.lastinspectfieldspecies IS NOT DISTINCT FROM p_lastinspectfieldspecies AND lv.lasttreatdate IS NOT DISTINCT FROM p_lasttreatdate AND lv.lasttreatproduct IS NOT DISTINCT FROM p_lasttreatproduct AND lv.lasttreatqty IS NOT DISTINCT FROM p_lasttreatqty AND lv.lasttreatqtyunit IS NOT DISTINCT FROM p_lasttreatqtyunit AND lv.lastinspectactivity IS NOT DISTINCT FROM p_lastinspectactivity AND lv.lasttreatactivity IS NOT DISTINCT FROM p_lasttreatactivity AND lv.lastinspectconditions IS NOT DISTINCT FROM p_lastinspectconditions AND lv.waterorigin IS NOT DISTINCT FROM p_waterorigin AND lv.x IS NOT DISTINCT FROM p_x AND lv.y IS NOT DISTINCT FROM p_y AND lv.assignedtech IS NOT DISTINCT FROM p_assignedtech AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction AND lv.deactivate_reason IS NOT DISTINCT FROM p_deactivate_reason AND lv.scalarpriority IS NOT DISTINCT FROM p_scalarpriority AND lv.sourcestatus IS NOT DISTINCT FROM p_sourcestatus 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.pointlocation 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.pointlocation
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.pointlocation (
		objectid,
		name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, nextactiondatescheduled, larvinspectinterval, zone2, locationnumber, globalid, stype, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, lastinspectactivity, lasttreatactivity, lastinspectconditions, waterorigin, x, y, assignedtech, creationdate, creator, editdate, editor, jurisdiction, deactivate_reason, scalarpriority, sourcestatus, 
		VERSION
	) VALUES (
		p_objectid,
		p_name, p_zone, p_habitat, p_priority, p_usetype, p_active, p_description, p_accessdesc, p_comments, p_symbology, p_externalid, p_nextactiondatescheduled, p_larvinspectinterval, p_zone2, p_locationnumber, p_globalid, p_stype, p_lastinspectdate, p_lastinspectbreeding, p_lastinspectavglarvae, p_lastinspectavgpupae, p_lastinspectlstages, p_lastinspectactiontaken, p_lastinspectfieldspecies, p_lasttreatdate, p_lasttreatproduct, p_lasttreatqty, p_lasttreatqtyunit, p_lastinspectactivity, p_lasttreatactivity, p_lastinspectconditions, p_waterorigin, p_x, p_y, p_assignedtech, p_creationdate, p_creator, p_editdate, p_editor, p_jurisdiction, p_deactivate_reason, p_scalarpriority, p_sourcestatus, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_polygonlocation(
	p_objectid bigint,
	p_name varchar,p_zone varchar,p_habitat varchar,p_priority varchar,p_usetype varchar,p_active smallint,p_description varchar,p_accessdesc varchar,p_comments varchar,p_symbology varchar,p_externalid varchar,p_acres double precision,p_nextactiondatescheduled timestamp,p_larvinspectinterval smallint,p_zone2 varchar,p_locationnumber integer,p_globalid uuid,p_lastinspectdate timestamp,p_lastinspectbreeding varchar,p_lastinspectavglarvae double precision,p_lastinspectavgpupae double precision,p_lastinspectlstages varchar,p_lastinspectactiontaken varchar,p_lastinspectfieldspecies varchar,p_lasttreatdate timestamp,p_lasttreatproduct varchar,p_lasttreatqty double precision,p_lasttreatqtyunit varchar,p_hectares double precision,p_lastinspectactivity varchar,p_lasttreatactivity varchar,p_lastinspectconditions varchar,p_waterorigin varchar,p_filter varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_jurisdiction varchar,p_shape__area double precision,p_shape__length double precision
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.polygonlocation lv 
		WHERE lv.objectid = p_objectid
		AND lv.name IS NOT DISTINCT FROM p_name AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.habitat IS NOT DISTINCT FROM p_habitat AND lv.priority IS NOT DISTINCT FROM p_priority AND lv.usetype IS NOT DISTINCT FROM p_usetype AND lv.active IS NOT DISTINCT FROM p_active AND lv.description IS NOT DISTINCT FROM p_description AND lv.accessdesc IS NOT DISTINCT FROM p_accessdesc AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.symbology IS NOT DISTINCT FROM p_symbology AND lv.externalid IS NOT DISTINCT FROM p_externalid AND lv.acres IS NOT DISTINCT FROM p_acres AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled AND lv.larvinspectinterval IS NOT DISTINCT FROM p_larvinspectinterval AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.lastinspectdate IS NOT DISTINCT FROM p_lastinspectdate AND lv.lastinspectbreeding IS NOT DISTINCT FROM p_lastinspectbreeding AND lv.lastinspectavglarvae IS NOT DISTINCT FROM p_lastinspectavglarvae AND lv.lastinspectavgpupae IS NOT DISTINCT FROM p_lastinspectavgpupae AND lv.lastinspectlstages IS NOT DISTINCT FROM p_lastinspectlstages AND lv.lastinspectactiontaken IS NOT DISTINCT FROM p_lastinspectactiontaken AND lv.lastinspectfieldspecies IS NOT DISTINCT FROM p_lastinspectfieldspecies AND lv.lasttreatdate IS NOT DISTINCT FROM p_lasttreatdate AND lv.lasttreatproduct IS NOT DISTINCT FROM p_lasttreatproduct AND lv.lasttreatqty IS NOT DISTINCT FROM p_lasttreatqty AND lv.lasttreatqtyunit IS NOT DISTINCT FROM p_lasttreatqtyunit AND lv.hectares IS NOT DISTINCT FROM p_hectares AND lv.lastinspectactivity IS NOT DISTINCT FROM p_lastinspectactivity AND lv.lasttreatactivity IS NOT DISTINCT FROM p_lasttreatactivity AND lv.lastinspectconditions IS NOT DISTINCT FROM p_lastinspectconditions AND lv.waterorigin IS NOT DISTINCT FROM p_waterorigin AND lv.filter IS NOT DISTINCT FROM p_filter AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction AND lv.shape__area IS NOT DISTINCT FROM p_shape__area AND lv.shape__length IS NOT DISTINCT FROM p_shape__length 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.polygonlocation 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.polygonlocation
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.polygonlocation (
		objectid,
		name, zone, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, acres, nextactiondatescheduled, larvinspectinterval, zone2, locationnumber, globalid, lastinspectdate, lastinspectbreeding, lastinspectavglarvae, lastinspectavgpupae, lastinspectlstages, lastinspectactiontaken, lastinspectfieldspecies, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, hectares, lastinspectactivity, lasttreatactivity, lastinspectconditions, waterorigin, filter, creationdate, creator, editdate, editor, jurisdiction, shape__area, shape__length, 
		VERSION
	) VALUES (
		p_objectid,
		p_name, p_zone, p_habitat, p_priority, p_usetype, p_active, p_description, p_accessdesc, p_comments, p_symbology, p_externalid, p_acres, p_nextactiondatescheduled, p_larvinspectinterval, p_zone2, p_locationnumber, p_globalid, p_lastinspectdate, p_lastinspectbreeding, p_lastinspectavglarvae, p_lastinspectavgpupae, p_lastinspectlstages, p_lastinspectactiontaken, p_lastinspectfieldspecies, p_lasttreatdate, p_lasttreatproduct, p_lasttreatqty, p_lasttreatqtyunit, p_hectares, p_lastinspectactivity, p_lasttreatactivity, p_lastinspectconditions, p_waterorigin, p_filter, p_creationdate, p_creator, p_editdate, p_editor, p_jurisdiction, p_shape__area, p_shape__length, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_pool(
	p_objectid bigint,
	p_trapdata_id uuid,p_datesent timestamp,p_survtech varchar,p_datetested timestamp,p_testtech varchar,p_comments varchar,p_sampleid varchar,p_processed smallint,p_lab_id uuid,p_testmethod varchar,p_diseasetested varchar,p_diseasepos varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_lab varchar,p_poolyear smallint,p_gatewaysync smallint,p_vectorsurvcollectionid varchar,p_vectorsurvpoolid varchar,p_vectorsurvtrapdataid varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.pool lv 
		WHERE lv.objectid = p_objectid
		AND lv.trapdata_id IS NOT DISTINCT FROM p_trapdata_id AND lv.datesent IS NOT DISTINCT FROM p_datesent AND lv.survtech IS NOT DISTINCT FROM p_survtech AND lv.datetested IS NOT DISTINCT FROM p_datetested AND lv.testtech IS NOT DISTINCT FROM p_testtech AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.sampleid IS NOT DISTINCT FROM p_sampleid AND lv.processed IS NOT DISTINCT FROM p_processed AND lv.lab_id IS NOT DISTINCT FROM p_lab_id AND lv.testmethod IS NOT DISTINCT FROM p_testmethod AND lv.diseasetested IS NOT DISTINCT FROM p_diseasetested AND lv.diseasepos IS NOT DISTINCT FROM p_diseasepos AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.lab IS NOT DISTINCT FROM p_lab AND lv.poolyear IS NOT DISTINCT FROM p_poolyear AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync AND lv.vectorsurvcollectionid IS NOT DISTINCT FROM p_vectorsurvcollectionid AND lv.vectorsurvpoolid IS NOT DISTINCT FROM p_vectorsurvpoolid AND lv.vectorsurvtrapdataid IS NOT DISTINCT FROM p_vectorsurvtrapdataid AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
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
		trapdata_id, datesent, survtech, datetested, testtech, comments, sampleid, processed, lab_id, testmethod, diseasetested, diseasepos, globalid, created_user, created_date, last_edited_user, last_edited_date, lab, poolyear, gatewaysync, vectorsurvcollectionid, vectorsurvpoolid, vectorsurvtrapdataid, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_trapdata_id, p_datesent, p_survtech, p_datetested, p_testtech, p_comments, p_sampleid, p_processed, p_lab_id, p_testmethod, p_diseasetested, p_diseasepos, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_lab, p_poolyear, p_gatewaysync, p_vectorsurvcollectionid, p_vectorsurvpoolid, p_vectorsurvtrapdataid, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_pooldetail(
	p_objectid bigint,
	p_trapdata_id uuid,p_pool_id uuid,p_species varchar,p_females smallint,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.pooldetail lv 
		WHERE lv.objectid = p_objectid
		AND lv.trapdata_id IS NOT DISTINCT FROM p_trapdata_id AND lv.pool_id IS NOT DISTINCT FROM p_pool_id AND lv.species IS NOT DISTINCT FROM p_species AND lv.females IS NOT DISTINCT FROM p_females AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.pooldetail 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.pooldetail
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.pooldetail (
		objectid,
		trapdata_id, pool_id, species, females, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_trapdata_id, p_pool_id, p_species, p_females, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_proposedtreatmentarea(
	p_objectid bigint,
	p_method varchar,p_comments varchar,p_zone varchar,p_reviewed smallint,p_reviewedby varchar,p_revieweddate timestamp,p_zone2 varchar,p_completeddate timestamp,p_completedby varchar,p_completed smallint,p_issprayroute smallint,p_name varchar,p_acres double precision,p_globalid uuid,p_exported smallint,p_targetproduct varchar,p_targetapprate double precision,p_hectares double precision,p_lasttreatactivity varchar,p_lasttreatdate timestamp,p_lasttreatproduct varchar,p_lasttreatqty double precision,p_lasttreatqtyunit varchar,p_priority varchar,p_duedate timestamp,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_targetspecies varchar,p_shape__area double precision,p_shape__length double precision
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.proposedtreatmentarea lv 
		WHERE lv.objectid = p_objectid
		AND lv.method IS NOT DISTINCT FROM p_method AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.reviewed IS NOT DISTINCT FROM p_reviewed AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.completeddate IS NOT DISTINCT FROM p_completeddate AND lv.completedby IS NOT DISTINCT FROM p_completedby AND lv.completed IS NOT DISTINCT FROM p_completed AND lv.issprayroute IS NOT DISTINCT FROM p_issprayroute AND lv.name IS NOT DISTINCT FROM p_name AND lv.acres IS NOT DISTINCT FROM p_acres AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.exported IS NOT DISTINCT FROM p_exported AND lv.targetproduct IS NOT DISTINCT FROM p_targetproduct AND lv.targetapprate IS NOT DISTINCT FROM p_targetapprate AND lv.hectares IS NOT DISTINCT FROM p_hectares AND lv.lasttreatactivity IS NOT DISTINCT FROM p_lasttreatactivity AND lv.lasttreatdate IS NOT DISTINCT FROM p_lasttreatdate AND lv.lasttreatproduct IS NOT DISTINCT FROM p_lasttreatproduct AND lv.lasttreatqty IS NOT DISTINCT FROM p_lasttreatqty AND lv.lasttreatqtyunit IS NOT DISTINCT FROM p_lasttreatqtyunit AND lv.priority IS NOT DISTINCT FROM p_priority AND lv.duedate IS NOT DISTINCT FROM p_duedate AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.targetspecies IS NOT DISTINCT FROM p_targetspecies AND lv.shape__area IS NOT DISTINCT FROM p_shape__area AND lv.shape__length IS NOT DISTINCT FROM p_shape__length 
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
		method, comments, zone, reviewed, reviewedby, revieweddate, zone2, completeddate, completedby, completed, issprayroute, name, acres, globalid, exported, targetproduct, targetapprate, hectares, lasttreatactivity, lasttreatdate, lasttreatproduct, lasttreatqty, lasttreatqtyunit, priority, duedate, creationdate, creator, editdate, editor, targetspecies, shape__area, shape__length, 
		VERSION
	) VALUES (
		p_objectid,
		p_method, p_comments, p_zone, p_reviewed, p_reviewedby, p_revieweddate, p_zone2, p_completeddate, p_completedby, p_completed, p_issprayroute, p_name, p_acres, p_globalid, p_exported, p_targetproduct, p_targetapprate, p_hectares, p_lasttreatactivity, p_lasttreatdate, p_lasttreatproduct, p_lasttreatqty, p_lasttreatqtyunit, p_priority, p_duedate, p_creationdate, p_creator, p_editdate, p_editor, p_targetspecies, p_shape__area, p_shape__length, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_qamosquitoinspection(
	p_objectid bigint,
	p_posdips smallint,p_actiontaken varchar,p_comments varchar,p_avetemp double precision,p_windspeed double precision,p_raingauge double precision,p_globalid uuid,p_startdatetime timestamp,p_enddatetime timestamp,p_winddir varchar,p_reviewed smallint,p_reviewedby varchar,p_revieweddate timestamp,p_locationname varchar,p_zone varchar,p_recordstatus smallint,p_zone2 varchar,p_lr smallint,p_negdips smallint,p_totalacres double precision,p_acresbreeding double precision,p_fish smallint,p_sitetype varchar,p_breedingpotential varchar,p_movingwater smallint,p_nowaterever smallint,p_mosquitohabitat varchar,p_habvalue1 smallint,p_habvalue1percent smallint,p_habvalue2 smallint,p_habvalue2percent smallint,p_potential smallint,p_larvaepresent smallint,p_larvaeinsidetreatedarea smallint,p_larvaeoutsidetreatedarea smallint,p_larvaereason varchar,p_aquaticorganisms varchar,p_vegetation varchar,p_sourcereduction varchar,p_waterpresent smallint,p_watermovement1 varchar,p_watermovement1percent smallint,p_watermovement2 varchar,p_watermovement2percent smallint,p_soilconditions varchar,p_waterduration varchar,p_watersource varchar,p_waterconditions varchar,p_adultactivity smallint,p_linelocid uuid,p_pointlocid uuid,p_polygonlocid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_fieldtech varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.qamosquitoinspection lv 
		WHERE lv.objectid = p_objectid
		AND lv.posdips IS NOT DISTINCT FROM p_posdips AND lv.actiontaken IS NOT DISTINCT FROM p_actiontaken AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.avetemp IS NOT DISTINCT FROM p_avetemp AND lv.windspeed IS NOT DISTINCT FROM p_windspeed AND lv.raingauge IS NOT DISTINCT FROM p_raingauge AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime AND lv.winddir IS NOT DISTINCT FROM p_winddir AND lv.reviewed IS NOT DISTINCT FROM p_reviewed AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate AND lv.locationname IS NOT DISTINCT FROM p_locationname AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.lr IS NOT DISTINCT FROM p_lr AND lv.negdips IS NOT DISTINCT FROM p_negdips AND lv.totalacres IS NOT DISTINCT FROM p_totalacres AND lv.acresbreeding IS NOT DISTINCT FROM p_acresbreeding AND lv.fish IS NOT DISTINCT FROM p_fish AND lv.sitetype IS NOT DISTINCT FROM p_sitetype AND lv.breedingpotential IS NOT DISTINCT FROM p_breedingpotential AND lv.movingwater IS NOT DISTINCT FROM p_movingwater AND lv.nowaterever IS NOT DISTINCT FROM p_nowaterever AND lv.mosquitohabitat IS NOT DISTINCT FROM p_mosquitohabitat AND lv.habvalue1 IS NOT DISTINCT FROM p_habvalue1 AND lv.habvalue1percent IS NOT DISTINCT FROM p_habvalue1percent AND lv.habvalue2 IS NOT DISTINCT FROM p_habvalue2 AND lv.habvalue2percent IS NOT DISTINCT FROM p_habvalue2percent AND lv.potential IS NOT DISTINCT FROM p_potential AND lv.larvaepresent IS NOT DISTINCT FROM p_larvaepresent AND lv.larvaeinsidetreatedarea IS NOT DISTINCT FROM p_larvaeinsidetreatedarea AND lv.larvaeoutsidetreatedarea IS NOT DISTINCT FROM p_larvaeoutsidetreatedarea AND lv.larvaereason IS NOT DISTINCT FROM p_larvaereason AND lv.aquaticorganisms IS NOT DISTINCT FROM p_aquaticorganisms AND lv.vegetation IS NOT DISTINCT FROM p_vegetation AND lv.sourcereduction IS NOT DISTINCT FROM p_sourcereduction AND lv.waterpresent IS NOT DISTINCT FROM p_waterpresent AND lv.watermovement1 IS NOT DISTINCT FROM p_watermovement1 AND lv.watermovement1percent IS NOT DISTINCT FROM p_watermovement1percent AND lv.watermovement2 IS NOT DISTINCT FROM p_watermovement2 AND lv.watermovement2percent IS NOT DISTINCT FROM p_watermovement2percent AND lv.soilconditions IS NOT DISTINCT FROM p_soilconditions AND lv.waterduration IS NOT DISTINCT FROM p_waterduration AND lv.watersource IS NOT DISTINCT FROM p_watersource AND lv.waterconditions IS NOT DISTINCT FROM p_waterconditions AND lv.adultactivity IS NOT DISTINCT FROM p_adultactivity AND lv.linelocid IS NOT DISTINCT FROM p_linelocid AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid AND lv.polygonlocid IS NOT DISTINCT FROM p_polygonlocid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.qamosquitoinspection 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.qamosquitoinspection
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.qamosquitoinspection (
		objectid,
		posdips, actiontaken, comments, avetemp, windspeed, raingauge, globalid, startdatetime, enddatetime, winddir, reviewed, reviewedby, revieweddate, locationname, zone, recordstatus, zone2, lr, negdips, totalacres, acresbreeding, fish, sitetype, breedingpotential, movingwater, nowaterever, mosquitohabitat, habvalue1, habvalue1percent, habvalue2, habvalue2percent, potential, larvaepresent, larvaeinsidetreatedarea, larvaeoutsidetreatedarea, larvaereason, aquaticorganisms, vegetation, sourcereduction, waterpresent, watermovement1, watermovement1percent, watermovement2, watermovement2percent, soilconditions, waterduration, watersource, waterconditions, adultactivity, linelocid, pointlocid, polygonlocid, created_user, created_date, last_edited_user, last_edited_date, fieldtech, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_posdips, p_actiontaken, p_comments, p_avetemp, p_windspeed, p_raingauge, p_globalid, p_startdatetime, p_enddatetime, p_winddir, p_reviewed, p_reviewedby, p_revieweddate, p_locationname, p_zone, p_recordstatus, p_zone2, p_lr, p_negdips, p_totalacres, p_acresbreeding, p_fish, p_sitetype, p_breedingpotential, p_movingwater, p_nowaterever, p_mosquitohabitat, p_habvalue1, p_habvalue1percent, p_habvalue2, p_habvalue2percent, p_potential, p_larvaepresent, p_larvaeinsidetreatedarea, p_larvaeoutsidetreatedarea, p_larvaereason, p_aquaticorganisms, p_vegetation, p_sourcereduction, p_waterpresent, p_watermovement1, p_watermovement1percent, p_watermovement2, p_watermovement2percent, p_soilconditions, p_waterduration, p_watersource, p_waterconditions, p_adultactivity, p_linelocid, p_pointlocid, p_polygonlocid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_fieldtech, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
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
	p_nextactiondatescheduled timestamp without time zone ,
	p_locationnumber integer,
	p_lastinspectdate timestamp without time zone,
	p_lastinspectspecies varchar,
	p_lastinspectaction varchar,
	p_lastinspectconditions varchar,
	p_lastinspectrodentevidence varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp without time zone,
	p_last_edited_user varchar,
	p_last_edited_date timestamp without time zone,
	p_creationdate timestamp without time zone,
	p_creator varchar,
	p_editdate timestamp without time zone,
	p_editor varchar,
	p_jurisdiction varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.rodentlocation lv 
		WHERE lv.objectid = p_objectid
		AND lv.locationname IS NOT DISTINCT FROM p_locationname AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.habitat IS NOT DISTINCT FROM p_habitat AND lv.priority IS NOT DISTINCT FROM p_priority AND lv.usetype IS NOT DISTINCT FROM p_usetype AND lv.active IS NOT DISTINCT FROM p_active AND lv.description IS NOT DISTINCT FROM p_description AND lv.accessdesc IS NOT DISTINCT FROM p_accessdesc AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.symbology IS NOT DISTINCT FROM p_symbology AND lv.externalid IS NOT DISTINCT FROM p_externalid AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber AND lv.lastinspectdate IS NOT DISTINCT FROM p_lastinspectdate AND lv.lastinspectspecies IS NOT DISTINCT FROM p_lastinspectspecies AND lv.lastinspectaction IS NOT DISTINCT FROM p_lastinspectaction AND lv.lastinspectconditions IS NOT DISTINCT FROM p_lastinspectconditions AND lv.lastinspectrodentevidence IS NOT DISTINCT FROM p_lastinspectrodentevidence AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction 
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
		locationname, zone, zone2, habitat, priority, usetype, active, description, accessdesc, comments, symbology, externalid, nextactiondatescheduled, locationnumber, lastinspectdate, lastinspectspecies, lastinspectaction, lastinspectconditions, lastinspectrodentevidence, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, jurisdiction, 
		VERSION
	) VALUES (
		p_objectid,
		p_locationname, p_zone, p_zone2, p_habitat, p_priority, p_usetype, p_active, p_description, p_accessdesc, p_comments, p_symbology, p_externalid, p_nextactiondatescheduled, p_locationnumber, p_lastinspectdate, p_lastinspectspecies, p_lastinspectaction, p_lastinspectconditions, p_lastinspectrodentevidence, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_creationdate, p_creator, p_editdate, p_editor, p_jurisdiction, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_samplecollection(
	p_objectid bigint,
	p_loc_id uuid,p_startdatetime timestamp,p_enddatetime timestamp,p_sitecond varchar,p_sampleid varchar,p_survtech varchar,p_datesent timestamp,p_datetested timestamp,p_testtech varchar,p_comments varchar,p_processed smallint,p_sampletype varchar,p_samplecond varchar,p_species varchar,p_sex varchar,p_avetemp double precision,p_windspeed double precision,p_winddir varchar,p_raingauge double precision,p_activity varchar,p_testmethod varchar,p_diseasetested varchar,p_diseasepos varchar,p_reviewed smallint,p_reviewedby varchar,p_revieweddate timestamp,p_locationname varchar,p_zone varchar,p_recordstatus smallint,p_zone2 varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_lab varchar,p_fieldtech varchar,p_flockid uuid,p_samplecount smallint,p_chickenid uuid,p_gatewaysync smallint,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.samplecollection lv 
		WHERE lv.objectid = p_objectid
		AND lv.loc_id IS NOT DISTINCT FROM p_loc_id AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime AND lv.sitecond IS NOT DISTINCT FROM p_sitecond AND lv.sampleid IS NOT DISTINCT FROM p_sampleid AND lv.survtech IS NOT DISTINCT FROM p_survtech AND lv.datesent IS NOT DISTINCT FROM p_datesent AND lv.datetested IS NOT DISTINCT FROM p_datetested AND lv.testtech IS NOT DISTINCT FROM p_testtech AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.processed IS NOT DISTINCT FROM p_processed AND lv.sampletype IS NOT DISTINCT FROM p_sampletype AND lv.samplecond IS NOT DISTINCT FROM p_samplecond AND lv.species IS NOT DISTINCT FROM p_species AND lv.sex IS NOT DISTINCT FROM p_sex AND lv.avetemp IS NOT DISTINCT FROM p_avetemp AND lv.windspeed IS NOT DISTINCT FROM p_windspeed AND lv.winddir IS NOT DISTINCT FROM p_winddir AND lv.raingauge IS NOT DISTINCT FROM p_raingauge AND lv.activity IS NOT DISTINCT FROM p_activity AND lv.testmethod IS NOT DISTINCT FROM p_testmethod AND lv.diseasetested IS NOT DISTINCT FROM p_diseasetested AND lv.diseasepos IS NOT DISTINCT FROM p_diseasepos AND lv.reviewed IS NOT DISTINCT FROM p_reviewed AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate AND lv.locationname IS NOT DISTINCT FROM p_locationname AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.lab IS NOT DISTINCT FROM p_lab AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech AND lv.flockid IS NOT DISTINCT FROM p_flockid AND lv.samplecount IS NOT DISTINCT FROM p_samplecount AND lv.chickenid IS NOT DISTINCT FROM p_chickenid AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.samplecollection 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.samplecollection
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.samplecollection (
		objectid,
		loc_id, startdatetime, enddatetime, sitecond, sampleid, survtech, datesent, datetested, testtech, comments, processed, sampletype, samplecond, species, sex, avetemp, windspeed, winddir, raingauge, activity, testmethod, diseasetested, diseasepos, reviewed, reviewedby, revieweddate, locationname, zone, recordstatus, zone2, globalid, created_user, created_date, last_edited_user, last_edited_date, lab, fieldtech, flockid, samplecount, chickenid, gatewaysync, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_loc_id, p_startdatetime, p_enddatetime, p_sitecond, p_sampleid, p_survtech, p_datesent, p_datetested, p_testtech, p_comments, p_processed, p_sampletype, p_samplecond, p_species, p_sex, p_avetemp, p_windspeed, p_winddir, p_raingauge, p_activity, p_testmethod, p_diseasetested, p_diseasepos, p_reviewed, p_reviewedby, p_revieweddate, p_locationname, p_zone, p_recordstatus, p_zone2, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_lab, p_fieldtech, p_flockid, p_samplecount, p_chickenid, p_gatewaysync, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_samplelocation(
	p_objectid bigint,
	p_name varchar,p_zone varchar,p_habitat varchar,p_priority varchar,p_usetype varchar,p_active smallint,p_description varchar,p_accessdesc varchar,p_comments varchar,p_externalid varchar,p_nextactiondatescheduled timestamp,p_zone2 varchar,p_locationnumber integer,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_gatewaysync smallint,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.samplelocation lv 
		WHERE lv.objectid = p_objectid
		AND lv.name IS NOT DISTINCT FROM p_name AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.habitat IS NOT DISTINCT FROM p_habitat AND lv.priority IS NOT DISTINCT FROM p_priority AND lv.usetype IS NOT DISTINCT FROM p_usetype AND lv.active IS NOT DISTINCT FROM p_active AND lv.description IS NOT DISTINCT FROM p_description AND lv.accessdesc IS NOT DISTINCT FROM p_accessdesc AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.externalid IS NOT DISTINCT FROM p_externalid AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.samplelocation 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.samplelocation
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.samplelocation (
		objectid,
		name, zone, habitat, priority, usetype, active, description, accessdesc, comments, externalid, nextactiondatescheduled, zone2, locationnumber, globalid, created_user, created_date, last_edited_user, last_edited_date, gatewaysync, creationdate, creator, editdate, editor, 
		VERSION
	) VALUES (
		p_objectid,
		p_name, p_zone, p_habitat, p_priority, p_usetype, p_active, p_description, p_accessdesc, p_comments, p_externalid, p_nextactiondatescheduled, p_zone2, p_locationnumber, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_gatewaysync, p_creationdate, p_creator, p_editdate, p_editor, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_servicerequest(
	p_objectid bigint,
	p_recdatetime timestamp,p_source varchar,p_entrytech varchar,p_priority varchar,p_supervisor varchar,p_assignedtech varchar,p_status varchar,p_clranon smallint,p_clrfname varchar,p_clrphone1 varchar,p_clrphone2 varchar,p_clremail varchar,p_clrcompany varchar,p_clraddr1 varchar,p_clraddr2 varchar,p_clrcity varchar,p_clrstate varchar,p_clrzip varchar,p_clrother varchar,p_clrcontpref varchar,p_reqcompany varchar,p_reqaddr1 varchar,p_reqaddr2 varchar,p_reqcity varchar,p_reqstate varchar,p_reqzip varchar,p_reqcrossst varchar,p_reqsubdiv varchar,p_reqmapgrid varchar,p_reqpermission smallint,p_reqtarget varchar,p_reqdescr varchar,p_reqnotesfortech varchar,p_reqnotesforcust varchar,p_reqfldnotes varchar,p_reqprogramactions varchar,p_datetimeclosed timestamp,p_techclosed varchar,p_sr_number integer,p_reviewed smallint,p_reviewedby varchar,p_revieweddate timestamp,p_accepted smallint,p_accepteddate timestamp,p_rejectedby varchar,p_rejecteddate timestamp,p_rejectedreason varchar,p_duedate timestamp,p_acceptedby varchar,p_comments varchar,p_estcompletedate timestamp,p_nextaction varchar,p_recordstatus smallint,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_firstresponsedate timestamp,p_responsedaycount smallint,p_allowed varchar,p_xvalue varchar,p_yvalue varchar,p_validx varchar,p_validy varchar,p_externalid varchar,p_externalerror varchar,p_pointlocid uuid,p_notified smallint,p_notifieddate timestamp,p_scheduled smallint,p_scheduleddate timestamp,p_dog integer,p_schedule_period varchar,p_schedule_notes varchar,p_spanish integer,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_issuesreported varchar,p_jurisdiction varchar,p_notificationtimestamp varchar,p_zone varchar,p_zone2 varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.servicerequest lv 
		WHERE lv.objectid = p_objectid
		AND lv.recdatetime IS NOT DISTINCT FROM p_recdatetime AND lv.source IS NOT DISTINCT FROM p_source AND lv.entrytech IS NOT DISTINCT FROM p_entrytech AND lv.priority IS NOT DISTINCT FROM p_priority AND lv.supervisor IS NOT DISTINCT FROM p_supervisor AND lv.assignedtech IS NOT DISTINCT FROM p_assignedtech AND lv.status IS NOT DISTINCT FROM p_status AND lv.clranon IS NOT DISTINCT FROM p_clranon AND lv.clrfname IS NOT DISTINCT FROM p_clrfname AND lv.clrphone1 IS NOT DISTINCT FROM p_clrphone1 AND lv.clrphone2 IS NOT DISTINCT FROM p_clrphone2 AND lv.clremail IS NOT DISTINCT FROM p_clremail AND lv.clrcompany IS NOT DISTINCT FROM p_clrcompany AND lv.clraddr1 IS NOT DISTINCT FROM p_clraddr1 AND lv.clraddr2 IS NOT DISTINCT FROM p_clraddr2 AND lv.clrcity IS NOT DISTINCT FROM p_clrcity AND lv.clrstate IS NOT DISTINCT FROM p_clrstate AND lv.clrzip IS NOT DISTINCT FROM p_clrzip AND lv.clrother IS NOT DISTINCT FROM p_clrother AND lv.clrcontpref IS NOT DISTINCT FROM p_clrcontpref AND lv.reqcompany IS NOT DISTINCT FROM p_reqcompany AND lv.reqaddr1 IS NOT DISTINCT FROM p_reqaddr1 AND lv.reqaddr2 IS NOT DISTINCT FROM p_reqaddr2 AND lv.reqcity IS NOT DISTINCT FROM p_reqcity AND lv.reqstate IS NOT DISTINCT FROM p_reqstate AND lv.reqzip IS NOT DISTINCT FROM p_reqzip AND lv.reqcrossst IS NOT DISTINCT FROM p_reqcrossst AND lv.reqsubdiv IS NOT DISTINCT FROM p_reqsubdiv AND lv.reqmapgrid IS NOT DISTINCT FROM p_reqmapgrid AND lv.reqpermission IS NOT DISTINCT FROM p_reqpermission AND lv.reqtarget IS NOT DISTINCT FROM p_reqtarget AND lv.reqdescr IS NOT DISTINCT FROM p_reqdescr AND lv.reqnotesfortech IS NOT DISTINCT FROM p_reqnotesfortech AND lv.reqnotesforcust IS NOT DISTINCT FROM p_reqnotesforcust AND lv.reqfldnotes IS NOT DISTINCT FROM p_reqfldnotes AND lv.reqprogramactions IS NOT DISTINCT FROM p_reqprogramactions AND lv.datetimeclosed IS NOT DISTINCT FROM p_datetimeclosed AND lv.techclosed IS NOT DISTINCT FROM p_techclosed AND lv.sr_number IS NOT DISTINCT FROM p_sr_number AND lv.reviewed IS NOT DISTINCT FROM p_reviewed AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate AND lv.accepted IS NOT DISTINCT FROM p_accepted AND lv.accepteddate IS NOT DISTINCT FROM p_accepteddate AND lv.rejectedby IS NOT DISTINCT FROM p_rejectedby AND lv.rejecteddate IS NOT DISTINCT FROM p_rejecteddate AND lv.rejectedreason IS NOT DISTINCT FROM p_rejectedreason AND lv.duedate IS NOT DISTINCT FROM p_duedate AND lv.acceptedby IS NOT DISTINCT FROM p_acceptedby AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.estcompletedate IS NOT DISTINCT FROM p_estcompletedate AND lv.nextaction IS NOT DISTINCT FROM p_nextaction AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.firstresponsedate IS NOT DISTINCT FROM p_firstresponsedate AND lv.responsedaycount IS NOT DISTINCT FROM p_responsedaycount AND lv.allowed IS NOT DISTINCT FROM p_allowed AND lv.xvalue IS NOT DISTINCT FROM p_xvalue AND lv.yvalue IS NOT DISTINCT FROM p_yvalue AND lv.validx IS NOT DISTINCT FROM p_validx AND lv.validy IS NOT DISTINCT FROM p_validy AND lv.externalid IS NOT DISTINCT FROM p_externalid AND lv.externalerror IS NOT DISTINCT FROM p_externalerror AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid AND lv.notified IS NOT DISTINCT FROM p_notified AND lv.notifieddate IS NOT DISTINCT FROM p_notifieddate AND lv.scheduled IS NOT DISTINCT FROM p_scheduled AND lv.scheduleddate IS NOT DISTINCT FROM p_scheduleddate AND lv.dog IS NOT DISTINCT FROM p_dog AND lv.schedule_period IS NOT DISTINCT FROM p_schedule_period AND lv.schedule_notes IS NOT DISTINCT FROM p_schedule_notes AND lv.spanish IS NOT DISTINCT FROM p_spanish AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.issuesreported IS NOT DISTINCT FROM p_issuesreported AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction AND lv.notificationtimestamp IS NOT DISTINCT FROM p_notificationtimestamp AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.servicerequest 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.servicerequest
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.servicerequest (
		objectid,
		recdatetime, source, entrytech, priority, supervisor, assignedtech, status, clranon, clrfname, clrphone1, clrphone2, clremail, clrcompany, clraddr1, clraddr2, clrcity, clrstate, clrzip, clrother, clrcontpref, reqcompany, reqaddr1, reqaddr2, reqcity, reqstate, reqzip, reqcrossst, reqsubdiv, reqmapgrid, reqpermission, reqtarget, reqdescr, reqnotesfortech, reqnotesforcust, reqfldnotes, reqprogramactions, datetimeclosed, techclosed, sr_number, reviewed, reviewedby, revieweddate, accepted, accepteddate, rejectedby, rejecteddate, rejectedreason, duedate, acceptedby, comments, estcompletedate, nextaction, recordstatus, globalid, created_user, created_date, last_edited_user, last_edited_date, firstresponsedate, responsedaycount, allowed, xvalue, yvalue, validx, validy, externalid, externalerror, pointlocid, notified, notifieddate, scheduled, scheduleddate, dog, schedule_period, schedule_notes, spanish, creationdate, creator, editdate, editor, issuesreported, jurisdiction, notificationtimestamp, zone, zone2, 
		VERSION
	) VALUES (
		p_objectid,
		p_recdatetime, p_source, p_entrytech, p_priority, p_supervisor, p_assignedtech, p_status, p_clranon, p_clrfname, p_clrphone1, p_clrphone2, p_clremail, p_clrcompany, p_clraddr1, p_clraddr2, p_clrcity, p_clrstate, p_clrzip, p_clrother, p_clrcontpref, p_reqcompany, p_reqaddr1, p_reqaddr2, p_reqcity, p_reqstate, p_reqzip, p_reqcrossst, p_reqsubdiv, p_reqmapgrid, p_reqpermission, p_reqtarget, p_reqdescr, p_reqnotesfortech, p_reqnotesforcust, p_reqfldnotes, p_reqprogramactions, p_datetimeclosed, p_techclosed, p_sr_number, p_reviewed, p_reviewedby, p_revieweddate, p_accepted, p_accepteddate, p_rejectedby, p_rejecteddate, p_rejectedreason, p_duedate, p_acceptedby, p_comments, p_estcompletedate, p_nextaction, p_recordstatus, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_firstresponsedate, p_responsedaycount, p_allowed, p_xvalue, p_yvalue, p_validx, p_validy, p_externalid, p_externalerror, p_pointlocid, p_notified, p_notifieddate, p_scheduled, p_scheduleddate, p_dog, p_schedule_period, p_schedule_notes, p_spanish, p_creationdate, p_creator, p_editdate, p_editor, p_issuesreported, p_jurisdiction, p_notificationtimestamp, p_zone, p_zone2, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_speciesabundance(
	p_objectid bigint,
	p_trapdata_id uuid,p_species varchar,p_males smallint,p_unknown smallint,p_bloodedfem smallint,p_gravidfem smallint,p_larvae smallint,p_poolstogen smallint,p_processed smallint,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_pupae smallint,p_eggs smallint,p_females integer,p_total integer,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_yearweek integer,p_globalzscore double precision,p_r7score double precision,p_r8score double precision,p_h3r7 varchar,p_h3r8 varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.speciesabundance lv 
		WHERE lv.objectid = p_objectid
		AND lv.trapdata_id IS NOT DISTINCT FROM p_trapdata_id AND lv.species IS NOT DISTINCT FROM p_species AND lv.males IS NOT DISTINCT FROM p_males AND lv.unknown IS NOT DISTINCT FROM p_unknown AND lv.bloodedfem IS NOT DISTINCT FROM p_bloodedfem AND lv.gravidfem IS NOT DISTINCT FROM p_gravidfem AND lv.larvae IS NOT DISTINCT FROM p_larvae AND lv.poolstogen IS NOT DISTINCT FROM p_poolstogen AND lv.processed IS NOT DISTINCT FROM p_processed AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.pupae IS NOT DISTINCT FROM p_pupae AND lv.eggs IS NOT DISTINCT FROM p_eggs AND lv.females IS NOT DISTINCT FROM p_females AND lv.total IS NOT DISTINCT FROM p_total AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.yearweek IS NOT DISTINCT FROM p_yearweek AND lv.globalzscore IS NOT DISTINCT FROM p_globalzscore AND lv.r7score IS NOT DISTINCT FROM p_r7score AND lv.r8score IS NOT DISTINCT FROM p_r8score AND lv.h3r7 IS NOT DISTINCT FROM p_h3r7 AND lv.h3r8 IS NOT DISTINCT FROM p_h3r8 
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
		trapdata_id, species, males, unknown, bloodedfem, gravidfem, larvae, poolstogen, processed, globalid, created_user, created_date, last_edited_user, last_edited_date, pupae, eggs, females, total, creationdate, creator, editdate, editor, yearweek, globalzscore, r7score, r8score, h3r7, h3r8, 
		VERSION
	) VALUES (
		p_objectid,
		p_trapdata_id, p_species, p_males, p_unknown, p_bloodedfem, p_gravidfem, p_larvae, p_poolstogen, p_processed, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_pupae, p_eggs, p_females, p_total, p_creationdate, p_creator, p_editdate, p_editor, p_yearweek, p_globalzscore, p_r7score, p_r8score, p_h3r7, p_h3r8, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_stormdrain(
	p_objectid bigint,
	p_nexttreatmentdate timestamp,p_lasttreatdate timestamp,p_lastaction varchar,p_symbology varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_laststatus varchar,p_zone varchar,p_zone2 varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_type varchar,p_jurisdiction varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.stormdrain lv 
		WHERE lv.objectid = p_objectid
		AND lv.nexttreatmentdate IS NOT DISTINCT FROM p_nexttreatmentdate AND lv.lasttreatdate IS NOT DISTINCT FROM p_lasttreatdate AND lv.lastaction IS NOT DISTINCT FROM p_lastaction AND lv.symbology IS NOT DISTINCT FROM p_symbology AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.laststatus IS NOT DISTINCT FROM p_laststatus AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.type IS NOT DISTINCT FROM p_type AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction 
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
		nexttreatmentdate, lasttreatdate, lastaction, symbology, globalid, created_user, created_date, last_edited_user, last_edited_date, laststatus, zone, zone2, creationdate, creator, editdate, editor, type, jurisdiction, 
		VERSION
	) VALUES (
		p_objectid,
		p_nexttreatmentdate, p_lasttreatdate, p_lastaction, p_symbology, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_laststatus, p_zone, p_zone2, p_creationdate, p_creator, p_editdate, p_editor, p_type, p_jurisdiction, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_timecard(
	p_objectid bigint,
	p_activity varchar,p_startdatetime timestamp,p_enddatetime timestamp,p_comments varchar,p_externalid varchar,p_equiptype varchar,p_locationname varchar,p_zone varchar,p_zone2 varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_linelocid uuid,p_pointlocid uuid,p_polygonlocid uuid,p_lclocid uuid,p_samplelocid uuid,p_srid uuid,p_traplocid uuid,p_fieldtech varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_rodentlocid uuid
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.timecard lv 
		WHERE lv.objectid = p_objectid
		AND lv.activity IS NOT DISTINCT FROM p_activity AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.externalid IS NOT DISTINCT FROM p_externalid AND lv.equiptype IS NOT DISTINCT FROM p_equiptype AND lv.locationname IS NOT DISTINCT FROM p_locationname AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.linelocid IS NOT DISTINCT FROM p_linelocid AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid AND lv.polygonlocid IS NOT DISTINCT FROM p_polygonlocid AND lv.lclocid IS NOT DISTINCT FROM p_lclocid AND lv.samplelocid IS NOT DISTINCT FROM p_samplelocid AND lv.srid IS NOT DISTINCT FROM p_srid AND lv.traplocid IS NOT DISTINCT FROM p_traplocid AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.rodentlocid IS NOT DISTINCT FROM p_rodentlocid 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.timecard 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.timecard
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.timecard (
		objectid,
		activity, startdatetime, enddatetime, comments, externalid, equiptype, locationname, zone, zone2, globalid, created_user, created_date, last_edited_user, last_edited_date, linelocid, pointlocid, polygonlocid, lclocid, samplelocid, srid, traplocid, fieldtech, creationdate, creator, editdate, editor, rodentlocid, 
		VERSION
	) VALUES (
		p_objectid,
		p_activity, p_startdatetime, p_enddatetime, p_comments, p_externalid, p_equiptype, p_locationname, p_zone, p_zone2, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_linelocid, p_pointlocid, p_polygonlocid, p_lclocid, p_samplelocid, p_srid, p_traplocid, p_fieldtech, p_creationdate, p_creator, p_editdate, p_editor, p_rodentlocid, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_trapdata(
	p_objectid bigint,
	p_traptype varchar,p_trapactivitytype varchar,p_startdatetime timestamp,p_enddatetime timestamp,p_comments varchar,p_idbytech varchar,p_sortbytech varchar,p_processed smallint,p_sitecond varchar,p_locationname varchar,p_recordstatus smallint,p_reviewed smallint,p_reviewedby varchar,p_revieweddate timestamp,p_trapcondition varchar,p_trapnights smallint,p_zone varchar,p_zone2 varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_srid uuid,p_fieldtech varchar,p_gatewaysync smallint,p_loc_id uuid,p_voltage double precision,p_winddir varchar,p_windspeed double precision,p_avetemp double precision,p_raingauge double precision,p_lr smallint,p_field integer,p_vectorsurvtrapdataid varchar,p_vectorsurvtraplocationid varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_lure varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.trapdata lv 
		WHERE lv.objectid = p_objectid
		AND lv.traptype IS NOT DISTINCT FROM p_traptype AND lv.trapactivitytype IS NOT DISTINCT FROM p_trapactivitytype AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.idbytech IS NOT DISTINCT FROM p_idbytech AND lv.sortbytech IS NOT DISTINCT FROM p_sortbytech AND lv.processed IS NOT DISTINCT FROM p_processed AND lv.sitecond IS NOT DISTINCT FROM p_sitecond AND lv.locationname IS NOT DISTINCT FROM p_locationname AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus AND lv.reviewed IS NOT DISTINCT FROM p_reviewed AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate AND lv.trapcondition IS NOT DISTINCT FROM p_trapcondition AND lv.trapnights IS NOT DISTINCT FROM p_trapnights AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.srid IS NOT DISTINCT FROM p_srid AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync AND lv.loc_id IS NOT DISTINCT FROM p_loc_id AND lv.voltage IS NOT DISTINCT FROM p_voltage AND lv.winddir IS NOT DISTINCT FROM p_winddir AND lv.windspeed IS NOT DISTINCT FROM p_windspeed AND lv.avetemp IS NOT DISTINCT FROM p_avetemp AND lv.raingauge IS NOT DISTINCT FROM p_raingauge AND lv.lr IS NOT DISTINCT FROM p_lr AND lv.field IS NOT DISTINCT FROM p_field AND lv.vectorsurvtrapdataid IS NOT DISTINCT FROM p_vectorsurvtrapdataid AND lv.vectorsurvtraplocationid IS NOT DISTINCT FROM p_vectorsurvtraplocationid AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.lure IS NOT DISTINCT FROM p_lure 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.trapdata 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.trapdata
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.trapdata (
		objectid,
		traptype, trapactivitytype, startdatetime, enddatetime, comments, idbytech, sortbytech, processed, sitecond, locationname, recordstatus, reviewed, reviewedby, revieweddate, trapcondition, trapnights, zone, zone2, globalid, created_user, created_date, last_edited_user, last_edited_date, srid, fieldtech, gatewaysync, loc_id, voltage, winddir, windspeed, avetemp, raingauge, lr, field, vectorsurvtrapdataid, vectorsurvtraplocationid, creationdate, creator, editdate, editor, lure, 
		VERSION
	) VALUES (
		p_objectid,
		p_traptype, p_trapactivitytype, p_startdatetime, p_enddatetime, p_comments, p_idbytech, p_sortbytech, p_processed, p_sitecond, p_locationname, p_recordstatus, p_reviewed, p_reviewedby, p_revieweddate, p_trapcondition, p_trapnights, p_zone, p_zone2, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_srid, p_fieldtech, p_gatewaysync, p_loc_id, p_voltage, p_winddir, p_windspeed, p_avetemp, p_raingauge, p_lr, p_field, p_vectorsurvtrapdataid, p_vectorsurvtraplocationid, p_creationdate, p_creator, p_editdate, p_editor, p_lure, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_traplocation(
	p_objectid bigint,
	p_name varchar,p_zone varchar,p_habitat varchar,p_priority varchar,p_usetype varchar,p_active smallint,p_description varchar,p_accessdesc varchar,p_comments varchar,p_externalid varchar,p_nextactiondatescheduled timestamp,p_zone2 varchar,p_locationnumber integer,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_gatewaysync smallint,p_route integer,p_set_dow integer,p_route_order integer,p_vectorsurvsiteid varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_h3r7 varchar,p_h3r8 varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.traplocation lv 
		WHERE lv.objectid = p_objectid
		AND lv.name IS NOT DISTINCT FROM p_name AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.habitat IS NOT DISTINCT FROM p_habitat AND lv.priority IS NOT DISTINCT FROM p_priority AND lv.usetype IS NOT DISTINCT FROM p_usetype AND lv.active IS NOT DISTINCT FROM p_active AND lv.description IS NOT DISTINCT FROM p_description AND lv.accessdesc IS NOT DISTINCT FROM p_accessdesc AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.externalid IS NOT DISTINCT FROM p_externalid AND lv.nextactiondatescheduled IS NOT DISTINCT FROM p_nextactiondatescheduled AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.locationnumber IS NOT DISTINCT FROM p_locationnumber AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.gatewaysync IS NOT DISTINCT FROM p_gatewaysync AND lv.route IS NOT DISTINCT FROM p_route AND lv.set_dow IS NOT DISTINCT FROM p_set_dow AND lv.route_order IS NOT DISTINCT FROM p_route_order AND lv.vectorsurvsiteid IS NOT DISTINCT FROM p_vectorsurvsiteid AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.h3r7 IS NOT DISTINCT FROM p_h3r7 AND lv.h3r8 IS NOT DISTINCT FROM p_h3r8 
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
		name, zone, habitat, priority, usetype, active, description, accessdesc, comments, externalid, nextactiondatescheduled, zone2, locationnumber, globalid, created_user, created_date, last_edited_user, last_edited_date, gatewaysync, route, set_dow, route_order, vectorsurvsiteid, creationdate, creator, editdate, editor, h3r7, h3r8, 
		VERSION
	) VALUES (
		p_objectid,
		p_name, p_zone, p_habitat, p_priority, p_usetype, p_active, p_description, p_accessdesc, p_comments, p_externalid, p_nextactiondatescheduled, p_zone2, p_locationnumber, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_gatewaysync, p_route, p_set_dow, p_route_order, p_vectorsurvsiteid, p_creationdate, p_creator, p_editdate, p_editor, p_h3r7, p_h3r8, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_treatment(
	p_objectid bigint,
	p_activity varchar,p_treatarea double precision,p_areaunit varchar,p_product varchar,p_qty double precision,p_qtyunit varchar,p_method varchar,p_equiptype varchar,p_comments varchar,p_avetemp double precision,p_windspeed double precision,p_winddir varchar,p_raingauge double precision,p_startdatetime timestamp,p_enddatetime timestamp,p_insp_id uuid,p_reviewed smallint,p_reviewedby varchar,p_revieweddate timestamp,p_locationname varchar,p_zone varchar,p_warningoverride smallint,p_recordstatus smallint,p_zone2 varchar,p_treatacres double precision,p_tirecount smallint,p_cbcount smallint,p_containercount smallint,p_globalid uuid,p_treatmentlength double precision,p_treatmenthours double precision,p_treatmentlengthunits varchar,p_linelocid uuid,p_pointlocid uuid,p_polygonlocid uuid,p_srid uuid,p_sdid uuid,p_barrierrouteid uuid,p_ulvrouteid uuid,p_fieldtech varchar,p_ptaid uuid,p_flowrate double precision,p_habitat varchar,p_treathectares double precision,p_invloc varchar,p_temp_sitecond varchar,p_sitecond varchar,p_totalcostprodcut double precision,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_targetspecies varchar
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.treatment lv 
		WHERE lv.objectid = p_objectid
		AND lv.activity IS NOT DISTINCT FROM p_activity AND lv.treatarea IS NOT DISTINCT FROM p_treatarea AND lv.areaunit IS NOT DISTINCT FROM p_areaunit AND lv.product IS NOT DISTINCT FROM p_product AND lv.qty IS NOT DISTINCT FROM p_qty AND lv.qtyunit IS NOT DISTINCT FROM p_qtyunit AND lv.method IS NOT DISTINCT FROM p_method AND lv.equiptype IS NOT DISTINCT FROM p_equiptype AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.avetemp IS NOT DISTINCT FROM p_avetemp AND lv.windspeed IS NOT DISTINCT FROM p_windspeed AND lv.winddir IS NOT DISTINCT FROM p_winddir AND lv.raingauge IS NOT DISTINCT FROM p_raingauge AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime AND lv.insp_id IS NOT DISTINCT FROM p_insp_id AND lv.reviewed IS NOT DISTINCT FROM p_reviewed AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate AND lv.locationname IS NOT DISTINCT FROM p_locationname AND lv.zone IS NOT DISTINCT FROM p_zone AND lv.warningoverride IS NOT DISTINCT FROM p_warningoverride AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus AND lv.zone2 IS NOT DISTINCT FROM p_zone2 AND lv.treatacres IS NOT DISTINCT FROM p_treatacres AND lv.tirecount IS NOT DISTINCT FROM p_tirecount AND lv.cbcount IS NOT DISTINCT FROM p_cbcount AND lv.containercount IS NOT DISTINCT FROM p_containercount AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.treatmentlength IS NOT DISTINCT FROM p_treatmentlength AND lv.treatmenthours IS NOT DISTINCT FROM p_treatmenthours AND lv.treatmentlengthunits IS NOT DISTINCT FROM p_treatmentlengthunits AND lv.linelocid IS NOT DISTINCT FROM p_linelocid AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid AND lv.polygonlocid IS NOT DISTINCT FROM p_polygonlocid AND lv.srid IS NOT DISTINCT FROM p_srid AND lv.sdid IS NOT DISTINCT FROM p_sdid AND lv.barrierrouteid IS NOT DISTINCT FROM p_barrierrouteid AND lv.ulvrouteid IS NOT DISTINCT FROM p_ulvrouteid AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech AND lv.ptaid IS NOT DISTINCT FROM p_ptaid AND lv.flowrate IS NOT DISTINCT FROM p_flowrate AND lv.habitat IS NOT DISTINCT FROM p_habitat AND lv.treathectares IS NOT DISTINCT FROM p_treathectares AND lv.invloc IS NOT DISTINCT FROM p_invloc AND lv.temp_sitecond IS NOT DISTINCT FROM p_temp_sitecond AND lv.sitecond IS NOT DISTINCT FROM p_sitecond AND lv.totalcostprodcut IS NOT DISTINCT FROM p_totalcostprodcut AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.targetspecies IS NOT DISTINCT FROM p_targetspecies 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.treatment 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.treatment
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.treatment (
		objectid,
		activity, treatarea, areaunit, product, qty, qtyunit, method, equiptype, comments, avetemp, windspeed, winddir, raingauge, startdatetime, enddatetime, insp_id, reviewed, reviewedby, revieweddate, locationname, zone, warningoverride, recordstatus, zone2, treatacres, tirecount, cbcount, containercount, globalid, treatmentlength, treatmenthours, treatmentlengthunits, linelocid, pointlocid, polygonlocid, srid, sdid, barrierrouteid, ulvrouteid, fieldtech, ptaid, flowrate, habitat, treathectares, invloc, temp_sitecond, sitecond, totalcostprodcut, creationdate, creator, editdate, editor, targetspecies, 
		VERSION
	) VALUES (
		p_objectid,
		p_activity, p_treatarea, p_areaunit, p_product, p_qty, p_qtyunit, p_method, p_equiptype, p_comments, p_avetemp, p_windspeed, p_winddir, p_raingauge, p_startdatetime, p_enddatetime, p_insp_id, p_reviewed, p_reviewedby, p_revieweddate, p_locationname, p_zone, p_warningoverride, p_recordstatus, p_zone2, p_treatacres, p_tirecount, p_cbcount, p_containercount, p_globalid, p_treatmentlength, p_treatmenthours, p_treatmentlengthunits, p_linelocid, p_pointlocid, p_polygonlocid, p_srid, p_sdid, p_barrierrouteid, p_ulvrouteid, p_fieldtech, p_ptaid, p_flowrate, p_habitat, p_treathectares, p_invloc, p_temp_sitecond, p_sitecond, p_totalcostprodcut, p_creationdate, p_creator, p_editdate, p_editor, p_targetspecies, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_treatmentarea(
	p_objectid bigint,
	p_treat_id uuid,p_session_id uuid,p_treatdate timestamp,p_comments varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_notified smallint,p_type varchar,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_shape__area double precision,p_shape__length double precision
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.treatmentarea lv 
		WHERE lv.objectid = p_objectid
		AND lv.treat_id IS NOT DISTINCT FROM p_treat_id AND lv.session_id IS NOT DISTINCT FROM p_session_id AND lv.treatdate IS NOT DISTINCT FROM p_treatdate AND lv.comments IS NOT DISTINCT FROM p_comments AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.notified IS NOT DISTINCT FROM p_notified AND lv.type IS NOT DISTINCT FROM p_type AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.shape__area IS NOT DISTINCT FROM p_shape__area AND lv.shape__length IS NOT DISTINCT FROM p_shape__length 
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
		treat_id, session_id, treatdate, comments, globalid, created_user, created_date, last_edited_user, last_edited_date, notified, type, creationdate, creator, editdate, editor, shape__area, shape__length, 
		VERSION
	) VALUES (
		p_objectid,
		p_treat_id, p_session_id, p_treatdate, p_comments, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_notified, p_type, p_creationdate, p_creator, p_editdate, p_editor, p_shape__area, p_shape__length, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_zones2(
	p_objectid bigint,
	p_name varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_shape__area double precision,p_shape__length double precision
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.zones2 lv 
		WHERE lv.objectid = p_objectid
		AND lv.name IS NOT DISTINCT FROM p_name AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.shape__area IS NOT DISTINCT FROM p_shape__area AND lv.shape__length IS NOT DISTINCT FROM p_shape__length 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.zones2 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.zones2
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.zones2 (
		objectid,
		name, globalid, created_user, created_date, last_edited_user, last_edited_date, creationdate, creator, editdate, editor, shape__area, shape__length, 
		VERSION
	) VALUES (
		p_objectid,
		p_name, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_creationdate, p_creator, p_editdate, p_editor, p_shape__area, p_shape__length, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_zones(
	p_objectid bigint,
	p_name varchar,p_globalid uuid,p_created_user varchar,p_created_date timestamp,p_last_edited_user varchar,p_last_edited_date timestamp,p_active integer,p_creationdate timestamp,p_creator varchar,p_editdate timestamp,p_editor varchar,p_shape__area double precision,p_shape__length double precision
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.zones lv 
		WHERE lv.objectid = p_objectid
		AND lv.name IS NOT DISTINCT FROM p_name AND lv.globalid IS NOT DISTINCT FROM p_globalid AND lv.created_user IS NOT DISTINCT FROM p_created_user AND lv.created_date IS NOT DISTINCT FROM p_created_date AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date AND lv.active IS NOT DISTINCT FROM p_active AND lv.creationdate IS NOT DISTINCT FROM p_creationdate AND lv.creator IS NOT DISTINCT FROM p_creator AND lv.editdate IS NOT DISTINCT FROM p_editdate AND lv.editor IS NOT DISTINCT FROM p_editor AND lv.shape__area IS NOT DISTINCT FROM p_shape__area AND lv.shape__length IS NOT DISTINCT FROM p_shape__length 
		ORDER BY VERSION DESC LIMIT 1
	) INTO v_changes_exist;
	
	-- If no changes, return false with current version
	IF NOT v_changes_exist THEN
		RETURN QUERY 
			SELECT 
				FALSE AS row_inserted, 
				(SELECT VERSION FROM fieldseeker.zones 
				 WHERE objectid = p_objectid ORDER BY VERSION DESC LIMIT 1) AS version_num;
		RETURN;
	END IF;
	
	-- Calculate next version
	SELECT COALESCE(MAX(VERSION) + 1, 1) INTO v_next_version
	FROM fieldseeker.zones
	WHERE objectid = p_objectid;
	
	-- Insert new version
	INSERT INTO fieldseeker.zones (
		objectid,
		name, globalid, created_user, created_date, last_edited_user, last_edited_date, active, creationdate, creator, editdate, editor, shape__area, shape__length, 
		VERSION
	) VALUES (
		p_objectid,
		p_name, p_globalid, p_created_user, p_created_date, p_last_edited_user, p_last_edited_date, p_active, p_creationdate, p_creator, p_editdate, p_editor, p_shape__area, p_shape__length, 
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
-- +goose Down
DROP FUNCTION fieldseeker.insert_containerrelate;
DROP FUNCTION fieldseeker.insert_fieldscoutinglog;
DROP FUNCTION fieldseeker.insert_habitatrelate;
DROP FUNCTION fieldseeker.insert_inspectionsample;
DROP FUNCTION fieldseeker.insert_inspectionsampledetail;
DROP FUNCTION fieldseeker.insert_linelocation;
DROP FUNCTION fieldseeker.insert_locationtracking;
DROP FUNCTION fieldseeker.insert_mosquitoinspection;
DROP FUNCTION fieldseeker.insert_pointlocation;
DROP FUNCTION fieldseeker.insert_polygonlocation;
DROP FUNCTION fieldseeker.insert_pool;
DROP FUNCTION fieldseeker.insert_pooldetail;
DROP FUNCTION fieldseeker.insert_proposedtreatmentarea;
DROP FUNCTION fieldseeker.insert_qamosquitoinspection;
DROP FUNCTION fieldseeker.insert_rodentlocation;
DROP FUNCTION fieldseeker.insert_samplecollection;
DROP FUNCTION fieldseeker.insert_samplelocation;
DROP FUNCTION fieldseeker.insert_servicerequest;
DROP FUNCTION fieldseeker.insert_speciesabundance;
DROP FUNCTION fieldseeker.insert_stormdrain;
DROP FUNCTION fieldseeker.insert_timecard;
DROP FUNCTION fieldseeker.insert_trapdata;
DROP FUNCTION fieldseeker.insert_traplocation;
DROP FUNCTION fieldseeker.insert_treatment;
DROP FUNCTION fieldseeker.insert_treatmentarea;
DROP FUNCTION fieldseeker.insert_zones2;
DROP FUNCTION fieldseeker.insert_zones;
