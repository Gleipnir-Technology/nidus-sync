-- +goose Up
ALTER TYPE PoolConditionType ADD VALUE 'unknown' AFTER 'false pool';
ALTER TABLE fileupload.pool DROP COLUMN condition;
DROP TYPE fileupload.PoolConditionType;
ALTER TABLE fileupload.pool ADD COLUMN condition PoolConditionType NOT NULL;
-- +goose Down
ALTER TABLE fileupload.pool DROP COLUMN condition;
CREATE TYPE fileupload.PoolConditionType AS ENUM (
	'green',
	'murky',
	'blue',
	'dry',
	'false pool',
	'unknown'
);
ALTER TABLE fileupload.pool ADD COLUMN condition fileupload.PoolConditionType NOT NULL;
