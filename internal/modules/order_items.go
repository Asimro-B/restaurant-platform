package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/models"
	pusherclient "restaurant-platform/internal/pusher"
)

func (m *WebModule) CreateorderItem(ctx context.Context, arg models.CreateOrderItemReq) (models.OrderItem, error) {
	response, err := m.persistenceDB.CreateOrderItem(ctx, db.CreateOrderItemParams{
		TenantID:   arg.TenantID,
		OrderID:    arg.OrderID,
		MenuItemID: arg.MenuItemID,
		Quantity:   int32(arg.Quantity),
		UnitPrice:  arg.UnitPrice,
		TotalPrice: arg.TotalPrice,
		Notes:      models.ToPGText(arg.Notes),
	})
	if err != nil {
		return models.OrderItem{}, err
	}

	result := models.ConvertOrderItemModel(response)

	return result, nil
}

func (m *WebModule) NotifyKitchen(ctx context.Context, orderID, tenantID int64) error {
	return pusherclient.Publish(
		pusherclient.KitchenChannel(tenantID),
		"order.Created",
		map[string]interface{}{
			"order_id":  orderID,
			"tenant_id": tenantID,
			"message":   "new order recieved",
		},
	)
}

func (m *WebModule) NotifyWaiter(ctx context.Context, orderID, tenantID int64) error {
	return pusherclient.Publish(
		pusherclient.KitchenChannel(tenantID),
		"ordr.ready",
		map[string]interface{}{
			"order_id":  orderID,
			"tenant_id": tenantID,
			"message":   "order is ready for picking",
		},
	)
}
