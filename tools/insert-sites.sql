BEGIN;
INSERT INTO site(address_id, created, creator_id, file_id, id, notes, organization_id, owner_name, owner_phone_e164, parcel_id, resident_owned, tags, version)
VALUES (:address_id, NOW(), :user_id, NULL, DEFAULT, '', :organization_id, '', NULL, :parcel_id, NULL, '', 1);
COMMIT;
