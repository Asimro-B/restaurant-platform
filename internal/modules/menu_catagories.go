package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/models"
)

func (m *WebModule) CreateMenuCategory(ctx context.Context, arg models.CreateMenuCategoryReq) (models.MenuCategory, error) {
	response, err := m.persistenceDB.CreateMenuCategory(ctx, db.CreateMenuCategoryParams{
		TenantID:    arg.TenantID,
		MenuID:      arg.MenuID,
		Name:        arg.Name,
		Description: models.ToPGText(arg.Description),
		SortOrder:   int32(arg.SortOrder),
		IsActive:    arg.IsActive,
	})
	if err != nil {
		return models.MenuCategory{}, err
	}

	result := models.ConvertMenuCategoryModel(response)
	return result, nil
}

func (m *WebModule) ListMenuCategories(ctx context.Context, arg models.ListMenuCategoriesReq) ([]models.MenuCategory, int64, error) {
	response, err := m.persistenceDB.ListMenuCategories(ctx, db.ListMenuCategoriesParams{
		TenantID: arg.TenantID,
		MenuID:   arg.MenuID,
		Limit:    int32(arg.Limit),
		Offset:   int32(arg.Offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := m.persistenceDB.CountMenuCategories(ctx, db.CountMenuCategoriesParams{
		MenuID:   arg.MenuID,
		TenantID: arg.TenantID,
	})
	if err != nil {
		return nil, 0, err
	}

	result := models.ConvertMenuCategoryCategoryToModles(response)

	return result, total, nil
}

func (m *WebModule) GetMenuCategoryByID(ctx context.Context, tenantID, menuID, id int64) (models.MenuCategory, error) {
	response, err := m.persistenceDB.GetMenuCategoryByID(ctx, db.GetMenuCategoryByIDParams{
		TenantID: tenantID,
		MenuID:   menuID,
		ID:       id,
	})
	if err != nil {
		return models.MenuCategory{}, err
	}

	result := models.ConvertMenuCategoryModel(response)

	return result, nil
}

func (m *WebModule) UpdateMenuCategory(ctx context.Context, arg models.UpdateMenuCategoryReq) (models.MenuCategory, error) {
	response, err := m.persistenceDB.UpdateMenuCategory(ctx, db.UpdateMenuCategoryParams{
		ID:          arg.ID,
		MenuID:      arg.MenuID,
		TenantID:    arg.TenantID,
		Name:        arg.Name,
		Description: models.ToPGText(arg.Description),
		SortOrder:   int32(arg.SortOrder),
		IsActive:    arg.IsActive,
	})
	if err != nil {
		return models.MenuCategory{}, err
	}

	result := models.ConvertMenuCategoryModel(response)

	return result, nil
}

func (m *WebModule) DeleteMenuCategory(ctx context.Context, id, menuID, tenantID int64) error {
	err := m.persistenceDB.DeleteMenuCategory(ctx, db.DeleteMenuCategoryParams{
		ID:       id,
		MenuID:   menuID,
		TenantID: tenantID,
	})
	if err != nil {
		return err
	}

	return nil
}
