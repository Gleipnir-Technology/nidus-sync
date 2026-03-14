-- +goose Up
ALTER TYPE LeadType ADD VALUE 'unknown' BEFORE 'green-pool';
ALTER TYPE LeadType ADD VALUE 'publicreport-nuisance' AFTER 'green-pool';
ALTER TYPE LeadType ADD VALUE 'publicreport-water' AFTER 'publicreport-nuisance';
