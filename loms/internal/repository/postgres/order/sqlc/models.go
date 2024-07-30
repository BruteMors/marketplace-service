// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "new"
	OrderStatusAwaitingpayment OrderStatus = "awaiting payment"
	OrderStatusFailed          OrderStatus = "failed"
	OrderStatusPayed           OrderStatus = "payed"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

func (e *OrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatus(s)
	case string:
		*e = OrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatus: %T", src)
	}
	return nil
}

type NullOrderStatus struct {
	OrderStatus OrderStatus
	Valid       bool // Valid is true if OrderStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderStatus) Scan(value interface{}) error {
	if value == nil {
		ns.OrderStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrderStatus), nil
}

type Item struct {
	Sku        int32
	TotalCount int32
	Reserved   int32
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

type Order struct {
	OrderID   int64
	UserID    int32
	Status    OrderStatus
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type OrderStatusChangedEvent struct {
	ID      int64
	OrderID int64
	Status  OrderStatus
	At      pgtype.Timestamp
	Sent    bool
}

type OrdersToItem struct {
	ID      int32
	OrderID int64
	ItemSku int32
	Count   int32
}