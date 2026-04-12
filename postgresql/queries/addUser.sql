-- name: AddUser :one
INSERT INTO users (name, registration_token) 
VALUES ($1, $2) RETURNING *;
