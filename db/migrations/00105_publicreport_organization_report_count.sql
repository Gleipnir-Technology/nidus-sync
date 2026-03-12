-- +goose Up
ALTER TABLE publicreport.nuisance ADD COLUMN reviewed TIMESTAMP WITHOUT TIME ZONE;
ALTER TABLE publicreport.nuisance ADD COLUMN reviewer_id INTEGER REFERENCES user_(id);
ALTER TABLE publicreport.water ADD COLUMN reviewed TIMESTAMP WITHOUT TIME ZONE;
ALTER TABLE publicreport.water ADD COLUMN reviewer_id INTEGER REFERENCES user_(id);

CREATE VIEW publicreport.organization_report_count AS
SELECT 
    o.id AS organization_id,
    COUNT(n.id) FILTER (WHERE n.status = 'reported') AS nuisance_reported,
    COUNT(n.id) FILTER (WHERE n.status = 'reviewed') AS nuisance_reviewed,
    COUNT(n.id) FILTER (WHERE n.status = 'scheduled') AS nuisance_scheduled,
    COUNT(n.id) FILTER (WHERE n.status = 'treated') AS nuisance_treated,
    COUNT(n.id) FILTER (WHERE n.status = 'invalidated') AS nuisance_invalidated,
    COUNT(w.id) FILTER (WHERE w.status = 'reported') AS water_reported,
    COUNT(w.id) FILTER (WHERE w.status = 'reviewed') AS water_reviewed,
    COUNT(w.id) FILTER (WHERE w.status = 'scheduled') AS water_scheduled,
    COUNT(w.id) FILTER (WHERE w.status = 'treated') AS water_treated,
    COUNT(w.id) FILTER (WHERE w.status = 'invalidated') AS water_invalidated
FROM organization o
LEFT JOIN publicreport.nuisance n ON o.id = n.organization_id
LEFT JOIN publicreport.water w ON o.id = w.organization_id
GROUP BY o.id;

-- +goose Down
DROP VIEW publicreport.organization_report_count;
ALTER TABLE publicreport.water DROP COLUMN reviewer_id;
ALTER TABLE publicreport.water DROP COLUMN reviewed;
ALTER TABLE publicreport.nuisance DROP COLUMN reviewer_id;
ALTER TABLE publicreport.nuisance DROP COLUMN reviewed;
