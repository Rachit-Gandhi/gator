-- name: GetFeedByUrl :one
SELECT * FROM feeds
where url = $1;