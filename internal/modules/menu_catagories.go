package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/cache"
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

	// Invalidate menu categories cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantMenuCategoriesPattern(arg.TenantID))
	result := models.ConvertMenuCategoryModel(response)
	return result, nil
}

func (m *WebModule) ListMenuCategories(ctx context.Context, arg models.ListMenuCategoriesReq) ([]models.MenuCategory, int64, error) {
	// Try cache first
	cacheKey := cache.MecnuCategoriesKey(arg.TenantID, arg.MenuID)
	var cached struct {
		MenuCategories []models.MenuCategory `json:"menu_categories"`
		Total          int64                 `json:"total"`
	}
	if err := cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached.MenuCategories, cached.Total, nil
	}

	// Cache miss hit db
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

	// store in cache
	_ = cache.Set(ctx, cacheKey, struct {
		MenuCategoies []models.MenuCategory `json:"menu_categories"`
		Total         int64                 `json:"total"`
	}{result, total}, cache.TTLMenu)

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

	// Invalidate menu categories cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantMenuCategoriesPattern(arg.TenantID))
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

	// Invalidate menu categories cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantMenuCategoriesPattern(tenantID))
	return nil
}

func (m *WebModule) RestoreMenuCategory(ctx context.Context, id, menuID, tenantID int64) error {
	err := m.persistenceDB.RestoreMenuCategory(ctx, db.RestoreMenuCategoryParams{
		ID:       id,
		MenuID:   menuID,
		TenantID: tenantID,
	})
	if err != nil {
		return err
	}

	// Invalidate menu categories cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantMenuCategoriesPattern(tenantID))
	return nil
}
