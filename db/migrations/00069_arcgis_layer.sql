-- +goose Up
CREATE TABLE arcgis.feature_service (
	extent box2d NOT NULL,
	item_id TEXT NOT NULL,
	spatial_reference INTEGER NOT NULL,
	url TEXT NOT NULL,
	PRIMARY KEY(item_id)
);
CREATE TABLE arcgis.layer (
	extent box2d NOT NULL,
	feature_service_item_id TEXT NOT NULL REFERENCES arcgis.feature_service(item_id),
	index_ INTEGER NOT NULL,
	PRIMARY KEY(feature_service_item_id, index_)
);
CREATE TYPE arcgis.FieldType AS ENUM (
	'esriFieldTypeSmallInteger', -- 16-bit Integer.
	'esriFieldTypeInteger', -- 32-bit Integer.
	'esriFieldTypeSingle', -- Single-precision floating-point number.
	'esriFieldTypeDouble', -- Double-precision floating-point number.
	'esriFieldTypeString', -- Character string.
	'esriFieldTypeDate', -- Date.
	'esriFieldTypeOID', -- Integer representing an object identifier. 32-bit OID has a length of 4 bytes, and 64-bit OID has a length of 8 bytes.
	'esriFieldTypeGeometry', -- Geometry.
	'esriFieldTypeBlob', -- Binary Large Object.
	'esriFieldTypeRaster', -- Raster.
	'esriFieldTypeGUID', -- Globally Unique Identifier.
	'esriFieldTypeGlobalID', -- Esri Global ID.
	'esriFieldTypeXML', -- XML Document.
	'esriFieldTypeBigInteger' -- 64-bit Integer.
);
CREATE TABLE arcgis.layer_field (
	layer_feature_service_item_id TEXT NOT NULL,
	layer_index INTEGER NOT NULL,
	name TEXT NOT NULL,
	type_ arcgis.FieldType NOT NULL,
	FOREIGN KEY(layer_feature_service_item_id, layer_index) REFERENCES arcgis.layer(feature_service_item_id, index_),
	PRIMARY KEY(layer_feature_service_item_id, layer_index, name)
);
CREATE TYPE arcgis.MappingDestinationParcel AS ENUM (
	'apn',
	'description'
);
CREATE TABLE arcgis.parcel_mapping (
	destination arcgis.MappingDestinationParcel NOT NULL,
	layer_feature_service_item_id TEXT NOT NULL,
	layer_index INTEGER NOT NULL,
	layer_field_name TEXT NOT NULL,
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	FOREIGN KEY(layer_feature_service_item_id, layer_index, layer_field_name) REFERENCES arcgis.layer_field(layer_feature_service_item_id, layer_index, name),
	PRIMARY KEY(organization_id, destination)
);

CREATE TYPE arcgis.MappingDestinationAddress AS ENUM (
	'country',
	'locality',
	'postal_code',
	'street',
	'unit'
);
CREATE TABLE arcgis.address_mapping (
	destination arcgis.MappingDestinationAddress NOT NULL,
	layer_feature_service_item_id TEXT NOT NULL,
	layer_index INTEGER NOT NULL,
	layer_field_name TEXT NOT NULL,
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	FOREIGN KEY(layer_feature_service_item_id, layer_index, layer_field_name) REFERENCES arcgis.layer_field(layer_feature_service_item_id, layer_index, name),
	PRIMARY KEY(organization_id, destination)
);
-- +goose Down
DROP TABLE arcgis.address_mapping;
DROP TYPE arcgis.MappingDestinationAddress;
DROP TABLE arcgis.parcel_mapping;
DROP TYPE arcgis.MappingDestinationParcel;
DROP TABLE arcgis.layer_field;
DROP TYPE arcgis.FieldType;
DROP TABLE arcgis.layer;
DROP TABLE arcgis.feature_service;
