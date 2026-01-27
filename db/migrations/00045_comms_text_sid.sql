-- +goose Up
ALTER TABLE comms.text_log ADD COLUMN twilio_sid TEXT UNIQUE;
ALTER TABLE comms.text_log ADD COLUMN twilio_status TEXT;
UPDATE comms.text_log SET twilio_status = '';
ALTER TABLE comms.text_log ALTER COLUMN twilio_status SET NOT NULL;
-- +goose Down
ALTER TABLE comms.text_log DROP COLUMN twilio_status;
ALTER TABLE comms.text_log DROP COLUMN twilio_sid;
