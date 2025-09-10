-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs
DROP COLUMN origin_id,
ADD COLUMN origin_id UUID NULL REFERENCES companies(id);

CREATE UNIQUE INDEX jobs_name_idx ON categories (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs
DROP COLUMN origin_id,
ADD COLUMN origin_id INTEGER;

DROP INDEX jobs_name_idx;
-- +goose StatementEnd
