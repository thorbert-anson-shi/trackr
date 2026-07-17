-- name: UpdateUser :one
UPDATE users 
SET firebase_id = $2 
WHERE id = $1 
RETURNING *;
