-- +goose Up
CREATE TABLE communication (
	closed TIMESTAMP WITHOUT TIME ZONE,
	closed_by INTEGER REFERENCES user_(id),
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	id SERIAL,
	invalidated TIMESTAMP WITHOUT TIME ZONE,
	invalidated_by INTEGER REFERENCES user_(id),
	opened TIMESTAMP WITHOUT TIME ZONE,
	opened_by INTEGER REFERENCES user_(id),
	response_email_log_id INTEGER REFERENCES comms.email_log(id),
	response_text_log_id INTEGER REFERENCES comms.text_log(id),
	set_pending TIMESTAMP WITHOUT TIME ZONE,
	set_pending_by INTEGER REFERENCES user_(id),
	source_email_log_id INTEGER REFERENCES comms.email_log(id),
	source_report_id INTEGER REFERENCES publicreport.report(id),
	source_text_log_id INTEGER REFERENCES comms.text_log(id),
	PRIMARY KEY(id)
);
-- +goose Down
DROP TABLE communication;
