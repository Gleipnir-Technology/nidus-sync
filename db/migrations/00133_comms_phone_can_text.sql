-- +goose Up
ALTER TABLE comms.phone ADD COLUMN can_sms BOOLEAN;
UPDATE comms.phone SET can_sms = TRUE;
ALTER TABLE comms.phone ALTER COLUMN can_sms SET NOT NULL;
ALTER TABLE publicreport.report ADD COLUMN reporter_phone_can_sms BOOLEAN;
UPDATE publicreport.report SET reporter_phone_can_sms = TRUE;
ALTER TABLE publicreport.report ALTER COLUMN reporter_phone_can_sms SET NOT NULL;
-- +goose Down
ALTER TABLE comms.phone DROP COLUMN can_sms;
ALTER TABLE publicreport.report DROAP COLUMN reporter_phone_can_sms;
