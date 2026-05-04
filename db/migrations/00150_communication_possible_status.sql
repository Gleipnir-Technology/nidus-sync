-- +goose Up
ALTER TABLE communication 
	ADD COLUMN set_possible_issue TIMESTAMP WITHOUT TIME ZONE,
	ADD COLUMN set_possible_issue_by INTEGER REFERENCES user_(id),
	ADD COLUMN set_possible_resolved TIMESTAMP WITHOUT TIME ZONE,
	ADD COLUMN set_possible_resolved_by INTEGER REFERENCES user_(id);
