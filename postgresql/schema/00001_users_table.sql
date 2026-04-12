-- +goose Up
CREATE TABLE users (
	id serial PRIMARY KEY,
	name varchar(50)
);

-- +goose Down
DROP TABLE users;
