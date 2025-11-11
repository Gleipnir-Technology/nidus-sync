-- +goose Up
CREATE TYPE NotificationType AS ENUM (
	'oauth_token_invalidated');

CREATE TABLE notification (
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES user_(id),
	message TEXT,
	link TEXT,
	type NotificationType);



-- +goose Down
DROP TABLE notification;
DROP TYPE NotificationType;
