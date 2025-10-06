-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs ADD COLUMN company_name VARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs DROP COLUMN company_name;
-- +goose StatementEnd
