
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_treatment(
	p_objectid bigint,
	
	p_activity varchar,
	p_treatarea double precision,
	p_areaunit varchar,
	p_product varchar,
	p_qty double precision,
	p_qtyunit varchar,
	p_method varchar,
	p_equiptype varchar,
	p_comments varchar,
	p_avetemp double precision,
	p_windspeed double precision,
	p_winddir varchar,
	p_raingauge double precision,
	p_startdatetime timestamp,
	p_enddatetime timestamp,
	p_insp_id uuid,
	p_reviewed smallint,
	p_reviewedby varchar,
	p_revieweddate timestamp,
	p_locationname varchar,
	p_zone varchar,
	p_warningoverride smallint,
	p_recordstatus smallint,
	p_zone2 varchar,
	p_treatacres double precision,
	p_tirecount smallint,
	p_cbcount smallint,
	p_containercount smallint,
	p_globalid uuid,
	p_treatmentlength double precision,
	p_treatmenthours double precision,
	p_treatmentlengthunits varchar,
	p_linelocid uuid,
	p_pointlocid uuid,
	p_polygonlocid uuid,
	p_srid uuid,
	p_sdid uuid,
	p_barrierrouteid uuid,
	p_ulvrouteid uuid,
	p_fieldtech varchar,
	p_ptaid uuid,
	p_flowrate double precision,
	p_habitat varchar,
	p_treathectares double precision,
	p_invloc varchar,
	p_temp_sitecond varchar,
	p_sitecond varchar,
	p_totalcostprodcut double precision,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_targetspecies varchar,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.treatment lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.activity IS NOT DISTINCT FROM p_activity 
		AND lv.treatarea IS NOT DISTINCT FROM p_treatarea 
		AND lv.areaunit IS NOT DISTINCT FROM p_areaunit 
		AND lv.product IS NOT DISTINCT FROM p_product 
		AND lv.qty IS NOT DISTINCT FROM p_qty 
		AND lv.qtyunit IS NOT DISTINCT FROM p_qtyunit 
		AND lv.method IS NOT DISTINCT FROM p_method 
		AND lv.equiptype IS NOT DISTINCT FROM p_equiptype 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.avetemp IS NOT DISTINCT FROM p_avetemp 
		AND lv.windspeed IS NOT DISTINCT FROM p_windspeed 
		AND lv.winddir IS NOT DISTINCT FROM p_winddir 
		AND lv.raingauge IS NOT DISTINCT FROM p_raingauge 
		AND lv.startdatetime IS NOT DISTINCT FROM p_startdatetime 
		AND lv.enddatetime IS NOT DISTINCT FROM p_enddatetime 
		AND lv.insp_id IS NOT DISTINCT FROM p_insp_id 
		AND lv.reviewed IS NOT DISTINCT FROM p_reviewed 
		AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby 
		AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate 
		AND lv.locationname IS NOT DISTINCT FROM p_locationname 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.warningoverride IS NOT DISTINCT FROM p_warningoverride 
		AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.treatacres IS NOT DISTINCT FROM p_treatacres 
		AND lv.tirecount IS NOT DISTINCT FROM p_tirecount 
		AND lv.cbcount IS NOT DISTINCT FROM p_cbcount 
		AND lv.containercount IS NOT DISTINCT FROM p_containercount 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.treatmentlength IS NOT DISTINCT FROM p_treatmentlength 
		AND lv.treatmenthours IS NOT DISTINCT FROM p_treatmenthours 
		AND lv.treatmentlengthunits IS NOT DISTINCT FROM p_treatmentlengthunits 
		AND lv.linelocid IS NOT DISTINCT FROM p_linelocid 
		AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid 
		AND lv.polygonlocid IS NOT DISTINCT FROM p_polygonlocid 
		AND lv.srid IS NOT DISTINCT FROM p_srid 
		AND lv.sdid IS NOT DISTINCT FROM p_sdid 
		AND lv.barrierrouteid IS NOT DISTINCT FROM p_barrierrouteid 
		AND lv.ulvrouteid IS NOT DISTINCT FROM p_ulvrouteid 
		AND lv.fieldtech IS NOT DISTINCT FROM p_fieldtech 
		AND lv.ptaid IS NOT DISTINCT FROM p_ptaid 
		AND lv.flowrate IS NOT DISTINCT FROM p_flowrate 
		AND lv.habitat IS NOT DISTINCT FROM p_habitat 
		AND lv.treathectares IS NOT DISTINCT FROM p_treathectares 
		AND lv.invloc IS NOT DISTINCT FROM p_invloc 
		AND lv.temp_sitecond IS NOT DISTINCT FROM p_temp_sitecond 
		AND lv.sitecond IS NOT DISTINCT FROM p_sitecond 
		AND lv.totalcostprodcut IS NOT DISTINCT FROM p_totalcostprodcut 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.targetspecies IS NOT DISTINCT FROM p_targetspecies 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
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
		
		activity,
		treatarea,
		areaunit,
		product,
		qty,
		qtyunit,
		method,
		equiptype,
		comments,
		avetemp,
		windspeed,
		winddir,
		raingauge,
		startdatetime,
		enddatetime,
		insp_id,
		reviewed,
		reviewedby,
		revieweddate,
		locationname,
		zone,
		warningoverride,
		recordstatus,
		zone2,
		treatacres,
		tirecount,
		cbcount,
		containercount,
		globalid,
		treatmentlength,
		treatmenthours,
		treatmentlengthunits,
		linelocid,
		pointlocid,
		polygonlocid,
		srid,
		sdid,
		barrierrouteid,
		ulvrouteid,
		fieldtech,
		ptaid,
		flowrate,
		habitat,
		treathectares,
		invloc,
		temp_sitecond,
		sitecond,
		totalcostprodcut,
		creationdate,
		creator,
		editdate,
		editor,
		targetspecies,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_activity,
		p_treatarea,
		p_areaunit,
		p_product,
		p_qty,
		p_qtyunit,
		p_method,
		p_equiptype,
		p_comments,
		p_avetemp,
		p_windspeed,
		p_winddir,
		p_raingauge,
		p_startdatetime,
		p_enddatetime,
		p_insp_id,
		p_reviewed,
		p_reviewedby,
		p_revieweddate,
		p_locationname,
		p_zone,
		p_warningoverride,
		p_recordstatus,
		p_zone2,
		p_treatacres,
		p_tirecount,
		p_cbcount,
		p_containercount,
		p_globalid,
		p_treatmentlength,
		p_treatmenthours,
		p_treatmentlengthunits,
		p_linelocid,
		p_pointlocid,
		p_polygonlocid,
		p_srid,
		p_sdid,
		p_barrierrouteid,
		p_ulvrouteid,
		p_fieldtech,
		p_ptaid,
		p_flowrate,
		p_habitat,
		p_treathectares,
		p_invloc,
		p_temp_sitecond,
		p_sitecond,
		p_totalcostprodcut,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_targetspecies,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
