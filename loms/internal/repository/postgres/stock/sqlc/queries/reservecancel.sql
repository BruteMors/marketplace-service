-- name: ReserveCancel :exec
WITH unnested_data AS (
  SELECT unnest(@sku::int[]) AS sku, unnest(@count::int[]) AS count
)
UPDATE items
SET reserved = GREATEST(reserved - unnested_data.count, 0)
FROM unnested_data
WHERE items.sku = unnested_data.sku;
