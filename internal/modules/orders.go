package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/models"
)

func (m *WebModule) CreateOrderWithItems(
	ctx context.Context,
	params models.CreateOrderReq,
	items []models.CreateOrderItemReq,
) (*models.Order, error) {

	// 1. Insert the order row
	order, err := m.persistenceDB.CreateOrder(ctx, db.CreateOrderParams{
		TenantID:    params.TenantID,
		TableID:     params.TableID,
		UserID:      params.UserID,
		Notes:       models.ToPGText(params.Notes),
		TotalAmount: params.TotalAmount,
		ReferenceID: models.ToPGText(params.ReferenceID),
	})
	if err != nil {
		return nil, err
	}

	// 2. Insert each order item using the new order's ID
	for _, item := range items {
		_, err := m.persistenceDB.CreateOrderItem(ctx, db.CreateOrderItemParams{
			TenantID:   params.TenantID,
			OrderID:    order.ID,
			MenuItemID: item.MenuItemID,
			Quantity:   int32(item.Quantity),
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
			Notes:      models.ToPGText(item.Notes),
		})
		if err != nil {
			return nil, err
		}
	}

	result := models.ConvertOrderModel(order)
	return &result, nil
}

func (m *WebModule) ListOrders(ctx context.Context, arg models.ListOrdersReq) ([]models.Order, int64, error) {
	response, err := m.persistenceDB.ListOrders(ctx, db.ListOrdersParams{
		TenantID: arg.TenantID,
		Column2:  arg.Status,
		Column3:  arg.TableID,
		Limit:    int32(arg.Limit),
		Offset:   int32(arg.Offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := m.persistenceDB.CountOrders(ctx, db.CountOrdersParams{
		TenantID: arg.TenantID,
		Column2:  arg.Status,
		Column3:  arg.TableID,
	})
	if err != nil {
		return nil, 0, err
	}

	result := models.ConvertOrderModels(response)

	return result, total, nil
}

func (m *WebModule) GetOrderByID(ctx context.Context, id, tenantID, tableID, userID int64) (models.Order, error) {
	response, err := m.persistenceDB.GetOrderByID(ctx, db.GetOrderByIDParams{
		ID:       id,
		TenantID: tenantID,
		TableID:  tableID,
		UserID:   userID,
	})
	if err != nil {
		return models.Order{}, err
	}

	result := models.ConvertOrderModel(response)

	return result, nil
}

func (m *WebModule) UpdateOrderStatus(ctx context.Context, orderID, tenantID int64, status models.OrderStatus) error {
	_, err := m.persistenceDB.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		Status:   string(status),
		ID:       orderID,
		TenantID: tenantID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *WebModule) DeleteOrder(ctx context.Context, id, tenantID, tableID, userID int64) error {
	err := m.persistenceDB.DeleteOrder(ctx, db.DeleteOrderParams{
		ID:       id,
		TenantID: tenantID,
		TableID:  tableID,
		UserID:   userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *WebModule) GetOrderByReferenceID(ctx context.Context, referenceID string, tenantID int64) (models.Order, error) {
	response, err := m.persistenceDB.GetOrderByReferenceID(ctx, db.GetOrderByReferenceIDParams{
		ReferenceID: models.ToPGText(referenceID),
		TenantID:    tenantID,
	})
	if err != nil {
		return models.Order{}, err
	}

	result := models.ConvertOrderModel(response)

	return result, err
}
