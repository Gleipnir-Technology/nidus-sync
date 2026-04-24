-- +goose Up
ALTER TABLE publicreport.water ADD COLUMN duration publicreport.NuisanceDurationType;
UPDATE publicreport.water SET duration = 'none';
ALTER TABLE publicreport.water ALTER COLUMN duration SET NOT NULL;
-- +goose Down
ALTER TABLE publicreport.water DROP COLUMN duration;
