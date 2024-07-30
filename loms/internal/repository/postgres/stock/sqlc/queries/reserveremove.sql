-- name: ReserveRemove :exec
UPDATE items
SET
  reserved = GREATEST(reserved - data.count, 0),
  total_count = GREATEST(total_count - data.count, 0)
FROM (SELECT unnest(@sku::int[]) AS sku, unnest(@count::int[]) AS count) AS data
WHERE items.sku = data.sku;
