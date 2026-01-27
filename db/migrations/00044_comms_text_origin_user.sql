-- +goose Up
ALTER TYPE comms.TextOrigin ADD VALUE 'customer';
ALTER TYPE comms.TextOrigin ADD VALUE 'reiteration';
