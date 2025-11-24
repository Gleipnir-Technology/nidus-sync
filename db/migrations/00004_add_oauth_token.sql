-- +goose Up
CREATE TABLE oauth_token (
	id SERIAL PRIMARY KEY,
	access_token TEXT NOT NULL,
	expires TIMESTAMP NOT NULL,
	refresh_token TEXT NOT NULL,
	username TEXT NOT NULL,
	user_id INTEGER REFERENCES user_ (id) NOT NULL
);

-- +goose Down
DROP TABLE oauth_token;
