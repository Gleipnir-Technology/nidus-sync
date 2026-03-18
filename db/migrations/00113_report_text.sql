-- +goose Up
ALTER TABLE comms.text_job ADD COLUMN creator_id INTEGER REFERENCES user_(id);
ALTER TABLE comms.text_job ADD COLUMN report_id INTEGER REFERENCES publicreport.report(id);

CREATE TABLE report_text (
	creator_id INTEGER NOT NULL REFERENCES user_(id),
	report_id INTEGER NOT NULL REFERENCES publicreport.report(id),
	text_log_id INTEGER NOT NULL REFERENCES comms.text_log(id),
	PRIMARY KEY(creator_id, report_id, text_log_id)
);
-- +goose Down
DROP TABLE report_text;
ALTER TABLE comms.text_job DROP COLUMN report_id;
ALTER TABLE comms.text_job DROP COLUMN creator_id;

