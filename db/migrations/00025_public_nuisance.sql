-- +goose Up
CREATE TYPE publicreport.NuisanceDurationType AS ENUM (
	'none',
	'just-noticed',
	'few-days',
	'1-2-weeks',
	'2-4-weeks',
	'1-3-months',
	'seasonal'
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
CREATE TYPE publicreport.NuisancePreferredDateRangeType AS ENUM (
	'none',
	'any-time',
	'in-two-weeks',
	'next-week'
);
CREATE TYPE publicreport.NuisancePreferredTimeType AS ENUM (
	'none',
	'afternoon',
	'any-time',
	'morning'
);
CREATE TABLE publicreport.nuisance (
	id SERIAL PRIMARY KEY,
	additional_info TEXT NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	duration publicreport.NuisanceDurationType NOT NULL,
	email TEXT NOT NULL,
	inspection_type publicreport.NuisanceInspectionType NOT NULL,
	location publicreport.NuisanceLocationType NOT NULL,
	preferred_date_range publicreport.NuisancePreferredDateRangeType NOT NULL,
	preferred_time publicreport.NuisancePreferredTimeType NOT NULL,
	request_call BOOLEAN NOT NULL,
	severity SMALLINT NOT NULL,
	source_container BOOLEAN NOT NULL,
	source_description TEXT NOT NULL,
	source_roof BOOLEAN NOT NULL,
	source_stagnant BOOLEAN NOT NULL,
	time_of_day_day BOOLEAN NOT NULL,
	time_of_day_early BOOLEAN NOT NULL,
	time_of_day_evening BOOLEAN NOT NULL,
	time_of_day_night BOOLEAN NOT NULL,
	public_id TEXT NOT NULL UNIQUE,
	reporter_address TEXT NOT NULL,
	reporter_email TEXT NOT NULL,
	reporter_name TEXT NOT NULL,
	reporter_phone TEXT NOT NULL
);

-- +goose Down
DROP TABLE publicreport.nuisance;
DROP TYPE publicreport.NuisanceDurationType;
DROP TYPE publicreport.NuisanceLocationType;
DROP TYPE publicreport.NuisanceInspectionType;
DROP TYPE publicreport.NuisancePreferredDateRangeType;
DROP TYPE publicreport.NuisancePreferredTimeType;
