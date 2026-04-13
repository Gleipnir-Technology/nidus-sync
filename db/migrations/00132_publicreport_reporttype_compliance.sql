-- +goose Up
ALTER TYPE publicreport.ReportType ADD VALUE 'compliance' AFTER 'water';
