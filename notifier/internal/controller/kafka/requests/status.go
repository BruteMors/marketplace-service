package requests

import "time"

type StatusChangedEvent struct {
	ID      int64     `json:"id"`
	OrderID int64     `json:"order_id"`
	Status  Status    `json:"status"`
	At      time.Time `json:"at"`
}

type Status string

const (
	OrderStatusNew             Status = "new"
	OrderStatusAwaitingPayment Status = "awaiting payment"
	OrderStatusFailed          Status = "failed"
	OrderStatusPayed           Status = "payed"
	OrderStatusCancelled       Status = "cancelled"
)

func (s Status) String() string {
	return string(s)
}
