-- Prepared statement for conditional insert with versioning for fieldseeker.servicerequest
-- Only inserts a new version if data has changed

PREPARE insert_servicerequest_versioned(bigint, timestamp, fieldseeker.servicerequest_servicerequestsource_enum, varchar, fieldseeker.servicerequest_servicerequestpriority_enum, fieldseeker.servicerequest_servicerequest_supervisor_eba07b90_c885_4fe6_8080_7aa775403b9a_enum, fieldseeker.servicerequest_servicerequest_assignedtech_71d0d685_868f_4b7a_87e2_3661a3ee67c5_enum, fieldseeker.servicerequest_servicerequeststatus_enum, fieldseeker.servicerequest_notinuit_f_enum, varchar, varchar, varchar, varchar, varchar, varchar, varchar, varchar, fieldseeker.servicerequest_servicerequestregion_enum, varchar, varchar, fieldseeker.servicerequest_servicerequestcontactpreferences_enum, varchar, varchar, varchar, varchar, fieldseeker.servicerequest_servicerequestregion_enum, varchar, varchar, varchar, varchar, fieldseeker.servicerequest_notinuit_f_enum, fieldseeker.servicerequest_servicerequesttarget_enum, varchar, varchar, varchar, varchar, varchar, timestamp, varchar, integer, fieldseeker.servicerequest_notinuit_f_enum, varchar, timestamp, fieldseeker.servicerequest_notinuit_f_enum, timestamp, varchar, timestamp, fieldseeker.servicerequest_servicerequestrejectedreason_enum, timestamp, varchar, varchar, timestamp, fieldseeker.servicerequest_servicerequestnextaction_enum, smallint, uuid, varchar, timestamp, varchar, timestamp, timestamp, smallint, varchar, varchar, varchar, varchar, varchar, varchar, varchar, uuid, smallint, timestamp, smallint, timestamp, fieldseeker.servicerequest_servicerequest_dog_2b95ec97_1286_4fcd_88f4_f0e31113f696_enum, fieldseeker.servicerequest_servicerequest_schedule_period_3f40c046_afd1_4abd_8bf4_389650d29a49_enum, varchar, fieldseeker.servicerequest_servicerequest_spanish_aaa3dc66_9f9a_4527_8ecd_c9f76db33879_enum, timestamp, varchar, timestamp, varchar, fieldseeker.servicerequest_servicerequestissues_enum, varchar, varchar, varchar, varchar) AS
WITH
-- Get the current latest version of this record
latest_version AS (
  SELECT * FROM fieldseeker.servicerequest
  WHERE objectid = $1
  ORDER BY VERSION DESC
  LIMIT 1
),
-- Calculate the next version number
next_version AS (
  SELECT COALESCE(MAX(VERSION) + 1, 1) as version_num
  FROM fieldseeker.servicerequest
  WHERE objectid = $1
)
-- Perform conditional insert
INSERT INTO fieldseeker.servicerequest (
  objectid, recdatetime, source, entrytech, priority, supervisor, assignedtech, status, clranon, clrfname, clrphone1, clrphone2, clremail, clrcompany, clraddr1, clraddr2, clrcity, clrstate, clrzip, clrother, clrcontpref, reqcompany, reqaddr1, reqaddr2, reqcity, reqstate, reqzip, reqcrossst, reqsubdiv, reqmapgrid, reqpermission, reqtarget, reqdescr, reqnotesfortech, reqnotesforcust, reqfldnotes, reqprogramactions, datetimeclosed, techclosed, sr_number, reviewed, reviewedby, revieweddate, accepted, accepteddate, rejectedby, rejecteddate, rejectedreason, duedate, acceptedby, comments, estcompletedate, nextaction, recordstatus, globalid, created_user, created_date, last_edited_user, last_edited_date, firstresponsedate, responsedaycount, allowed, xvalue, yvalue, validx, validy, externalid, externalerror, pointlocid, notified, notifieddate, scheduled, scheduleddate, dog, schedule_period, schedule_notes, spanish, creationdate, creator, editdate, editor, issuesreported, jurisdiction, notificationtimestamp, zone, zone2,
  VERSION
)
SELECT
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86,
  v.version_num
