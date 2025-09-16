-- +goose Up
-- +goose StatementBegin
CREATE TYPE transaction_type AS ENUM ('debit', 'credit');

ALTER TABLE transactions
ADD COLUMN type transaction_type NOT NULL DEFAULT 'debit';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
DROP COLUMN type;

DROP TYPE transaction_type;
-- +goose StatementEnd
