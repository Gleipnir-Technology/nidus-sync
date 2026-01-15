-- +goose Up
ALTER TYPE H3AggregationType ADD VALUE 'Trap' AFTER 'ServiceRequest';

-- +goose Down
ALTER TYPE H3AggregationType DROP VALUE 'Trap';
