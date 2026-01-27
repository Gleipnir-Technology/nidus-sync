-- +goose Up
ALTER TABLE comms.text_log ADD COLUMN twilio_sid TEXT UNIQUE;
ALTER TABLE comms.text_log ADD COLUMN twilio_status TEXT;
UPDATE comms.text_log SET twilio_status = '';
ALTER TABLE comms.text_log ALTER COLUMN twilio_status SET NOT NULL;
ALTER TABLE comms.text_log ADD COLUMN is_visible_to_llm BOOLEAN;
UPDATE comms.text_log SET is_visible_to_llm = FALSE;
ALTER TABLE comms.text_log ALTER COLUMN is_visible_to_llm SET NOT NULL;
ALTER TYPE comms.TextOrigin ADD VALUE 'command-response';
-- +goose Down
ALTER TABLE comms.text_log DROP COLUMN is_visible_to_llm;
ALTER TABLE comms.text_log DROP COLUMN twilio_status;
ALTER TABLE comms.text_log DROP COLUMN twilio_sid;
