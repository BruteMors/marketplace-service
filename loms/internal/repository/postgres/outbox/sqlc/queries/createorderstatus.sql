-- name: CreateOrderStatusChangedEvent :exec
INSERT INTO order_status_changed_events (order_id, status)
VALUES ($1, $2);
