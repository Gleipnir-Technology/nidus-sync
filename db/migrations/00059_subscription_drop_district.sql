-- +goose Up
ALTER TABLE publicreport.subscribe_email DROP COLUMN district_id;
ALTER TABLE publicreport.subscribe_phone DROP COLUMN district_id;
