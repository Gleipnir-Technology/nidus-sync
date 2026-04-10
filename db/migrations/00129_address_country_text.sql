-- +goose Up
ALTER TABLE address 
ALTER COLUMN country 
TYPE TEXT 
USING country::TEXT;
DROP TYPE CountryType;
-- +goose Down

