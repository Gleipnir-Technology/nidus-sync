-- +goose Up
UPDATE address 
SET gid = gen_random_uuid() 
WHERE gid = '';
ALTER TABLE address ADD CONSTRAINT address_gid_unique UNIQUE (gid);
-- +goose Down
ALTER TABLE address DROP CONSTRAINT address_gid_unique;
