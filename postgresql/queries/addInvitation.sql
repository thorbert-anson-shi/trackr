-- name: AddInvitation :one
INSERT INTO invitations (
  code, expiry_date
) VALUES ( 
  $1, $2 
) RETURNING *;
