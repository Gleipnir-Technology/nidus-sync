-- +goose Up
ALTER TYPE fileupload.PoolConditionType ADD VALUE 'empty' AFTER 'blue';
