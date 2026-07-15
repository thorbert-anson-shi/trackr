-- name: ListLatestLocations :many
SELECT DISTINCT ON (user_id) * 
FROM locations 
ORDER BY user_id, timestamp DESC;
