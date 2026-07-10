package models

import (
	db "restaurant-platform/database/sqlc/gen"
	"time"

	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"   // waiter submitted
	OrderStatusConfirmed OrderStatus = "confirmed" // kitchen acknowledged
	OrderStatusPreparing OrderStatus = "preparing" // kitchen started
	OrderStatusReady     OrderStatus = "ready"     // kitchen done, waiter picks up
	OrderStatusServed    OrderStatus = "served"    // delivered to table
	OrderStatusPaid      OrderStatus = "paid"      // payment done
	OrderStatusCancelled OrderStatus = "cancelled"
)

// ValidTransitions defines allowed status moves
var ValidTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusCreated:   {OrderStatusConfirmed, OrderStatusCancelled},
	OrderStatusConfirmed: {OrderStatusPreparing, OrderStatusCancelled},
	OrderStatusPreparing: {OrderStatusReady, OrderStatusCancelled},
	OrderStatusReady:     {OrderStatusServed},
	OrderStatusServed:    {OrderStatusPaid},
	OrderStatusPaid:      {},
	OrderStatusCancelled: {},
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	allowed := ValidTransitions[s]
	for _, a := range allowed {
		if a == next {
			return true
		}
	}
	return false
}

type Order struct {
	ID          int64           `json:"id" db:"id"`
	TenantID    int64           `json:"tenant_id" db:"tenant_id"`
	TableID     int64           `json:"table_id" db:"table_id"`
	UserID      int64           `json:"user_id" db:"user_id"`
	Status      string          `json:"status" db:"status"`
	Notes       string          `json:"notes" db:"notes"`
	TotalAmount decimal.Decimal `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

func ConvertOrderModel(order db.Order) Order {
	return Order{
		ID:          order.ID,
		TenantID:    order.TenantID,
		TableID:     order.TableID,
		UserID:      order.UserID,
		Status:      order.Status,
		Notes:       order.Notes.String,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt.Time,
		UpdatedAt:   order.UpdatedAt.Time,
	}
}

func ConvertOrderModels(orders []db.Order) []Order {
	result := make([]Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, ConvertOrderModel(order))
	}
	return result
}

type ListOrdersReq struct {
	TenantID int64 `json:"tenant_id"`
	TableID  int64 `json:"table_id"`
	UserID   int64 `json:"user_id"`
	Page     int   `json:"page"`
	Limit    int   `json:"limit"`
	Offset   int   `json:"offset"`
}

type UpdateOrderStatusReq struct {
	Status   string `json:"status"`
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	TableID  int64  `json:"table_id"`
	UserID   int64  `json:"user_id"`
}

type CreateOrderInput struct {
	ReferenceID string           `json:"reference_id"`
	TenantID    int64            `json:"tenant_id"`
	TableID     int64            `json:"table_id"`
	UserID      int64            `json:"user_id"`
	Notes       string           `json:"notes"`
	Items       []OrderItemInput `json:"items"`
}

type OrderItemInput struct {
	MenuItemID int64  `json:"menu_item_id"`
	Quantity   int    `json:"quantity"`
	Notes      string `json:"notes"`
}

type CreateOrderResult struct {
	OrderID     int64   `json:"order_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

type CreateOrderReq struct {
	ReferenceID string          `json:"reference_id"`
	TenantID    int64           `json:"tenant_id"`
	TableID     int64           `json:"table_id"`
	UserID      int64           `json:"user_id"`
	Notes       string          `json:"notes"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Status      string          `json:"status"`
	Items       []OrderItemInput
}

type ValidateItemsResult struct {
	Items []CreateOrderItemReq `json:"items"`
	Total string               `json:"total"`
}

func (s OrderStatus) ValidateStatus() bool {
	switch s {
	case OrderStatusServed:
		return true
	}
	return false
}
