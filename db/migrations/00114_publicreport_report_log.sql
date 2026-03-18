-- +goose UP
CREATE TYPE publicreport.ReportLogType AS ENUM (
	'created',
	'invalidated',
	'message-email',
	'message-text',
	'reviewed',
	'scheduled',
	'treated'
);
CREATE TABLE publicreport.report_log (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	email_log_id INTEGER REFERENCES comms.email_log(id),
	id SERIAL NOT NULL,
	report_id INTEGER NOT NULL REFERENCES publicreport.report(id),
	text_log_id INTEGER REFERENCES comms.text_log(id),
	type_ publicreport.ReportLogType NOT NULL,
	user_id INTEGER REFERENCES user_(id),
	PRIMARY KEY (id)
);
INSERT INTO publicreport.report_log (
	created,
	email_log_id,
	report_id,
	text_log_id,
	type_,
	user_id
)
SELECT
	created,
	NULL,
	id,
	NULL,
	'created',
	NULL
FROM publicreport.report;
-- +goose Down
DROP TABLE publicreport.report_log;
DROP TYPE publicreport.ReportLogType;
