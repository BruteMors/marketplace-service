-- name: FetchNextOrderStatusChangedEvent :one
SELECT id, order_id, status, at
FROM order_status_changed_events
WHERE sent = FALSE
ORDER BY at ASC
LIMIT 1;
