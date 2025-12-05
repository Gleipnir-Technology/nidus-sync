
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_mosquitoinspection(
	p_objectid bigint,
	
	p_numdips smallint,
	p_activity varchar,
	p_breeding varchar,
	p_totlarvae smallint,
	p_totpupae smallint,
	p_eggs smallint,
	p_posdips smallint,
	p_adultact varchar,
	p_lstages varchar,
	p_domstage varchar,
	p_actiontaken varchar,
	p_comments varchar,
	p_avetemp double precision,
	p_windspeed double precision,
	p_raingauge double precision,
	p_startdatetime timestamp,
	p_enddatetime timestamp,
	p_winddir varchar,
	p_avglarvae double precision,
	p_avgpupae double precision,
	p_reviewed smallint,
	p_reviewedby varchar,
	p_revieweddate timestamp,
	p_locationname varchar,
	p_zone varchar,
	p_recordstatus smallint,
	p_zone2 varchar,
	p_personalcontact smallint,
	p_tirecount smallint,
	p_cbcount smallint,
	p_containercount smallint,
	p_fieldspecies varchar,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_linelocid uuid,
	p_pointlocid uuid,
	p_polygonlocid uuid,
	p_srid uuid,
	p_fieldtech varchar,
	p_larvaepresent smallint,
	p_pupaepresent smallint,
	p_sdid uuid,
	p_sitecond varchar,
	p_positivecontainercount smallint,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_jurisdiction varchar,
	p_visualmonitoring smallint,
	p_vmcomments varchar,
	p_adminaction varchar,
	p_ptaid uuid,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.mosquitoinspection lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.numdips IS NOT DISTINCT FROM p_numdips 
		AND lv.activity IS NOT DISTINCT FROM p_activity 
		AND lv.breeding IS NOT DISTINCT FROM p_breeding 
		AND lv.totlarvae IS NOT DISTINCT FROM p_totlarvae 
		AND lv.totpupae IS NOT DISTINCT FROM p_totpupae 
		AND lv.eggs IS NOT DISTINCT FROM p_eggs 
		AND lv.posdips IS NOT DISTINCT FROM p_posdips 
		AND lv.adultact IS NOT DISTINCT FROM p_adultact 
		AND lv.lstages IS NOT DISTINCT FROM p_lstages 
		AND lv.domstage IS NOT DISTINCT FROM p_domstage 
		AND lv.actiontaken IS NOT DISTINCT FROM p_actiontaken 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.avetemp IS NOT DISTINCT FROM p_avetemp 
		AND lv.windspeed IS NOT DISTINCT FROM p_windspeed 
		AND lv.raingauge IS NOT DISTINCT FROM p_raingauge 
		AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime 
		AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime 
		AND lv.winddir IS NOT DISTINCT FROM p_winddir 
		AND lv.avglarvae IS NOT DISTINCT FROM p_avglarvae 
		AND lv.avgpupae IS NOT DISTINCT FROM p_avgpupae 
		AND lv.reviewed IS NOT DISTINCT FROM p_reviewed 
		AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby 
		AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate 
		AND lv.locationname IS NOT DISTINCT FROM p_locationname 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.personalcontact IS NOT DISTINCT FROM p_personalcontact 
		AND lv.tirecount IS NOT DISTINCT FROM p_tirecount 
		AND lv.cbcount IS NOT DISTINCT FROM p_cbcount 
		AND lv.containercount IS NOT DISTINCT FROM p_containercount 
		AND lv.fieldspecies IS NOT DISTINCT FROM p_fieldspecies 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.linelocid IS NOT DISTINCT FROM p_linelocid 
		AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid 
		AND lv.polygonlocid IS NOT DISTINCT FROM p_polygonlocid 
		AND lv.srid IS NOT DISTINCT FROM p_srid 
		AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech 
		AND lv.larvaepresent IS NOT DISTINCT FROM p_larvaepresent 
		AND lv.pupaepresent IS NOT DISTINCT FROM p_pupaepresent 
		AND lv.sdid IS NOT DISTINCT FROM p_sdid 
		AND lv.sitecond IS NOT DISTINCT FROM p_sitecond 
		AND lv.positivecontainercount IS NOT DISTINCT FROM p_positivecontainercount 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction 
		AND lv.visualmonitoring IS NOT DISTINCT FROM p_visualmonitoring 
		AND lv.vmcomments IS NOT DISTINCT FROM p_vmcomments 
		AND lv.adminaction IS NOT DISTINCT FROM p_adminaction 
		AND lv.ptaid IS NOT DISTINCT FROM p_ptaid 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
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
		
		numdips,
		activity,
		breeding,
		totlarvae,
		totpupae,
		eggs,
		posdips,
		adultact,
		lstages,
		domstage,
		actiontaken,
		comments,
		avetemp,
		windspeed,
		raingauge,
		startdatetime,
		enddatetime,
		winddir,
		avglarvae,
		avgpupae,
		reviewed,
		reviewedby,
		revieweddate,
		locationname,
		zone,
		recordstatus,
		zone2,
		personalcontact,
		tirecount,
		cbcount,
		containercount,
		fieldspecies,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		linelocid,
		pointlocid,
		polygonlocid,
		srid,
		fieldtech,
		larvaepresent,
		pupaepresent,
		sdid,
		sitecond,
		positivecontainercount,
		creationdate,
		creator,
		editdate,
		editor,
		jurisdiction,
		visualmonitoring,
		vmcomments,
		adminaction,
		ptaid,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_numdips,
		p_activity,
		p_breeding,
		p_totlarvae,
		p_totpupae,
		p_eggs,
		p_posdips,
		p_adultact,
		p_lstages,
		p_domstage,
		p_actiontaken,
		p_comments,
		p_avetemp,
		p_windspeed,
		p_raingauge,
		p_startdatetime,
		p_enddatetime,
		p_winddir,
		p_avglarvae,
		p_avgpupae,
		p_reviewed,
		p_reviewedby,
		p_revieweddate,
		p_locationname,
		p_zone,
		p_recordstatus,
		p_zone2,
		p_personalcontact,
		p_tirecount,
		p_cbcount,
		p_containercount,
		p_fieldspecies,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_linelocid,
		p_pointlocid,
		p_polygonlocid,
		p_srid,
		p_fieldtech,
		p_larvaepresent,
		p_pupaepresent,
		p_sdid,
		p_sitecond,
		p_positivecontainercount,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_jurisdiction,
		p_visualmonitoring,
		p_vmcomments,
		p_adminaction,
		p_ptaid,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
