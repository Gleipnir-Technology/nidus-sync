-- +goose Up
CREATE TYPE publicreport.AccuracyType AS ENUM (
	'rooftop', -- Result intersects a known building/entrance.
	'parcel', -- Result is associated with one or more address within a specified polygonal boundary
	'point', -- Result is a known address point but does not intersect a known rooftop/parcel.
	'interpolated', -- Result position and existence are estimated based on nearby known addresses.
	'approximate', -- Result position is approximated by a 9-digit zipcode centroid.
	'intersection', -- For street type features only. The result is an intersection of 2 streets.
	'browser', -- added to signify that we got location from the browser
	'none' -- we have no accuarcy at all
);
COMMENT ON TYPE publicreport.AccuracyType IS 'most volues are determined by our geocoding API provider, mapbox. You can read more details at https://docs.mapbox.com/api/search/geocoding/#point-accuracy-for-address-features.';

ALTER TABLE publicreport.nuisance ADD COLUMN address_country TEXT;
ALTER TABLE publicreport.nuisance ADD COLUMN address_place TEXT;
ALTER TABLE publicreport.nuisance ADD COLUMN address_postcode TEXT;
ALTER TABLE publicreport.nuisance ADD COLUMN address_region TEXT;
ALTER TABLE publicreport.nuisance ADD COLUMN address_street TEXT;
ALTER TABLE publicreport.nuisance ADD COLUMN is_location_backyard BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN is_location_frontyard BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN is_location_garden BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN is_location_other BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN is_location_pool BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN map_zoom REAL;
ALTER TABLE publicreport.nuisance ADD COLUMN tod_early BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN tod_day BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN tod_evening BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN tod_night BOOLEAN;
ALTER TABLE publicreport.nuisance ADD COLUMN latlng_accuracy_type publicreport.AccuracyType;
ALTER TABLE publicreport.nuisance ADD COLUMN latlng_accuracy_value REAL;

UPDATE publicreport.nuisance SET address_country = '';
UPDATE publicreport.nuisance SET address_place = '';
UPDATE publicreport.nuisance SET address_postcode = '';
UPDATE publicreport.nuisance SET address_region = '';
UPDATE publicreport.nuisance SET address_street = '';
UPDATE publicreport.nuisance SET is_location_backyard = false;
UPDATE publicreport.nuisance SET is_location_frontyard = false;
UPDATE publicreport.nuisance SET is_location_garden = false;
UPDATE publicreport.nuisance SET is_location_other = false;
UPDATE publicreport.nuisance SET is_location_pool = false;
UPDATE publicreport.nuisance SET map_zoom = 0;
UPDATE publicreport.nuisance SET tod_early = false;
UPDATE publicreport.nuisance SET tod_day = false;
UPDATE publicreport.nuisance SET tod_evening = false;
UPDATE publicreport.nuisance SET tod_night = false;
UPDATE publicreport.nuisance SET latlng_accuracy_type = 'approximate';
UPDATE publicreport.nuisance SET latlng_accuracy_value = 0.0;

ALTER TABLE publicreport.nuisance ALTER COLUMN address_country SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN address_place SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN address_postcode SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN address_region SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN address_street SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN is_location_backyard SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN is_location_frontyard SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN is_location_garden SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN is_location_other SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN is_location_pool SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN map_zoom SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN tod_early SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN tod_day SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN tod_evening SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN tod_night SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN latlng_accuracy_type SET NOT NULL;
ALTER TABLE publicreport.nuisance ALTER COLUMN latlng_accuracy_value SET NOT NULL;

ALTER TABLE publicreport.nuisance DROP COLUMN source_location;

DROP TYPE publicreport.NuisanceInspectionType;
DROP TYPE publicreport.NuisanceLocationType;
DROP TYPE publicreport.NuisancePreferredDateRangeType;
DROP TYPE publicreport.NuisancePreferredTimeType;

-- +goose Down
CREATE TYPE publicreport.NuisancePreferredTimeType AS ENUM (
	'none',
	'afternoon',
	'any-time',
	'morning'
);
CREATE TYPE publicreport.NuisancePreferredDateRangeType AS ENUM (
	'none',
	'any-time',
	'in-two-weeks',
	'next-week'
);
CREATE TYPE publicreport.NuisanceLocationType AS ENUM (
	'none',
	'front-yard',
	'backyard',
	'patio',
	'garden',
	'pool-area',
	'throughout',
	'indoors',
	'other'
);
CREATE TYPE publicreport.NuisanceInspectionType AS ENUM (
	'neighborhood',
	'property'
);

ALTER TABLE publicreport.nuisance ADD COLUMN source_location publicreport.NuisanceLocationType;

ALTER TABLE publicreport.nuisance DROP COLUMN address_country;
ALTER TABLE publicreport.nuisance DROP COLUMN address_place;
ALTER TABLE publicreport.nuisance DROP COLUMN address_postcode;
ALTER TABLE publicreport.nuisance DROP COLUMN address_region;
ALTER TABLE publicreport.nuisance DROP COLUMN address_street;
ALTER TABLE publicreport.nuisance DROP COLUMN is_location_frontyard;
ALTER TABLE publicreport.nuisance DROP COLUMN is_location_backyard;
ALTER TABLE publicreport.nuisance DROP COLUMN is_location_garden;
ALTER TABLE publicreport.nuisance DROP COLUMN is_location_pool;
ALTER TABLE publicreport.nuisance DROP COLUMN is_location_other;
ALTER TABLE publicreport.nuisance DROP COLUMN map_zoom;
ALTER TABLE publicreport.nuisance DROP COLUMN tod_early;
ALTER TABLE publicreport.nuisance DROP COLUMN tod_day;
ALTER TABLE publicreport.nuisance DROP COLUMN tod_evening;
ALTER TABLE publicreport.nuisance DROP COLUMN tod_night;
ALTER TABLE publicreport.nuisance DROP COLUMN latlng_accuracy_type;
ALTER TABLE publicreport.nuisance DROP COLUMN latlng_accuracy_value;

DROP TYPE publicreport.AccuracyType;
