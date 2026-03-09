-- +goose Up
ALTER TABLE publicreport.nuisance RENAME COLUMN address TO address_raw;
ALTER TABLE publicreport.nuisance RENAME COLUMN address_place TO address_locality;
ALTER TABLE publicreport.nuisance RENAME COLUMN address_postcode TO address_postal_code;

ALTER TABLE publicreport.nuisance ADD COLUMN address_id INTEGER REFERENCES address(id);

ALTER TABLE publicreport.pool RENAME COLUMN address TO address_raw;
ALTER TABLE publicreport.pool RENAME COLUMN address_place TO address_locality;
ALTER TABLE publicreport.pool RENAME COLUMN address_post_code TO address_postal_code;

ALTER TABLE publicreport.pool ADD COLUMN address_id INTEGER REFERENCES address(id);

ALTER TABLE publicreport.pool RENAME TO water;

ALTER TABLE publicreport.pool_image RENAME COLUMN pool_id TO water_id;
ALTER TABLE publicreport.pool_image RENAME TO water_image;

ALTER TABLE publicreport.notify_email_pool RENAME COLUMN pool_id TO water_id;
ALTER TABLE publicreport.notify_email_pool RENAME TO notify_email_water;

ALTER TABLE publicreport.notify_phone_pool RENAME COLUMN pool_id TO water_id;
ALTER TABLE publicreport.notify_phone_pool RENAME TO notify_phone_water;
-- +goose Down
ALTER TABLE publicreport.notify_phone_water RENAME TO notify_phone_pool;
ALTER TABLE publicreport.notify_phone_pool RENAME COLUMN water_id TO pool_id;

ALTER TABLE publicreport.notify_email_water RENAME TO notify_email_pool;
ALTER TABLE publicreport.notify_email_pool RENAME COLUMN water_id TO pool_id;

ALTER TABLE publicreport.water_image RENAME COLUMN water_id TO pool_id;
ALTER TABLE publicreport.water_image RENAME TO pool_image;

ALTER TABLE publicreport.water RENAME TO pool;

ALTER TABLE publicreport.pool DROP COLUMN address_id;

ALTER TABLE publicreport.pool RENAME COLUMN address_postal_code TO address_post_code;
ALTER TABLE publicreport.pool RENAME COLUMN address_locality TO address_place;
ALTER TABLE publicreport.pool RENAME COLUMN address_raw TO address;

ALTER TABLE publicreport.nuisance DROP COLUMN address_id;

ALTER TABLE publicreport.nuisance RENAME COLUMN address_postal_code TO address_postcode;
ALTER TABLE publicreport.nuisance RENAME COLUMN address_locality TO address_place;
ALTER TABLE publicreport.nuisance RENAME COLUMN address_raw TO address;
