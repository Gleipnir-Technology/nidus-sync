-- +goose Up
ALTER TABLE organization ADD COLUMN is_catchall BOOLEAN;
UPDATE organization SET is_catchall = false;
ALTER TABLE organization ALTER COLUMN is_catchall SET NOT NULL;
CREATE UNIQUE INDEX only_one_catchall 
ON organization (is_catchall) 
WHERE is_catchall = true;
INSERT INTO organization(name, is_catchall) VALUES ('Gleipnir Catch-All', true);
UPDATE publicreport.nuisance SET organization_id = organization.id FROM organization WHERE organization.is_catchall = true;
ALTER TABLE publicreport.nuisance ALTER COLUMN organization_id SET NOT NULL;
UPDATE publicreport.water SET organization_id = organization.id FROM organization WHERE organization.is_catchall = true;
ALTER TABLE publicreport.water ALTER COLUMN organization_id SET NOT NULL;
-- +goose Down
ALTER TABLE organization DROP COLUMN is_catchall;
ALTER TABLE publicreport.nuisance ALTER COLUMN organization_id DROP NOT NULL;
ALTER TABLE publicreport.water ALTER COLUMN organization_id DROP NOT NULL;
