-- +goose Up
ALTER TABLE history_treatment ADD COLUMN created TIMESTAMP;
ALTER TABLE history_proposedtreatmentarea ADD COLUMN created TIMESTAMP;
ALTER TABLE history_polygonlocation ADD COLUMN created TIMESTAMP;

-- +goose Down
ALTER TABLE history_treatment DROP COLUMN created;
ALTER TABLE history_proposedtreatmentarea DROP COLUMN created;
ALTER TABLE history_polygonlocation DROP COLUMN created;
