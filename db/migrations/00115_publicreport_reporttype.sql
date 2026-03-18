-- +goose Up
ALTER TABLE publicreport.report DROP CONSTRAINT report_report_type_check;
CREATE TYPE publicreport.ReportType AS ENUM('nuisance', 'water');
ALTER TABLE publicreport.report ALTER COLUMN report_type TYPE publicreport.ReportType USING report_type::publicreport.ReportType;
-- +goose Down

