-- +goose Up
ALTER TYPE publicreport.ReportStatusType ADD VALUE 'invalidated' AFTER 'treated';
-- +goose Down
