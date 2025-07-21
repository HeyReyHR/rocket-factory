-- +goose Up
CREATE TYPE status AS ENUM ('unknown_status', 'pending_payment', 'cancelled', 'paid');

CREATE TYPE payment_method AS ENUM ('unknown_payment_method', 'credit_card', 'card', 'sbp', 'investor_money');


CREATE TABLE orders (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid TEXT NOT NULL,
    part_uuids TEXT[],
    total_price DOUBLE PRECISION,
    transaction_uuid UUID,
    status status NOT NULL DEFAULT 'pending_payment',
    payment_method payment_method,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS payment_method;