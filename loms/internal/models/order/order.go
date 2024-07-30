package order

import "time"

type Order struct {
	ID        int64
	Status    Status
	UserID    int64
	Items     []Item
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type NewOrder struct {
	User   int64
	Items  []Item
	Status Status
}

type Item struct {
	SKU   uint32
	Count uint16
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

type StatusChangedEvent struct {
	ID      int64     `json:"id"`
	OrderID int64     `json:"order_id"`
	Status  Status    `json:"status"`
	At      time.Time `json:"at"`
}
