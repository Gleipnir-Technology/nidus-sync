-- +goose Up
ALTER TABLE site ALTER COLUMN parcel_id DROP NOT NULL;
-- +goose Down
ALTER TABLE site ALTER COLUMN parcel_id ADD NOT NULL;
