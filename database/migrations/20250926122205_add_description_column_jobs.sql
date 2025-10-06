-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs ADD COLUMN job_description TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs DROP COLUMN job_description;
-- +goose StatementEnd
