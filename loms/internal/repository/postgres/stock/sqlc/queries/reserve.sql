-- name: GetItemsAvailability :many
SELECT sku, (total_count - reserved) AS available
FROM items
WHERE sku = ANY(@sku::int[])
  FOR UPDATE;

-- name: UpdateReservedItems :exec
WITH unnested_data AS (
  SELECT unnest(@sku::int[]) AS sku, unnest(@count::int[]) AS count
)
UPDATE items
SET reserved = items.reserved + ud.count
FROM unnested_data ud
WHERE items.sku = ud.sku AND (items.total_count - items.reserved) >= ud.count;