FROM next_version v
WHERE
  -- Only insert if no record exists yet OR data has changed
  NOT EXISTS (SELECT 1 FROM latest_version lv WHERE
    lv.objectid IS NOT DISTINCT FROM $1 AND
    lv.recdatetime IS NOT DISTINCT FROM $2 AND
    lv.source IS NOT DISTINCT FROM $3 AND
    lv.entrytech IS NOT DISTINCT FROM $4 AND
    lv.priority IS NOT DISTINCT FROM $5 AND
    lv.supervisor IS NOT DISTINCT FROM $6 AND
    lv.assignedtech IS NOT DISTINCT FROM $7 AND
    lv.status IS NOT DISTINCT FROM $8 AND
    lv.clranon IS NOT DISTINCT FROM $9 AND
    lv.clrfname IS NOT DISTINCT FROM $10 AND
    lv.clrphone1 IS NOT DISTINCT FROM $11 AND
    lv.clrphone2 IS NOT DISTINCT FROM $12 AND
    lv.clremail IS NOT DISTINCT FROM $13 AND
    lv.clrcompany IS NOT DISTINCT FROM $14 AND
    lv.clraddr1 IS NOT DISTINCT FROM $15 AND
    lv.clraddr2 IS NOT DISTINCT FROM $16 AND
    lv.clrcity IS NOT DISTINCT FROM $17 AND
    lv.clrstate IS NOT DISTINCT FROM $18 AND
    lv.clrzip IS NOT DISTINCT FROM $19 AND
    lv.clrother IS NOT DISTINCT FROM $20 AND
    lv.clrcontpref IS NOT DISTINCT FROM $21 AND
    lv.reqcompany IS NOT DISTINCT FROM $22 AND
    lv.reqaddr1 IS NOT DISTINCT FROM $23 AND
    lv.reqaddr2 IS NOT DISTINCT FROM $24 AND
    lv.reqcity IS NOT DISTINCT FROM $25 AND
    lv.reqstate IS NOT DISTINCT FROM $26 AND
    lv.reqzip IS NOT DISTINCT FROM $27 AND
    lv.reqcrossst IS NOT DISTINCT FROM $28 AND
    lv.reqsubdiv IS NOT DISTINCT FROM $29 AND
    lv.reqmapgrid IS NOT DISTINCT FROM $30 AND
    lv.reqpermission IS NOT DISTINCT FROM $31 AND
    lv.reqtarget IS NOT DISTINCT FROM $32 AND
    lv.reqdescr IS NOT DISTINCT FROM $33 AND
    lv.reqnotesfortech IS NOT DISTINCT FROM $34 AND
    lv.reqnotesforcust IS NOT DISTINCT FROM $35 AND
    lv.reqfldnotes IS NOT DISTINCT FROM $36 AND
    lv.reqprogramactions IS NOT DISTINCT FROM $37 AND
    lv.datetimeclosed IS NOT DISTINCT FROM $38 AND
    lv.techclosed IS NOT DISTINCT FROM $39 AND
    lv.sr_number IS NOT DISTINCT FROM $40 AND
    lv.reviewed IS NOT DISTINCT FROM $41 AND
    lv.reviewedby IS NOT DISTINCT FROM $42 AND
    lv.revieweddate IS NOT DISTINCT FROM $43 AND
    lv.accepted IS NOT DISTINCT FROM $44 AND
    lv.accepteddate IS NOT DISTINCT FROM $45 AND
    lv.rejectedby IS NOT DISTINCT FROM $46 AND
    lv.rejecteddate IS NOT DISTINCT FROM $47 AND
    lv.rejectedreason IS NOT DISTINCT FROM $48 AND
    lv.duedate IS NOT DISTINCT FROM $49 AND
    lv.acceptedby IS NOT DISTINCT FROM $50 AND
    lv.comments IS NOT DISTINCT FROM $51 AND
    lv.estcompletedate IS NOT DISTINCT FROM $52 AND
    lv.nextaction IS NOT DISTINCT FROM $53 AND
    lv.recordstatus IS NOT DISTINCT FROM $54 AND
    lv.globalid IS NOT DISTINCT FROM $55 AND
    lv.created_user IS NOT DISTINCT FROM $56 AND
    lv.created_date IS NOT DISTINCT FROM $57 AND
    lv.last_edited_user IS NOT DISTINCT FROM $58 AND
    lv.last_edited_date IS NOT DISTINCT FROM $59 AND
    lv.firstresponsedate IS NOT DISTINCT FROM $60 AND
    lv.responsedaycount IS NOT DISTINCT FROM $61 AND
    lv.allowed IS NOT DISTINCT FROM $62 AND
    lv.xvalue IS NOT DISTINCT FROM $63 AND
    lv.yvalue IS NOT DISTINCT FROM $64 AND
    lv.validx IS NOT DISTINCT FROM $65 AND
    lv.validy IS NOT DISTINCT FROM $66 AND
    lv.externalid IS NOT DISTINCT FROM $67 AND
    lv.externalerror IS NOT DISTINCT FROM $68 AND
    lv.pointlocid IS NOT DISTINCT FROM $69 AND
    lv.notified IS NOT DISTINCT FROM $70 AND
    lv.notifieddate IS NOT DISTINCT FROM $71 AND
    lv.scheduled IS NOT DISTINCT FROM $72 AND
    lv.scheduleddate IS NOT DISTINCT FROM $73 AND
    lv.dog IS NOT DISTINCT FROM $74 AND
    lv.schedule_period IS NOT DISTINCT FROM $75 AND
    lv.schedule_notes IS NOT DISTINCT FROM $76 AND
    lv.spanish IS NOT DISTINCT FROM $77 AND
    lv.creationdate IS NOT DISTINCT FROM $78 AND
    lv.creator IS NOT DISTINCT FROM $79 AND
    lv.editdate IS NOT DISTINCT FROM $80 AND
    lv.editor IS NOT DISTINCT FROM $81 AND
    lv.issuesreported IS NOT DISTINCT FROM $82 AND
    lv.jurisdiction IS NOT DISTINCT FROM $83 AND
    lv.notificationtimestamp IS NOT DISTINCT FROM $84 AND
    lv.zone IS NOT DISTINCT FROM $85 AND
    lv.zone2 IS NOT DISTINCT FROM $86
  )
