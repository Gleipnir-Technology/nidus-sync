-- +goose Up
ALTER TABLE publicreport.nuisance ADD COLUMN address_number TEXT;
UPDATE publicreport.nuisance SET address_number = '';
ALTER TABLE publicreport.nuisance ALTER COLUMN address_number SET NOT NULL;

ALTER TABLE publicreport.pool ADD COLUMN address_number TEXT;
UPDATE publicreport.pool SET address_number = '';
ALTER TABLE publicreport.pool ALTER COLUMN address_number SET NOT NULL;
-- +goose Down
ALTER TABLE publicreport.pool DROP COLUMN address_number;
ALTER TABLE publicreport.nuisance DROP COLUMN address_number;

