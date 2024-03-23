-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item_transactions (
     id SERIAL PRIMARY KEY,
     user_id INT,
     item_id INT,
     transaction_type VARCHAR(50),
     transaction_date TIMESTAMP,
     quantity INT,
     total_price INT,
     availability INT,
     FOREIGN KEY (user_id) REFERENCES users(id),
     FOREIGN KEY (item_id) REFERENCES items(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item_transactions;
-- +goose StatementEnd
