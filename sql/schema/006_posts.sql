-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    title VARCHAR(100) NOT NULL,
    post_url VARCHAR(500) NOT NULL UNIQUE,
    post_description VARCHAR(500) NOT NULL,
    published_at TIMESTAMP NOT NULL DEFAULT NOW(),
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;