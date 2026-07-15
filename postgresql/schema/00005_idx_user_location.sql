-- +goose Up
CREATE INDEX idx_location_user_timestamp ON locations (user_id, timestamp DESC);

-- +goose Down
DROP INDEX idx_location_user_timestamp;
