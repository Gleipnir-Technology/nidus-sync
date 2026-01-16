-- +goose Up
CREATE TABLE publicreport.image (
	id SERIAL,
	content_type TEXT NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	location GEOGRAPHY, 
	resolution_x INTEGER NOT NULL,
	resolution_y INTEGER NOT NULL,
	storage_uuid UUID NOT NULL,
	storage_size BIGINT NOT NULL,
	uploaded_filename TEXT NOT NULL,
	PRIMARY KEY(id)
);
CREATE TABLE publicreport.image_exif (
	image_id INTEGER NOT NULL REFERENCES publicreport.image(id),
	name TEXT NOT NULL,
	value TEXT NOT NULL,
	PRIMARY KEY(image_id, name, value)
);
CREATE TABLE publicreport.pool_image (
	image_id INTEGER NOT NULL REFERENCES publicreport.image(id),
	pool_id INTEGER NOT NULL REFERENCES publicreport.pool(id),
	PRIMARY KEY (image_id, pool_id)
);
CREATE TABLE publicreport.quick_image (
	image_id INTEGER NOT NULL REFERENCES publicreport.image(id),
	quick_id INTEGER NOT NULL REFERENCES publicreport.quick(id),
	PRIMARY KEY (image_id, quick_id)
);
DROP TABLE IF EXISTS publicreport.pool_photo;
DROP TABLE IF EXISTS publicreport.quick_photo;
-- +goose Down
DROP TABLE publicreport.quick_image;
DROP TABLE publicreport.pool_image;
DROP TABLE publicreport.image_exif;
DROP TABLE publicreport.image;
-- that's right, I'm not rebuilding the pool_photo or quick_photo tables because I'm lazy.
