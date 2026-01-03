-- name: GetFeedNameById :one
SELECT name FROM feeds
WHERE id = $1;