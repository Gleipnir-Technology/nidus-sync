-- +goose Up
ALTER TYPE comms.MessageTypeEmail ADD VALUE 'report-notification-confirmation' AFTER 'report-status-complete';
UPDATE comms.email_log SET source = 'report-notification-confirmation' WHERE source = 'report-subscription-confirmation';
