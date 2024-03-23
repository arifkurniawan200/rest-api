-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS table_history (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) NOT NULL,
    table_key int NOT NULL,
    data_before JSONB NOT NULL,
    data_after JSONB NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
