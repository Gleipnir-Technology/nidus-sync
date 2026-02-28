-- +goose Up
ALTER TABLE organization ADD COLUMN mailing_address_country TEXT;
ALTER TABLE organization ADD COLUMN mailing_address_state TEXT;
ALTER TABLE organization ADD COLUMN office_address_country TEXT;
ALTER TABLE organization ADD COLUMN office_address_state TEXT;

UPDATE organization SET mailing_address_country = 'USA';
UPDATE organization SET mailing_address_state = 'CA';
UPDATE organization SET office_address_country = 'USA';
UPDATE organization SET office_address_state = 'CA';

UPDATE organization
SET office_address_postal_code = split_part(office_address_postal_code, '.', 1)
WHERE office_address_postal_code LIKE '%.%';

UPDATE organization
SET mailing_address_postal_code = split_part(mailing_address_postal_code, '.', 1)
WHERE mailing_address_postal_code LIKE '%.%';
-- +goose Down
ALTER TABLE organization DROP COLUMN mailing_address_country;
ALTER TABLE organization DROP COLUMN mailing_address_state;
ALTER TABLE organization DROP COLUMN office_address_country;
ALTER TABLE organization DROP COLUMN office_address_state;
