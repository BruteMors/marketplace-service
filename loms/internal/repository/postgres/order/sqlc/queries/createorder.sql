-- name: CreateOrder :one
INSERT INTO "orders" (user_id, status, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING order_id;

-- name: InsertOrderItems :exec
INSERT INTO "orders_to_items" (order_id, item_sku, count)
SELECT $1, unnest(@item_sku::int[]), unnest(@count::int[]);

