-- +goose Up
INSERT INTO communication (
	closed,
	closed_by,
	created,
	--id,
	invalidated,
	invalidated_by,
	opened,
	opened_by,
	organization_id,
	response_email_log_id,
	response_text_log_id,
	set_pending,
	set_pending_by,
	source_email_log_id,
	source_report_id,
	source_text_log_id
) SELECT 
	NULL,
	NULL,
	created,
	NULL,
	NULL,
	NULL,
	NULL,
	organization_id,
	NULL,
	NULL,
	NULL,
	NULL,
	NULL,
	id,
	NULL
FROM publicreport.report;
-- +goose Down
DELETE FROM communication;
