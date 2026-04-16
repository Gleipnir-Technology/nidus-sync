-- +goose Up
ALTER TYPE JobType ADD VALUE 'compliance-mailer-send';
