
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fieldseeker.insert_servicerequest(
	p_objectid bigint,
	
	p_recdatetime timestamp,
	p_source varchar,
	p_entrytech varchar,
	p_priority varchar,
	p_supervisor varchar,
	p_assignedtech varchar,
	p_status varchar,
	p_clranon smallint,
	p_clrfname varchar,
	p_clrphone1 varchar,
	p_clrphone2 varchar,
	p_clremail varchar,
	p_clrcompany varchar,
	p_clraddr1 varchar,
	p_clraddr2 varchar,
	p_clrcity varchar,
	p_clrstate varchar,
	p_clrzip varchar,
	p_clrother varchar,
	p_clrcontpref varchar,
	p_reqcompany varchar,
	p_reqaddr1 varchar,
	p_reqaddr2 varchar,
	p_reqcity varchar,
	p_reqstate varchar,
	p_reqzip varchar,
	p_reqcrossst varchar,
	p_reqsubdiv varchar,
	p_reqmapgrid varchar,
	p_reqpermission smallint,
	p_reqtarget varchar,
	p_reqdescr varchar,
	p_reqnotesfortech varchar,
	p_reqnotesforcust varchar,
	p_reqfldnotes varchar,
	p_reqprogramactions varchar,
	p_datetimeclosed timestamp,
	p_techclosed varchar,
	p_sr_number integer,
	p_reviewed smallint,
	p_reviewedby varchar,
	p_revieweddate timestamp,
	p_accepted smallint,
	p_accepteddate timestamp,
	p_rejectedby varchar,
	p_rejecteddate timestamp,
	p_rejectedreason varchar,
	p_duedate timestamp,
	p_acceptedby varchar,
	p_comments varchar,
	p_estcompletedate timestamp,
	p_nextaction varchar,
	p_recordstatus smallint,
	p_globalid uuid,
	p_created_user varchar,
	p_created_date timestamp,
	p_last_edited_user varchar,
	p_last_edited_date timestamp,
	p_firstresponsedate timestamp,
	p_responsedaycount smallint,
	p_allowed varchar,
	p_xvalue varchar,
	p_yvalue varchar,
	p_validx varchar,
	p_validy varchar,
	p_externalid varchar,
	p_externalerror varchar,
	p_pointlocid uuid,
	p_notified smallint,
	p_notifieddate timestamp,
	p_scheduled smallint,
	p_scheduleddate timestamp,
	p_dog integer,
	p_schedule_period varchar,
	p_schedule_notes varchar,
	p_spanish integer,
	p_creationdate timestamp,
	p_creator varchar,
	p_editdate timestamp,
	p_editor varchar,
	p_issuesreported varchar,
	p_jurisdiction varchar,
	p_notificationtimestamp varchar,
	p_zone varchar,
	p_zone2 varchar,
	p_geometry jsonb,
	p_geospatial geometry
) RETURNS TABLE(row_inserted boolean, version_num integer) AS $$
DECLARE
	v_next_version integer;
	v_changes_exist boolean;
