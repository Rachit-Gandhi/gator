-- +goose Up
ALTER TABLE feeds
  ADD COLUMN created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
  ADD COLUMN updated_at  TIMESTAMP NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE feeds
  DROP COLUMN created_at,
  DROP COLUMN updated_at;