-- +goose Up
ALTER TYPE publicreport.AccuracyType ADD VALUE 'centroid' AFTER 'intersection';
