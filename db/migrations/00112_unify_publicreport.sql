-- +goose Up
CREATE TABLE publicreport.report (
	-- Address fields
	address_raw TEXT NOT NULL,
	address_number TEXT NOT NULL,
	address_street TEXT NOT NULL,
	address_locality TEXT NOT NULL,
	address_region TEXT NOT NULL,
	address_postal_code TEXT NOT NULL,
	address_country TEXT NOT NULL,
	address_id INTEGER REFERENCES address(id),

	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	-- Location fields
	location GEOMETRY(Point, 4326),
	h3cell h3index,
	id SERIAL,
	latlng_accuracy_type publicreport.accuracytype NOT NULL,
	latlng_accuracy_value REAL NOT NULL,
	map_zoom REAL NOT NULL,

	-- Organization
	organization_id INTEGER NOT NULL REFERENCES organization(id),

	public_id TEXT NOT NULL UNIQUE,

	-- Reporter information
	reporter_name TEXT NOT NULL,
	reporter_email TEXT NOT NULL,
	reporter_phone TEXT NOT NULL,
	reporter_contact_consent BOOLEAN,

	-- Report type discriminator
	report_type TEXT NOT NULL CHECK (report_type IN ('nuisance', 'water')),

	-- Review fields
	reviewed TIMESTAMP WITHOUT TIME ZONE,
	reviewer_id INTEGER REFERENCES user_(id),

	status publicreport.reportstatustype NOT NULL,
	PRIMARY KEY(id)
);

-- Insert nuisance reports into the report table
INSERT INTO publicreport.report (
	--id,
	public_id,
	created,
	status,
	reporter_name,
	reporter_email,
	reporter_phone,
	reporter_contact_consent,
	address_raw,
	address_number,
	address_street,
	address_locality,
	address_region,
	address_postal_code,
	address_country,
	address_id,
	location,
	h3cell,
	map_zoom,
	latlng_accuracy_type,
	latlng_accuracy_value,
	reviewed,
	reviewer_id,
	organization_id,
	report_type
)
SELECT 
	--id,
	public_id,
	created,
	status,
	COALESCE(reporter_name, ''),
	COALESCE(reporter_email, ''),
	COALESCE(reporter_phone, ''),
	reporter_contact_consent,
	address_raw,
	address_number,
	address_street,
	address_locality,
	address_region,
	address_postal_code,
	address_country,
	address_id,
	location,
	h3cell,
	map_zoom,
	latlng_accuracy_type,
	latlng_accuracy_value,
	reviewed,
	reviewer_id,
	organization_id,
	'nuisance' as report_type
FROM publicreport.nuisance;

-- Insert water reports into the report table
INSERT INTO publicreport.report (
	--id,
	public_id,
	created,
	status,
	reporter_name,
	reporter_email,
	reporter_phone,
	reporter_contact_consent,
	address_raw,
	address_number,
	address_street,
	address_locality,
	address_region,
	address_postal_code,
	address_country,
	address_id,
	location,
	h3cell,
	map_zoom,
	latlng_accuracy_type,
	latlng_accuracy_value,
	reviewed,
	reviewer_id,
	organization_id,
	report_type
)
SELECT 
	--id,
	public_id,
	created,
	status,
	COALESCE(reporter_name, ''),
	COALESCE(reporter_email, ''),
	COALESCE(reporter_phone, ''),
	reporter_contact_consent,
	address_raw,
	address_number,
	address_street,
	address_locality,
	address_region,
	address_postal_code,
	address_country,
	address_id,
	location,
	h3cell,
	map_zoom,
	-- Water table doesn't have these fields, so use defaults
	'none'::publicreport.accuracytype as latlng_accuracy_type,
	0.0 as latlng_accuracy_value,
	reviewed,
	reviewer_id,
	organization_id,
	'water' as report_type
FROM publicreport.water;

