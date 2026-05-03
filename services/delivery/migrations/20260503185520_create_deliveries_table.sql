-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS deliveries (
    id BIGSERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT deliveries_order_id_unique UNIQUE (order_id)
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

DROP TRIGGER IF EXISTS set_updated_at ON deliveries;
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON deliveries
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS deliveries;
DROP TRIGGER IF EXISTS set_updated_at ON deliveries;
DROP FUNCTION IF EXISTS trigger_set_updated_at();
DROP EXTENSION IF EXISTS "pgcrypto";
