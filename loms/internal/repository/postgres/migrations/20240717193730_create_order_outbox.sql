-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order_status_changed_events" (
                                                           id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                                           order_id BIGINT NOT NULL REFERENCES "orders" (order_id),
                                                           status order_status NOT NULL,
                                                           at TIMESTAMP NOT NULL DEFAULT NOW(),
                                                           sent BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order_status_changed_events";
-- +goose StatementEnd
