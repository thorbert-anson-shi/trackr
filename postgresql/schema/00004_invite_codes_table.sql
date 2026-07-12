-- +goose Up
CREATE TABLE invitations (
  id serial PRIMARY KEY,
  is_used boolean DEFAULT false NOT NULL,
  code varchar(26) NOT NULL,
  expiry_date timestamp(3) NOT NULL
);

-- +goose Down
DROP TABLE invitations;
