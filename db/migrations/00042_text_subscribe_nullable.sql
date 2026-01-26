-- +goose Up
ALTER TABLE comms.phone ALTER COLUMN is_subscribed DROP NOT NULL;
