-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT orders_uuid_unique UNIQUE (uuid)
);

-- trigger function to maintain updated_at
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS set_updated_at ON orders;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at ON orders;
DROP FUNCTION IF EXISTS trigger_set_updated_at();
DROP TABLE IF EXISTS orders;
DROP EXTENSION IF EXISTS "pgcrypto";