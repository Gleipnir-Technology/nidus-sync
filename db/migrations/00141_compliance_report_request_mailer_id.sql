-- +goose Up
ALTER TABLE compliance_report_request_mailer
	ADD COLUMN id SERIAL PRIMARY KEY;
-- +goose Down
ALTER TABLE compliance_report_request_mailer
	DROP COLUMN id;

