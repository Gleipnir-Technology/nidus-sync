-- +goose Up
ALTER TABLE user_ ADD CONSTRAINT user_username_unique UNIQUE (username);
