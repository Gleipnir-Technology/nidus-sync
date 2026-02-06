-- +goose Up
ALTER TABLE publicreport.nuisance ADD COLUMN reporter_contact_consent BOOLEAN;
ALTER TABLE publicreport.pool ADD COLUMN reporter_contact_consent BOOLEAN;

-- +goose Down
ALTER TABLE publicreport.pool DROP COLUMN reporter_contact_consent;
ALTER TABLE publicreport.nuisance DROP COLUMN reporter_contact_consent;
