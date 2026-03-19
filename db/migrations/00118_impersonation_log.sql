-- +goose Up
CREATE TABLE log_impersonation (
	begin_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	end_at TIMESTAMP WITHOUT TIME ZONE,
	id SERIAL NOT NULL,
	impersonator_id INTEGER NOT NULL REFERENCES user_(id),
	target_id INTEGER NOT NULL REFERENCES user_(id),
	PRIMARY KEY (id)
);
-- +goose Down
DROP TABLE log_impersonation;

