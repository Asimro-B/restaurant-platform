package module

import (
	"context"

	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/models"
)

func (m *WebModule) CreateMenuItem(ctx context.Context, arg models.CreateMenuItemReq) (models.MenuItem, error) {
	response, err := m.persistenceDB.CreateMenuItem(ctx, db.CreateMenuItemParams{
		TenantID:    arg.TenantID,
		CategoryID:  arg.CategoryID,
		MenuID:      arg.MenuID,
		Name:        arg.Name,
		Description: models.ToPGText(arg.Description),
		Price:       arg.Price,
		IsAvailable: arg.IsAvailable,
	})
	if err != nil {
		return models.MenuItem{}, err
	}

	result := models.ConvertMenuItemModel(response)
	return result, nil
}

func (m *WebModule) ListMenuItems(ctx context.Context, arg models.ListMenuItemsReq) ([]models.MenuItem, int64, error) {
	response, err := m.persistenceDB.ListMenuItems(ctx, db.ListMenuItemsParams{
		TenantID:   arg.TenantID,
		CategoryID: arg.CategoryID,
		MenuID:     arg.MenuID,
		Limit:      int32(arg.Limit),
		Offset:     int32(arg.Offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := m.persistenceDB.CountMenuItems(ctx, db.CountMenuItemsParams{
		TenantID:   arg.TenantID,
		CategoryID: arg.CategoryID,
		MenuID:     arg.MenuID,
	})
	if err != nil {
		return nil, 0, err
	}

	result := models.ConvertMenuItemsToModels(response)

	return result, total, nil
}

func (m *WebModule) GetMenuItemByID(ctx context.Context, tenantID, menuID, categoryID, id int64) (models.MenuItem, error) {
	response, err := m.persistenceDB.GetMenuItemByID(ctx, db.GetMenuItemByIDParams{
		TenantID:   tenantID,
		MenuID:     menuID,
		CategoryID: categoryID,
		ID:         id,
	})
	if err != nil {
		return models.MenuItem{}, err
	}

	result := models.ConvertMenuItemModel(response)

	return result, nil
}

func (m *WebModule) UpdateMenuItem(ctx context.Context, arg models.UpdateMenuItemReq) (models.MenuItem, error) {
	response, err := m.persistenceDB.UpdateMenuItem(ctx, db.UpdateMenuItemParams{
		ID:          arg.ID,
		CategoryID:  arg.CategoryID,
		MenuID:      arg.MenuID,
		TenantID:    arg.TenantID,
		Name:        arg.Name,
		Description: models.ToPGText(arg.Description),
		Price:       arg.Price,
		IsAvailable: arg.IsAvailable,
	})
	if err != nil {
		return models.MenuItem{}, err
	}

	result := models.ConvertMenuItemModel(response)

	return result, nil
}

func (m *WebModule) DeleteMenuItem(ctx context.Context, id, categoryID, menuID, tenantID int64) error {
	err := m.persistenceDB.DeleteMenuItem(ctx, db.DeleteMenuItemParams{
		ID:         id,
		CategoryID: categoryID,
		MenuID:     menuID,
		TenantID:   tenantID,
	})
	if err != nil {
		return err
	}

	return nil
}
