-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, post_url, post_description, published_at, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
ON CONFLICT (post_url) DO NOTHING
RETURNING *;