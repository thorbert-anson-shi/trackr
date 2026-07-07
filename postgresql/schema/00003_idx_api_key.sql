-- +goose Up
CREATE INDEX idx_api_key ON users (api_key);

-- +goose Down
DROP INDEX idx_api_key;
