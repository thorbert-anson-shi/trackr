-- name: ListFirebaseIDs :many
SELECT (firebase_id) 
FROM users 
WHERE is_admin = false;
