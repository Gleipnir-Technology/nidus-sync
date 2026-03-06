-- +goose Up
CREATE TABLE review_task (
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	creator_id INTEGER REFERENCES user_(id),
	id SERIAL NOT NULL,
	organization_id INTEGER NOT NULL REFERENCES organization(id),
	reviewed TIMESTAMP WITHOUT TIME ZONE,
	reviewer_id INTEGER REFERENCES user_(id),
	PRIMARY KEY(id)
);
CREATE TABLE review_task_pool (
	feature_pool INTEGER NOT NULL REFERENCES feature_pool(feature_id),
	location geometry(Point, 4326),
	geometry geometry(Polygon, 4326),
	review_task_id INTEGER NOT NULL REFERENCES review_task(id),
	PRIMARY KEY(review_task_id)
);
-- +goose Down
DROP TABLE review_task_pool;
DROP TABLE review_task;
