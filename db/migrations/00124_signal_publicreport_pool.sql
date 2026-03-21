-- +goose Up
ALTER TABLE signal ADD COLUMN feature_pool_feature_id INTEGER REFERENCES feature_pool(feature_id);
ALTER TABLE signal ADD COLUMN report_id INTEGER REFERENCES publicreport.report(id);
ALTER TABLE signal 
ADD CONSTRAINT check_exclusive_reference 
CHECK (
    (feature_pool_feature_id IS NULL OR report_id IS NULL)
);
-- +goose Down
ALTER TABLE signal DROP CONSTRAINT check_exclusive_reference;
ALTER TABLE signal DROP COLUMN report_id;
ALTER TABLE signal DROP COLUMN feature_pool_feature_id;
