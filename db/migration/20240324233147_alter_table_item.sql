-- +goose Up
-- +goose StatementBegin
ALTER TABLE items
ADD COLUMN is_active BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
