-- name: UpdateUser :one
UPDATE users 
SET registration_token = $2 
WHERE id = $1 
RETURNING *;
