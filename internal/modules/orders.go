package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/models"
)

func (m *WebModule) CreateOrder(ctx context.Context, arg models.CreateOrderReq) (models.Order, error) {
	response, err := m.persistenceDB.CreateOrder(ctx, db.CreateOrderParams{
		TenantID:    arg.TenantID,
		TableID:     arg.TableID,
		UserID:      arg.UserID,
		Notes:       models.ToPGText(arg.Notes),
		TotalAmount: arg.TotalAmount,
		Status:      arg.Status,
	})
	if err != nil {
		return models.Order{}, err
	}

	result := models.ConvertOrderModel(response)

	return result, nil
}

func (m *WebModule) ListOrders(ctx context.Context, arg models.ListOrdersReq) ([]models.Order, int64, error) {
	response, err := m.persistenceDB.ListOrders(ctx, db.ListOrdersParams{
		TenantID: arg.TenantID,
		TableID:  arg.TableID,
		UserID:   arg.UserID,
		Limit:    int32(arg.Limit),
		Offset:   int32(arg.Offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := m.persistenceDB.CountOrders(ctx, db.CountOrdersParams{
		TenantID: arg.TenantID,
		TableID:  arg.TableID,
		UserID:   arg.UserID,
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

func (m *WebModule) UpdateOrderStatus(ctx context.Context, arg models.UpdateOrderStatusReq) (models.Order, error) {
	response, err := m.persistenceDB.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		Status:   arg.Status,
		ID:       arg.ID,
		TenantID: arg.TenantID,
		TableID:  arg.TableID,
		UserID:   arg.UserID,
	})
	if err != nil {
		return models.Order{}, err
	}

	result := models.ConvertOrderModel(response)
	return result, nil
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
