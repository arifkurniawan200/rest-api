-- +goose Up
-- +goose StatementBegin
ALTER TABLE items
ADD COLUMN value VARCHAR(10) CHECK (value IN ('red', 'yellow', 'green'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
