-- +goose Up
CREATE TABLE locations (
	id serial PRIMARY KEY,
	user_id integer REFERENCES users(id) ON DELETE CASCADE,
	latitude real NOT NULL,
	longitude real NOT NULL,
	timestamp timestamp NOT NULL,
	accuracy real
);

-- +goose Down
DROP TABLE locations;
