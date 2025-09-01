-- +goose Up
ALTER TYPE status ADD VALUE IF NOT EXISTS 'assembled';

-- +goose Down
CREATE TYPE status_new AS ENUM ('unknown_status', 'pending_payment', 'cancelled', 'paid');

ALTER TABLE orders
  ALTER COLUMN status TYPE status_new USING status::text::status_new;

DROP TYPE status;

ALTER TYPE status_new RENAME TO status;
