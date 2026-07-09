package models

import (
	db "restaurant-platform/database/sqlc/gen"
	"time"

	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID         int64           `json:"id" db:"id"`
	TenantID   int64           `json:"tenant_id" db:"tenant_id"`
	OrderID    int64           `json:"order_id" db:"order_id"`
	MenuItemID int64           `json:"menu_item_id" db:"menu_item_id"`
	Quantity   int             `json:"quantity" db:"quantity"`
	UnitPrice  decimal.Decimal `unit_price:"id" db:"unit_price"`
	TotalPrice decimal.Decimal `json:"total_price" db:"total_price"`
	Notes      string          `json:"notes" db:"notes"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

type CreateOrderItemReq struct {
	TenantID   int64           `json:"tenant_id"`
	OrderID    int64           `json:"order_id"`
	MenuItemID int64           `json:"menu_item_id"`
	Quantity   int             `json:"quantity"`
	UnitPrice  decimal.Decimal `unit_price:"id"`
	TotalPrice decimal.Decimal `json:"total_price"`
	Notes      string          `json:"notes"`
}

func ConvertOrderItemModel(orderItem db.OrderItem) OrderItem {
	return OrderItem{
		ID:         orderItem.ID,
		TenantID:   orderItem.TenantID,
		OrderID:    orderItem.OrderID,
		MenuItemID: orderItem.MenuItemID,
		Quantity:   int(orderItem.Quantity),
		UnitPrice:  orderItem.UnitPrice,
		TotalPrice: orderItem.TotalPrice,
		Notes:      orderItem.Notes.String,
		CreatedAt:  orderItem.CreatedAt.Time,
	}
}

func ConvertOrderItemModels(orderItems []db.OrderItem) []OrderItem {
	result := make([]OrderItem, 0, len(orderItems))
	for _, orderItem := range orderItems {
		result = append(result, ConvertOrderItemModel(orderItem))
	}
	return result
}
