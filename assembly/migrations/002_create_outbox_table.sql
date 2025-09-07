-- +goose Up

CREATE TABLE IF NOT EXISTS outbox (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type TEXT NOT NULL,
    payload JSONB,
    status TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS outbox;