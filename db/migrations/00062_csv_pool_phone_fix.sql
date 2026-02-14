-- +goose Up
ALTER TABLE fileupload.pool DROP COLUMN property_owner_phone;
ALTER TABLE fileupload.pool DROP COLUMN resident_phone;
ALTER TABLE fileupload.pool ADD COLUMN property_owner_phone_e164 TEXT REFERENCES comms.phone(e164);
ALTER TABLE fileupload.pool ADD COLUMN resident_phone_e164 TEXT REFERENCES comms.phone(e164);
