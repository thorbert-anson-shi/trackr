-- +goose Up
CREATE TABLE users (
	id serial PRIMARY KEY,
	name varchar(50) NOT NULL,
	registration_token varchar,
	api_key varchar NOT NULL,
	is_admin bool DEFAULT false NOT NULL
);

-- +goose Down
DROP TABLE users;
