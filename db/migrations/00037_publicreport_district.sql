-- +goose Up
ALTER TABLE publicreport.quick ADD COLUMN organization_id INTEGER REFERENCES "public"."organization"(id);
ALTER TABLE publicreport.pool ADD COLUMN organization_id INTEGER REFERENCES "public"."organization"(id);
ALTER TABLE publicreport.nuisance ADD COLUMN organization_id INTEGER REFERENCES "public"."organization"(id);
-- +goose Down
ALTER TABLE publicreport.nuisance DROP COLUMN organization_id;
ALTER TABLE publicreport.pool DROP COLUMN organization_id;
ALTER TABLE publicreport.quick DROP COLUMN organization_id;
