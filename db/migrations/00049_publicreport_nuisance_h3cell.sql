-- +goose Up
ALTER TABLE publicreport.nuisance ADD COLUMN h3cell h3index;
CREATE TABLE publicreport.nuisance_image (
	image_id INTEGER NOT NULL REFERENCES publicreport.image(id),
	nuisance_id INTEGER NOT NULL REFERENCES publicreport.nuisance(id),
	PRIMARY KEY (image_id, nuisance_id)
);
-- +goose Down
DROP TABLE publicreport.nuisance_image;
ALTER TABLE publicreport.nuisance DROP COLUMN h3cell;
