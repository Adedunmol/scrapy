-- +goose Up
-- +goose StatementBegin
ALTER TABLE companies
ADD COLUMN user_id UUID REFERENCES users(id) UNIQUE,
ADD COLUMN email VARCHAR(255) UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE companies
    DROP COLUMN user_id UUID,
    DROP COLUMN email;
-- +goose StatementEnd
