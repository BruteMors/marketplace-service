-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('new', 'awaiting payment', 'failed', 'payed', 'cancelled');

CREATE TABLE IF NOT EXISTS "orders" (
                                      order_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                      user_id INTEGER NOT NULL,
                                      status order_status NOT NULL,
                                      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "items" (
                                     sku INTEGER PRIMARY KEY,
                                     total_count INTEGER NOT NULL,
                                     reserved INTEGER NOT NULL DEFAULT 0,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "orders_to_items" (
                                               id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                               order_id BIGINT NOT NULL REFERENCES "orders" (order_id),
                                               item_sku INTEGER NOT NULL REFERENCES "items" (sku),
                                               count INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "orders_to_items";
DROP TABLE IF EXISTS "items";
DROP TABLE IF EXISTS "orders";
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
