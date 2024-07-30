-- +goose Up
-- +goose StatementBegin
INSERT INTO items(sku, total_count) VALUES (1076963, 100), (1148162, 200);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM orders_to_items WHERE item_sku IN (1076963, 1148162);
DELETE FROM items WHERE sku IN (1076963, 1148162);
-- +goose StatementEnd
