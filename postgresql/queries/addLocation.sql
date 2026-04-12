-- name: AddLocation :one
INSERT INTO locations (
	user_id,
	latitude,
	longitude,
	timestamp,
	accuracy
) VALUES (
	$1, $2, $3, $4, $5
) RETURNING *;
