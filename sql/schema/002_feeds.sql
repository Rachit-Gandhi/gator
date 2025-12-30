-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    url VARCHAR(250) NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    CONSTRAINT fk_feeds_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
