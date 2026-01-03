-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows
where user_id = $1;