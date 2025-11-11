-- +goose Up
CREATE TYPE NotificationType AS ENUM (
	'oauth_token_invalidated');

CREATE TABLE notification (
	id SERIAL PRIMARY KEY,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	link TEXT NOT NULL,
	message TEXT NOT NULL,
	type NotificationType NOT NULL,
	user_id INTEGER REFERENCES user_(id) NOT NULL);



-- +goose Down
DROP TABLE notification;
DROP TYPE NotificationType;
