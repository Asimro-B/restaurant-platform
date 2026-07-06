package models

import (
	db "restaurant-platform/database/sqlc/gen"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID          int64
	TenantID    int64
	TableID     int64
	UserID      int64
	Status      string
	Notes       string
	TotalAmount decimal.Decimal
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateOrderReq struct {
	TenantID    int64           `json:"tenant_id"`
	TableID     int64           `json:"table_id"`
	UserID      int64           `json:"user_id"`
	Notes       string          `json:"notes"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Status      string          `json:"status"`
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
