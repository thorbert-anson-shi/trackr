-- name: AddUser :one
INSERT INTO users (name, firebase_id, api_key, is_admin) 
VALUES ($1, $2, $3, $4) 
RETURNING *;
