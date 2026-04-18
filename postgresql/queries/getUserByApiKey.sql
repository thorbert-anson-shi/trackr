-- name: GetUserByApiKey :one
SELECT * FROM users WHERE users.api_key = $1;
