-- +goose Up
CREATE TABLE users (
	id serial PRIMARY KEY,
	name varchar(50),
	registration_token varchar,
	api_key varchar,
	is_admin bool DEFAULT false
);

-- +goose Down
DROP TABLE users;