RETURNING *;

-- Example usage: EXECUTE insert_servicerequest_versioned(id, value1, value2, ...);

-- Parameters in order:
-- $1: OBJECTID (bigint)
-- $2: RECDATETIME (timestamp)
-- $3: SOURCE (fieldseeker.servicerequest_servicerequestsource_enum)
-- $4: ENTRYTECH (varchar)
-- $5: PRIORITY (fieldseeker.servicerequest_servicerequestpriority_enum)
-- $6: SUPERVISOR (fieldseeker.servicerequest_servicerequest_supervisor_eba07b90_c885_4fe6_8080_7aa775403b9a_enum)
-- $7: ASSIGNEDTECH (fieldseeker.servicerequest_servicerequest_assignedtech_71d0d685_868f_4b7a_87e2_3661a3ee67c5_enum)
-- $8: STATUS (fieldseeker.servicerequest_servicerequeststatus_enum)
-- $9: CLRANON (fieldseeker.servicerequest_notinuit_f_enum)
-- $10: CLRFNAME (varchar)
-- $11: CLRPHONE1 (varchar)
-- $12: CLRPHONE2 (varchar)
-- $13: CLREMAIL (varchar)
-- $14: CLRCOMPANY (varchar)
-- $15: CLRADDR1 (varchar)
-- $16: CLRADDR2 (varchar)
-- $17: CLRCITY (varchar)
-- $18: CLRSTATE (fieldseeker.servicerequest_servicerequestregion_enum)
-- $19: CLRZIP (varchar)
-- $20: CLROTHER (varchar)
-- $21: CLRCONTPREF (fieldseeker.servicerequest_servicerequestcontactpreferences_enum)
-- $22: REQCOMPANY (varchar)
-- $23: REQADDR1 (varchar)
-- $24: REQADDR2 (varchar)
-- $25: REQCITY (varchar)
-- $26: REQSTATE (fieldseeker.servicerequest_servicerequestregion_enum)
-- $27: REQZIP (varchar)
-- $28: REQCROSSST (varchar)
-- $29: REQSUBDIV (varchar)
-- $30: REQMAPGRID (varchar)
-- $31: REQPERMISSION (fieldseeker.servicerequest_notinuit_f_enum)
-- $32: REQTARGET (fieldseeker.servicerequest_servicerequesttarget_enum)
-- $33: REQDESCR (varchar)
-- $34: REQNOTESFORTECH (varchar)
-- $35: REQNOTESFORCUST (varchar)
-- $36: REQFLDNOTES (varchar)
-- $37: REQPROGRAMACTIONS (varchar)
-- $38: DATETIMECLOSED (timestamp)
-- $39: TECHCLOSED (varchar)
-- $40: SR_NUMBER (integer)
-- $41: REVIEWED (fieldseeker.servicerequest_notinuit_f_enum)
-- $42: REVIEWEDBY (varchar)
-- $43: REVIEWEDDATE (timestamp)
-- $44: ACCEPTED (fieldseeker.servicerequest_notinuit_f_enum)
-- $45: ACCEPTEDDATE (timestamp)
-- $46: REJECTEDBY (varchar)
-- $47: REJECTEDDATE (timestamp)
-- $48: REJECTEDREASON (fieldseeker.servicerequest_servicerequestrejectedreason_enum)
-- $49: DUEDATE (timestamp)
-- $50: ACCEPTEDBY (varchar)
-- $51: COMMENTS (varchar)
-- $52: ESTCOMPLETEDATE (timestamp)
-- $53: NEXTACTION (fieldseeker.servicerequest_servicerequestnextaction_enum)
-- $54: RECORDSTATUS (smallint)
-- $55: GlobalID (uuid)
-- $56: created_user (varchar)
-- $57: created_date (timestamp)
-- $58: last_edited_user (varchar)
-- $59: last_edited_date (timestamp)
-- $60: FIRSTRESPONSEDATE (timestamp)
-- $61: RESPONSEDAYCOUNT (smallint)
-- $62: ALLOWED (varchar)
-- $63: XVALUE (varchar)
-- $64: YVALUE (varchar)
-- $65: VALIDX (varchar)
-- $66: VALIDY (varchar)
-- $67: EXTERNALID (varchar)
-- $68: EXTERNALERROR (varchar)
-- $69: POINTLOCID (uuid)
-- $70: NOTIFIED (smallint)
-- $71: NOTIFIEDDATE (timestamp)
-- $72: SCHEDULED (smallint)
-- $73: SCHEDULEDDATE (timestamp)
-- $74: DOG (fieldseeker.servicerequest_servicerequest_dog_2b95ec97_1286_4fcd_88f4_f0e31113f696_enum)
-- $75: schedule_period (fieldseeker.servicerequest_servicerequest_schedule_period_3f40c046_afd1_4abd_8bf4_389650d29a49_enum)
-- $76: schedule_notes (varchar)
-- $77: Spanish (fieldseeker.servicerequest_servicerequest_spanish_aaa3dc66_9f9a_4527_8ecd_c9f76db33879_enum)
-- $78: CreationDate (timestamp)
-- $79: Creator (varchar)
-- $80: EditDate (timestamp)
-- $81: Editor (varchar)
-- $82: ISSUESREPORTED (fieldseeker.servicerequest_servicerequestissues_enum)
-- $83: JURISDICTION (varchar)
-- $84: NOTIFICATIONTIMESTAMP (varchar)
-- $85: ZONE (varchar)
-- $86: ZONE2 (varchar)
