-- name: InvalidateInvitationById :exec
UPDATE invitations
SET is_used = true
WHERE id = $1;
