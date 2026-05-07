-- +goose Up
INSERT INTO communication (
	created,
	--id,
	organization_id,
	response_email_log_id,
	response_text_log_id,
	source_email_log_id,
	source_report_id,
	source_text_log_id,
	status
) SELECT 
	created,
	organization_id,
	NULL,
	NULL,
	NULL,
	id,
	NULL,
	'new'
FROM publicreport.report;
INSERT INTO communication_log_entry (
	communication_id,
	created,
	--id,
	type_,
	user_
) SELECT 
	id,
	created,
	'created',
	NULL
FROM communication;
-- +goose Down
DELETE FROM communication_log_entry;
DELETE FROM communication;
