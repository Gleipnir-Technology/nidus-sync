-- +goose Up
ALTER TABLE comms.email_log ALTER COLUMN template_id SET NOT NULL;