-- Rename existing tables
ALTER TABLE publicreport.nuisance RENAME TO nuisance_old;
ALTER TABLE publicreport.water RENAME TO water_old;

-- Create new nuisance table with only specific fields
CREATE TABLE publicreport.nuisance (
	additional_info TEXT NOT NULL,
	duration publicreport.nuisancedurationtype NOT NULL,
	is_location_backyard BOOLEAN NOT NULL,
	is_location_frontyard BOOLEAN NOT NULL,
	is_location_garden BOOLEAN NOT NULL,
	is_location_other BOOLEAN NOT NULL,
	is_location_pool BOOLEAN NOT NULL,
	report_id INTEGER REFERENCES publicreport.report(id),
	source_container BOOLEAN NOT NULL,
	source_description TEXT NOT NULL,
	source_stagnant BOOLEAN NOT NULL,
	source_gutter BOOLEAN NOT NULL,
	tod_early BOOLEAN NOT NULL,
	tod_day BOOLEAN NOT NULL,
	tod_evening BOOLEAN NOT NULL,
	tod_night BOOLEAN NOT NULL,
	PRIMARY KEY(report_id)
);

-- Create new water table with only specific fields
CREATE TABLE publicreport.water (
	access_comments TEXT NOT NULL,
	access_gate BOOLEAN NOT NULL,
	access_fence BOOLEAN NOT NULL,
	access_locked BOOLEAN NOT NULL,
	access_dog BOOLEAN NOT NULL,
	access_other BOOLEAN NOT NULL,
	comments TEXT NOT NULL,
	is_reporter_confidential BOOLEAN NOT NULL,
	is_reporter_owner BOOLEAN NOT NULL,
	has_adult BOOLEAN NOT NULL,
	has_backyard_permission BOOLEAN NOT NULL,
	has_larvae BOOLEAN NOT NULL,
	has_pupae BOOLEAN NOT NULL,
	owner_email TEXT NOT NULL,
	owner_name TEXT NOT NULL,
	owner_phone TEXT NOT NULL,
	report_id INTEGER REFERENCES publicreport.report(id),
	
	PRIMARY KEY(report_id)
);
-- Migrate nuisance-specific data
INSERT INTO publicreport.nuisance (
	report_id,
	additional_info,
	duration,
	source_container,
	source_description,
	source_stagnant,
	source_gutter,
	is_location_backyard,
	is_location_frontyard,
	is_location_garden,
	is_location_other,
	is_location_pool,
	tod_early,
	tod_day,
	tod_evening,
	tod_night
)
SELECT 
	r.id,
	n.additional_info,
	n.duration,
	n.source_container,
	n.source_description,
	n.source_stagnant,
	n.source_gutter,
	n.is_location_backyard,
	n.is_location_frontyard,
	n.is_location_garden,
	n.is_location_other,
	n.is_location_pool,
	n.tod_early,
	n.tod_day,
	n.tod_evening,
	n.tod_night
FROM publicreport.nuisance_old n
JOIN publicreport.report r ON r.public_id = n.public_id AND r.report_type = 'nuisance';

-- Migrate water-specific data
INSERT INTO publicreport.water (
	report_id,
	access_comments,
	access_gate,
	access_fence,
	access_locked,
	access_dog,
	access_other,
	comments,
	has_adult,
	has_larvae,
	has_pupae,
	owner_email,
	owner_name,
	owner_phone,
	has_backyard_permission,
	is_reporter_confidential,
	is_reporter_owner
)
SELECT 
	r.id,
	w.access_comments,
	w.access_gate,
	w.access_fence,
	w.access_locked,
	w.access_dog,
	w.access_other,
	w.comments,
	w.has_adult,
	w.has_larvae,
	w.has_pupae,
	w.owner_email,
	w.owner_name,
	w.owner_phone,
	w.has_backyard_permission,
	w.is_reporter_confidential,
	w.is_reporter_owner
FROM publicreport.water_old w
JOIN publicreport.report r ON r.public_id = w.public_id AND r.report_type = 'water';

