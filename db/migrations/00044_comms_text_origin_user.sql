-- +goose Up
ALTER TYPE comms.TextOrigin ADD VALUE 'customer';
ALTER TYPE comms.TextOrigin ADD VALUE 'reiteration';
-- +goose Down
ALTER TYPE comms.TextOrigin DROP VALUE 'reiteration';
ALTER TYPE comms.TextOrigin DROP VALUE 'customer';
