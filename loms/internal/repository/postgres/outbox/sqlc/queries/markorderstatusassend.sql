-- name: MarkOrderStatusChangedEventAsSend :exec
UPDATE order_status_changed_events
SET sent = TRUE
WHERE id = $1;
