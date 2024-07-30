-- name: GetByID :one
SELECT
  o.order_id AS id,
  o.status,
  o.user_id AS user_id,
  o.created_at,
  o.updated_at,
  array_agg(i.item_sku)::int[] AS skus,
  array_agg(i.count)::int[] AS counts
FROM orders o
       JOIN orders_to_items i ON o.order_id = i.order_id
WHERE o.order_id = $1
GROUP BY o.order_id;
