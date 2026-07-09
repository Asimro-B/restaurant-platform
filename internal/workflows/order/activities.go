package orderworkflow

import (
	"context"
	"fmt"
	"restaurant-platform/internal/models"

	"github.com/shopspring/decimal"
)

type OrderActivities struct {
	module OrderModule // interface — keeps activities decoupled from concrete module
}

// interface so activities don't import the full module package
type OrderModule interface {
	GetMenuItemByID(ctx context.Context, id, tenantID int64) (*models.MenuItem, error)
	CreateOrderWithItems(ctx context.Context, order models.CreateOrderReq, items []models.CreateOrderItemReq) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID, tenantID int64, status models.OrderStatus) error
}

func NewOrderActivities(module OrderModule) *OrderActivities {
	return &OrderActivities{module: module}
}

func (a *OrderActivities) ValidateAndPriceItems(ctx context.Context, input models.CreateOrderInput) (*models.ValidateItemsResult, error) {
	var items []models.CreateOrderItemReq
	total := decimal.Zero

	for _, item := range input.Items {
		menuItem, err := a.module.GetMenuItemByID(ctx, item.MenuItemID, input.TenantID)
		if err != nil {
			return nil, fmt.Errorf("menu item %d not found", item.MenuItemID)
		}
		if !menuItem.IsAvailable {
			return nil, fmt.Errorf("menu item %s is not available", menuItem.Name)
		}

		unitPrice := menuItem.Price
		lineTotal := unitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))

		total = total.Add(lineTotal)

		items = append(items, models.CreateOrderItemReq{
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
			UnitPrice:  unitPrice,
			TotalPrice: lineTotal,
			Notes:      item.Notes,
		})
	}

	return &models.ValidateItemsResult{Items: items, Total: total.String()}, nil
}

func (a *OrderActivities) PersistOrder(ctx context.Context, input models.CreateOrderInput, result *models.ValidateItemsResult) (*models.Order, error) {
	total, err := decimal.NewFromString(result.Total)
	if err != nil {
		return nil, fmt.Errorf("invalid total amount: %w", err)
	}

	return a.module.CreateOrderWithItems(ctx, models.CreateOrderReq{
		ReferenceID: input.ReferenceID,
		TenantID:    input.TenantID,
		TableID:     input.TableID,
		UserID:      input.UserID,
		Notes:       input.Notes,
		TotalAmount: total,
	}, result.Items)
}

func (a *OrderActivities) ConfirmOrder(ctx context.Context, orderID, tenantID int64) error {
	return a.module.UpdateOrderStatus(ctx, orderID, tenantID, models.OrderStatusConfirmed)
}

func (a *OrderActivities) MarkPreparing(ctx context.Context, orderID, tenantID int64) error {
	return a.module.UpdateOrderStatus(ctx, orderID, tenantID, models.OrderStatusPreparing)
}

func (a *OrderActivities) MarkReady(ctx context.Context, orderID, tenantID int64) error {
	return a.module.UpdateOrderStatus(ctx, orderID, tenantID, models.OrderStatusReady)
}

func (a *OrderActivities) MarkServed(ctx context.Context, orderID, tenantID int64) error {
	return a.module.UpdateOrderStatus(ctx, orderID, tenantID, models.OrderStatusServed)
}