BEGIN
	-- Check if changes exist
	SELECT NOT EXISTS (
		SELECT 1 FROM fieldseeker.servicerequest lv 
		WHERE lv.objectid = p_objectid
		
		AND lv.recdatetime IS NOT DISTINCT FROM p_recdatetime 
		AND lv.source IS NOT DISTINCT FROM p_source 
		AND lv.entrytech IS NOT DISTINCT FROM p_entrytech 
		AND lv.priority IS NOT DISTINCT FROM p_priority 
		AND lv.supervisor IS NOT DISTINCT FROM p_supervisor 
		AND lv.assignedtech IS NOT DISTINCT FROM p_assignedtech 
		AND lv.status IS NOT DISTINCT FROM p_status 
		AND lv.clranon IS NOT DISTINCT FROM p_clranon 
		AND lv.clrfname IS NOT DISTINCT FROM p_clrfname 
		AND lv.clrphone1 IS NOT DISTINCT FROM p_clrphone1 
		AND lv.clrphone2 IS NOT DISTINCT FROM p_clrphone2 
		AND lv.clremail IS NOT DISTINCT FROM p_clremail 
		AND lv.clrcompany IS NOT DISTINCT FROM p_clrcompany 
		AND lv.clraddr1 IS NOT DISTINCT FROM p_clraddr1 
		AND lv.clraddr2 IS NOT DISTINCT FROM p_clraddr2 
		AND lv.clrcity IS NOT DISTINCT FROM p_clrcity 
		AND lv.clrstate IS NOT DISTINCT FROM p_clrstate 
		AND lv.clrzip IS NOT DISTINCT FROM p_clrzip 
		AND lv.clrother IS NOT DISTINCT FROM p_clrother 
		AND lv.clrcontpref IS NOT DISTINCT FROM p_clrcontpref 
		AND lv.reqcompany IS NOT DISTINCT FROM p_reqcompany 
		AND lv.reqaddr1 IS NOT DISTINCT FROM p_reqaddr1 
		AND lv.reqaddr2 IS NOT DISTINCT FROM p_reqaddr2 
		AND lv.reqcity IS NOT DISTINCT FROM p_reqcity 
		AND lv.reqstate IS NOT DISTINCT FROM p_reqstate 
		AND lv.reqzip IS NOT DISTINCT FROM p_reqzip 
		AND lv.reqcrossst IS NOT DISTINCT FROM p_reqcrossst 
		AND lv.reqsubdiv IS NOT DISTINCT FROM p_reqsubdiv 
		AND lv.reqmapgrid IS NOT DISTINCT FROM p_reqmapgrid 
		AND lv.reqpermission IS NOT DISTINCT FROM p_reqpermission 
		AND lv.reqtarget IS NOT DISTINCT FROM p_reqtarget 
		AND lv.reqdescr IS NOT DISTINCT FROM p_reqdescr 
		AND lv.reqnotesfortech IS NOT DISTINCT FROM p_reqnotesfortech 
		AND lv.reqnotesforcust IS NOT DISTINCT FROM p_reqnotesforcust 
		AND lv.reqfldnotes IS NOT DISTINCT FROM p_reqfldnotes 
		AND lv.reqprogramactions IS NOT DISTINCT FROM p_reqprogramactions 
		AND lv.datetimeclosed IS NOT DISTINCT FROM p_datetimeclosed 
		AND lv.techclosed IS NOT DISTINCT FROM p_techclosed 
		AND lv.sr_number IS NOT DISTINCT FROM p_sr_number 
		AND lv.reviewed IS NOT DISTINCT FROM p_reviewed 
		AND lv.reviewedby IS NOT DISTINCT FROM p_reviewedby 
		AND lv.revieweddate IS NOT DISTINCT FROM p_revieweddate 
		AND lv.accepted IS NOT DISTINCT FROM p_accepted 
		AND lv.accepteddate IS NOT DISTINCT FROM p_accepteddate 
		AND lv.rejectedby IS NOT DISTINCT FROM p_rejectedby 
		AND lv.rejecteddate IS NOT DISTINCT FROM p_rejecteddate 
		AND lv.rejectedreason IS NOT DISTINCT FROM p_rejectedreason 
		AND lv.duedate IS NOT DISTINCT FROM p_duedate 
		AND lv.acceptedby IS NOT DISTINCT FROM p_acceptedby 
		AND lv.comments IS NOT DISTINCT FROM p_comments 
		AND lv.estcompletedate IS NOT DISTINCT FROM p_estcompletedate 
		AND lv.nextaction IS NOT DISTINCT FROM p_nextaction 
		AND lv.recordstatus IS NOT DISTINCT FROM p_recordstatus 
		AND lv.globalid IS NOT DISTINCT FROM p_globalid 
		AND lv.created_user IS NOT DISTINCT FROM p_created_user 
		AND lv.created_date IS NOT DISTINCT FROM p_created_date 
		AND lv.last_edited_user IS NOT DISTINCT FROM p_last_edited_user 
		AND lv.last_edited_date IS NOT DISTINCT FROM p_last_edited_date 
		AND lv.firstresponsedate IS NOT DISTINCT FROM p_firstresponsedate 
		AND lv.responsedaycount IS NOT DISTINCT FROM p_responsedaycount 
		AND lv.allowed IS NOT DISTINCT FROM p_allowed 
		AND lv.xvalue IS NOT DISTINCT FROM p_xvalue 
		AND lv.yvalue IS NOT DISTINCT FROM p_yvalue 
		AND lv.validx IS NOT DISTINCT FROM p_validx 
		AND lv.validy IS NOT DISTINCT FROM p_validy 
		AND lv.externalid IS NOT DISTINCT FROM p_externalid 
		AND lv.externalerror IS NOT DISTINCT FROM p_externalerror 
		AND lv.pointlocid IS NOT DISTINCT FROM p_pointlocid 
		AND lv.notified IS NOT DISTINCT FROM p_notified 
		AND lv.notifieddate IS NOT DISTINCT FROM p_notifieddate 
		AND lv.scheduled IS NOT DISTINCT FROM p_scheduled 
		AND lv.scheduleddate IS NOT DISTINCT FROM p_scheduleddate 
		AND lv.dog IS NOT DISTINCT FROM p_dog 
		AND lv.schedule_period IS NOT DISTINCT FROM p_schedule_period 
		AND lv.schedule_notes IS NOT DISTINCT FROM p_schedule_notes 
		AND lv.spanish IS NOT DISTINCT FROM p_spanish 
		AND lv.creationdate IS NOT DISTINCT FROM p_creationdate 
		AND lv.creator IS NOT DISTINCT FROM p_creator 
		AND lv.editdate IS NOT DISTINCT FROM p_editdate 
		AND lv.editor IS NOT DISTINCT FROM p_editor 
		AND lv.issuesreported IS NOT DISTINCT FROM p_issuesreported 
		AND lv.jurisdiction IS NOT DISTINCT FROM p_jurisdiction 
		AND lv.notificationtimestamp IS NOT DISTINCT FROM p_notificationtimestamp 
		AND lv.zone IS NOT DISTINCT FROM p_zone 
		AND lv.zone2 IS NOT DISTINCT FROM p_zone2 
		AND lv.geometry IS NOT DISTINCT FROM p_geometry
		AND lv.geospatial IS NOT DISTINCT FROM p_geospatial
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
		
		recdatetime,
		source,
		entrytech,
		priority,
		supervisor,
		assignedtech,
		status,
		clranon,
		clrfname,
		clrphone1,
		clrphone2,
		clremail,
		clrcompany,
		clraddr1,
		clraddr2,
		clrcity,
		clrstate,
		clrzip,
		clrother,
		clrcontpref,
		reqcompany,
		reqaddr1,
		reqaddr2,
		reqcity,
		reqstate,
		reqzip,
		reqcrossst,
		reqsubdiv,
		reqmapgrid,
		reqpermission,
		reqtarget,
		reqdescr,
		reqnotesfortech,
		reqnotesforcust,
		reqfldnotes,
		reqprogramactions,
		datetimeclosed,
		techclosed,
		sr_number,
		reviewed,
		reviewedby,
		revieweddate,
		accepted,
		accepteddate,
		rejectedby,
		rejecteddate,
		rejectedreason,
		duedate,
		acceptedby,
		comments,
		estcompletedate,
		nextaction,
		recordstatus,
		globalid,
		created_user,
		created_date,
		last_edited_user,
		last_edited_date,
		firstresponsedate,
		responsedaycount,
		allowed,
		xvalue,
		yvalue,
		validx,
		validy,
		externalid,
		externalerror,
		pointlocid,
		notified,
		notifieddate,
		scheduled,
		scheduleddate,
		dog,
		schedule_period,
		schedule_notes,
		spanish,
		creationdate,
		creator,
		editdate,
		editor,
		issuesreported,
		jurisdiction,
		notificationtimestamp,
		zone,
		zone2,
		geometry,
		geospatial,
		VERSION
	) VALUES (
		p_objectid,
		
		p_recdatetime,
		p_source,
		p_entrytech,
		p_priority,
		p_supervisor,
		p_assignedtech,
		p_status,
		p_clranon,
		p_clrfname,
		p_clrphone1,
		p_clrphone2,
		p_clremail,
		p_clrcompany,
		p_clraddr1,
		p_clraddr2,
		p_clrcity,
		p_clrstate,
		p_clrzip,
		p_clrother,
		p_clrcontpref,
		p_reqcompany,
		p_reqaddr1,
		p_reqaddr2,
		p_reqcity,
		p_reqstate,
		p_reqzip,
		p_reqcrossst,
		p_reqsubdiv,
		p_reqmapgrid,
		p_reqpermission,
		p_reqtarget,
		p_reqdescr,
		p_reqnotesfortech,
		p_reqnotesforcust,
		p_reqfldnotes,
		p_reqprogramactions,
		p_datetimeclosed,
		p_techclosed,
		p_sr_number,
		p_reviewed,
		p_reviewedby,
		p_revieweddate,
		p_accepted,
		p_accepteddate,
		p_rejectedby,
		p_rejecteddate,
		p_rejectedreason,
		p_duedate,
		p_acceptedby,
		p_comments,
		p_estcompletedate,
		p_nextaction,
		p_recordstatus,
		p_globalid,
		p_created_user,
		p_created_date,
		p_last_edited_user,
		p_last_edited_date,
		p_firstresponsedate,
		p_responsedaycount,
		p_allowed,
		p_xvalue,
		p_yvalue,
		p_validx,
		p_validy,
		p_externalid,
		p_externalerror,
		p_pointlocid,
		p_notified,
		p_notifieddate,
		p_scheduled,
		p_scheduleddate,
		p_dog,
		p_schedule_period,
		p_schedule_notes,
		p_spanish,
		p_creationdate,
		p_creator,
		p_editdate,
		p_editor,
		p_issuesreported,
		p_jurisdiction,
		p_notificationtimestamp,
		p_zone,
		p_zone2,
		p_geometry,
		p_geospatial,
		v_next_version
	);
	
	-- Return success with new version
	RETURN QUERY SELECT TRUE AS row_inserted, v_next_version AS version_num;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
