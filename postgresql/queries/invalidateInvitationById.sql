-- name: InvalidateInvitationById :exec
UPDATE invitations
SET invitations.is_used = true
WHERE invitations.id = $1;
