-- +goose Up
ALTER TYPE fileupload.PoolConditionType RENAME VALUE 'empty' TO 'dry';
ALTER TYPE fileupload.PoolConditionType ADD VALUE 'false pool' AFTER 'dry';
ALTER TYPE fileupload.FileStatusType ADD VALUE 'committed' AFTER 'uploaded';
ALTER TYPE fileupload.FileStatusType ADD VALUE 'discarded' AFTER 'committed';

