-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255),
     rating INT,
     category VARCHAR(255),
     image_url TEXT,
     reputation INT,
     price INT,
     availability INT
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
