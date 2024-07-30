-- name: GetBySKU :one
SELECT sku, total_count, reserved
FROM items
WHERE sku = $1;
