-- name: GetInvitationByCode :one
SELECT * FROM invitations WHERE invitations.code = $1;
