-- +goose Up
-- Step 1: Drop foreign key constraints that reference the composite primary key
ALTER TABLE feature DROP CONSTRAINT feature_site_id_site_version_fkey;
ALTER TABLE lead DROP CONSTRAINT lead_site_id_site_version_fkey;
ALTER TABLE resident DROP CONSTRAINT resident_site_id_site_version_fkey;

-- Step 2: Drop the existing composite primary key
ALTER TABLE site DROP CONSTRAINT site_pkey;

-- Step 3: Create new primary key on just the id column
ALTER TABLE site ADD PRIMARY KEY (id);

-- Step 4: Create a unique constraint on (id, version) to maintain uniqueness
ALTER TABLE site ADD CONSTRAINT site_id_version_unique UNIQUE (id, version);

-- Step 5: Recreate the foreign key constraints
ALTER TABLE feature 
  ADD CONSTRAINT feature_site_id_fkey 
  FOREIGN KEY (site_id)
  REFERENCES site(id);
ALTER TABLE feature
  DROP COLUMN site_version;

ALTER TABLE lead 
  ADD CONSTRAINT lead_site_id_fkey 
  FOREIGN KEY (site_id)
  REFERENCES site(id);
ALTER TABLE lead
  DROP COLUMN site_version;

ALTER TABLE resident 
  ADD CONSTRAINT resident_site_id_fkey 
  FOREIGN KEY (site_id) 
  REFERENCES site(id);
ALTER TABLE resident
  DROP COLUMN site_version;
