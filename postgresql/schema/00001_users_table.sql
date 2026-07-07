-- +goose Up
CREATE TABLE users (
	id serial PRIMARY KEY,
	name varchar(50),
	registration_token varchar,
	api_key varchar,
	is_admin bool DEFAULT false
);

CREATE INDEX idx_api_key ON users (api_key);

-- +goose Down
DROP TABLE users;
