-- +goose Up
ALTER TABLE review_task_pool ADD COLUMN condition PoolConditionType;
CREATE TYPE ReviewTaskResolutionType AS ENUM (
	'committed',
	'discarded'
);
ALTER TABLE review_task ADD COLUMN resolution ReviewTaskResolutionType;
-- +goose Down
ALTER TABLE review_task DROP COLUMN resolution;
DROP TYPE ReviewTaskResolutionType
ALTER TABLE review_task_pool DROP COLUMN condition;