-- Create new unified report_image junction table
CREATE TABLE publicreport.report_image (
	image_id INTEGER NOT NULL REFERENCES publicreport.image(id),
	report_id INTEGER NOT NULL REFERENCES publicreport.report(id),
	PRIMARY KEY (image_id, report_id)
);

-- Update nuisance_image table
ALTER TABLE publicreport.nuisance_image RENAME TO nuisance_image_old;

INSERT INTO publicreport.report_image (image_id, report_id)
SELECT ni.image_id, r.id as report_id
FROM publicreport.nuisance_image_old ni
JOIN publicreport.nuisance_old n ON n.id = ni.nuisance_id
JOIN publicreport.report r ON r.public_id = n.public_id AND r.report_type = 'nuisance';

-- Update water_image table
ALTER TABLE publicreport.water_image RENAME TO water_image_old;

INSERT INTO publicreport.report_image (image_id, report_id)
SELECT wi.image_id, r.id as report_id
FROM publicreport.water_image_old wi
JOIN publicreport.water_old w ON w.id = wi.water_id
JOIN publicreport.report r ON r.public_id = w.public_id AND r.report_type = 'water';

-- Create new unified notify_report_email
CREATE TABLE publicreport.notify_email (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	email_address TEXT NOT NULL REFERENCES comms.email_contact(address),
	report_id INTEGER NOT NULL REFERENCES publicreport.report(id),
	PRIMARY KEY (report_id, email_address)
);
-- Update notify_email_nuisance table
ALTER TABLE publicreport.notify_email_nuisance RENAME TO notify_email_nuisance_old;

INSERT INTO publicreport.notify_email (created, deleted, report_id, email_address)
SELECT nen.created, nen.deleted, r.id as report_id, nen.email_address
FROM publicreport.notify_email_nuisance_old nen
JOIN publicreport.nuisance_old n ON n.id = nen.nuisance_id
JOIN publicreport.report r ON r.public_id = n.public_id AND r.report_type = 'nuisance';

-- Update notify_email_water table
ALTER TABLE publicreport.notify_email_water RENAME TO notify_email_water_old;

INSERT INTO publicreport.notify_email (created, deleted, report_id, email_address)
SELECT new.created, new.deleted, r.id as report_id, new.email_address
FROM publicreport.notify_email_water_old new
JOIN publicreport.water_old w ON w.id = new.water_id
JOIN publicreport.report r ON r.public_id = w.public_id AND r.report_type = 'water';

-- Update notify_phone_nuisance table
ALTER TABLE publicreport.notify_phone_nuisance RENAME TO notify_phone_nuisance_old;

CREATE TABLE publicreport.notify_phone (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	deleted TIMESTAMP WITHOUT TIME ZONE,
	phone_e164 TEXT NOT NULL REFERENCES comms.phone(e164),
	report_id INTEGER NOT NULL REFERENCES publicreport.report(id),
	PRIMARY KEY (report_id, phone_e164)
);

INSERT INTO publicreport.notify_phone (created, deleted, report_id, phone_e164)
SELECT npn.created, npn.deleted, r.id as report_id, npn.phone_e164
FROM publicreport.notify_phone_nuisance_old npn
JOIN publicreport.nuisance_old n ON n.id = npn.nuisance_id
JOIN publicreport.report r ON r.public_id = n.public_id AND r.report_type = 'nuisance';

-- Update notify_phone_water table
ALTER TABLE publicreport.notify_phone_water RENAME TO notify_phone_water_old;

INSERT INTO publicreport.notify_phone (created, deleted, report_id, phone_e164)
SELECT npw.created, npw.deleted, r.id as report_id, npw.phone_e164
FROM publicreport.notify_phone_water_old npw
JOIN publicreport.water_old w ON w.id = npw.water_id
JOIN publicreport.report r ON r.public_id = w.public_id AND r.report_type = 'water';

