-- +goose Up
CREATE TABLE locations (
	id serial PRIMARY KEY,
	user_id integer REFERENCES users(id) ON DELETE CASCADE,
	latitude real,
	longitude real,
	timestamp timestamp,
	accuracy real
);

-- +goose Down
DROP TABLE locations;
