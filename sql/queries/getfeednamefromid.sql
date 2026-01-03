-- name: GetFeedNameById :one
SELECT name FROM feeds
where id = $1;