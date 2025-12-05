
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_qamosquitoinspection(
	p_objectid bigint,
	
	p_posdips smallint,
	p_actiontaken varchar,
	p_comments varchar,
	p_avetemp double precision,
	p_windspeed double precision,
	p_raingauge double precision,
	p_globalid uuid,
	p_startdatetime timestamp,
	p_enddatetime timestamp,
	p_winddir varchar,
	p_reviewed smallint,
	p_reviewedby varchar,
	p_revieweddate timestamp,
	p_locationname varchar,
	p_zone varchar,
	p_recordstatus smallint,
	p_zone2 varchar,
	p_lr smallint,
	p_negdips smallint,
	p_totalacres double precision,
	p_acresbreeding double precision,
	p_fish smallint,
	p_sitetype varchar,
	p_breedingpotential varchar,
	p_movingwater smallint,
	p_nowaterever smallint,
	p_mosquitohabitat varchar,
	p_habvalue1 smallint,
	p_habvalue1percent smallint,
	p_habvalue2 smallint,
	p_habvalue2percent smallint,
	p_potential smallint,
	p_larvaepresent smallint,
	p_larvaeinsidetreatedarea smallint,
	p_larvaeoutsidetreatedarea smallint,
	p_larvaereason varchar,
	p_aquaticorganisms varchar,
	p_vegetation varchar,
	p_sourcereduction varchar,
	p_waterpresent smallint,
	p_watermovement1 varchar,
	p_watermovement1percent smallint,
	p_watermovement2 varchar,
	p_watermovement2percent smallint,
	p_soilconditions varchar,
	p_waterduration varchar,
	p_watersource varchar,
	p_waterconditions varchar,
	p_adultactivity smallint,
	p_linelocid uuid,
	p_pointlocid uuid,
	p_polygonlocid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_fieldtech varchar,
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
		SELECT 1 FROM fieldseeker.qamosquitoinspection lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.posdips IS NOT DISTINCT FROM p_posdips 
		AND lv.actiontaken IS NOT DISTINCT FROM p_actiontaken 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.avetemp IS NOT DISTINCT FROM p_avetemp 
		AND lv.windspeed IS NOT DISTINCT FROM p_windspeed 
		AND lv.raingauge IS NOT DISTINCT FROM p_raingauge 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime 
		AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime 
		AND lv.winddir IS NOT DISTINCT FROM p_winddir 
		AND lv.reviewed IS NOT DISTINCT FROM p_reviewed 
		AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby 
		AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate 
		AND lv.locationname IS NOT DISTINCT FROM p_locationname 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.lr IS NOT DISTINCT FROM p_lr 
		AND lv.negdips IS NOT DISTINCT FROM p_negdips 
		AND lv.totalacres IS NOT DISTINCT FROM p_totalacres 
		AND lv.acresbreeding IS NOT DISTINCT FROM p_acresbreeding 
		AND lv.fish IS NOT DISTINCT FROM p_fish 
		AND lv.sitetype IS NOT DISTINCT FROM p_sitetype 
		AND lv.breedingpotential IS NOT DISTINCT FROM p_breedingpotential 
		AND lv.movingwater IS NOT DISTINCT FROM p_movingwater 
		AND lv.nowaterever IS NOT DISTINCT FROM p_nowaterever 
		AND lv.mosquitohabitat IS NOT DISTINCT FROM p_mosquitohabitat 
		AND lv.habvalue1 IS NOT DISTINCT FROM p_habvalue1 
		AND lv.habvalue1percent IS NOT DISTINCT FROM p_habvalue1percent 
		AND lv.habvalue2 IS NOT DISTINCT FROM p_habvalue2 
		AND lv.habvalue2percent IS NOT DISTINCT FROM p_habvalue2percent 
		AND lv.potential IS NOT DISTINCT FROM p_potential 
		AND lv.larvaepresent IS NOT DISTINCT FROM p_larvaepresent 
		AND lv.larvaeinsidetreatedarea IS NOT DISTINCT FROM p_larvaeinsidetreatedarea 
		AND lv.larvaeoutsidetreatedarea IS NOT DISTINCT FROM p_larvaeoutsidetreatedarea 
		AND lv.larvaereason IS NOT DISTINCT FROM p_larvaereason 
		AND lv.aquaticorganisms IS NOT DISTINCT FROM p_aquaticorganisms 
		AND lv.vegetation IS NOT DISTINCT FROM p_vegetation 
		AND lv.sourcereduction IS NOT DISTINCT FROM p_sourcereduction 
		AND lv.waterpresent IS NOT DISTINCT FROM p_waterpresent 
		AND lv.watermovement1 IS NOT DISTINCT FROM p_watermovement1 
		AND lv.watermovement1percent IS NOT DISTINCT FROM p_watermovement1percent 
		AND lv.watermovement2 IS NOT DISTINCT FROM p_watermovement2 
		AND lv.watermovement2percent IS NOT DISTINCT FROM p_watermovement2percent 
		AND lv.soilconditions IS NOT DISTINCT FROM p_soilconditions 
		AND lv.waterduration IS NOT DISTINCT FROM p_waterduration 
		AND lv.watersource IS NOT DISTINCT FROM p_watersource 
		AND lv.waterconditions IS NOT DISTINCT FROM p_waterconditions 
		AND lv.adultactivity IS NOT DISTINCT FROM p_adultactivity 
		AND lv.linelocid IS NOT DISTINCT FROM p_linelocid 
		AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid 
		AND lv.polygonlocid IS NOT DISTINCT FROM p_polygonlocid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech 
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
		
		posdips,
		actiontaken,
		comments,
		avetemp,
		windspeed,
		raingauge,
		globalid,
		startdatetime,
		enddatetime,
		winddir,
		reviewed,
		reviewedby,
		revieweddate,
		locationname,
		zone,
		recordstatus,
		zone2,
		lr,
		negdips,
		totalacres,
		acresbreeding,
		fish,
		sitetype,
		breedingpotential,
		movingwater,
		nowaterever,
		mosquitohabitat,
		habvalue1,
		habvalue1percent,
		habvalue2,
		habvalue2percent,
		potential,
		larvaepresent,
		larvaeinsidetreatedarea,
		larvaeoutsidetreatedarea,
		larvaereason,
		aquaticorganisms,
		vegetation,
		sourcereduction,
		waterpresent,
		watermovement1,
		watermovement1percent,
		watermovement2,
		watermovement2percent,
		soilconditions,
		waterduration,
		watersource,
		waterconditions,
		adultactivity,
		linelocid,
		pointlocid,
		polygonlocid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		fieldtech,
		creationdate,
		creator,
		editdate,
		editor,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_posdips,
		p_actiontaken,
		p_comments,
		p_avetemp,
		p_windspeed,
		p_raingauge,
		p_globalid,
		p_startdatetime,
		p_enddatetime,
		p_winddir,
		p_reviewed,
		p_reviewedby,
		p_revieweddate,
		p_locationname,
		p_zone,
		p_recordstatus,
		p_zone2,
		p_lr,
		p_negdips,
		p_totalacres,
		p_acresbreeding,
		p_fish,
		p_sitetype,
		p_breedingpotential,
		p_movingwater,
		p_nowaterever,
		p_mosquitohabitat,
		p_habvalue1,
		p_habvalue1percent,
		p_habvalue2,
		p_habvalue2percent,
		p_potential,
		p_larvaepresent,
		p_larvaeinsidetreatedarea,
		p_larvaeoutsidetreatedarea,
		p_larvaereason,
		p_aquaticorganisms,
		p_vegetation,
		p_sourcereduction,
		p_waterpresent,
		p_watermovement1,
		p_watermovement1percent,
		p_watermovement2,
		p_watermovement2percent,
		p_soilconditions,
		p_waterduration,
		p_watersource,
		p_waterconditions,
		p_adultactivity,
		p_linelocid,
		p_pointlocid,
		p_polygonlocid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_fieldtech,
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
