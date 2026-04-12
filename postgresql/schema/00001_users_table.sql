-- +goose Up
CREATE TABLE users (
	id serial PRIMARY KEY,
	name varchar(50),
	registration_token varchar
);

-- +goose Down
DROP TABLE users;
