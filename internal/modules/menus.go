package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/cache"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"time"
)

func (m *WebModule) CreateMenu(ctx context.Context, req models.CreateMenuReq) (models.Menu, error) {
	response, err := m.persistenceDB.CreateMenu(ctx, db.CreateMenuParams{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: models.ToPGText(req.Description),
		IsActive:    req.IsActive,
	})

	if err != nil {
		return models.Menu{}, err
	}

	// Invalidate menus cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantMenusPattern(req.TenantID))
	result := models.ConvertMenuModel(response)
	return result, nil
}

func (m *WebModule) ListMenus(ctx context.Context, arg models.ListMenusReq) ([]models.Menu, int64, error) {
	// Try cache first
	cacheKey := cache.MenusKey(arg.TenantID)
	start := time.Now()
	var cached struct {
		Menus []models.Menu `json:"menus"`
		Total int64         `json:"total"`
	}
	if err := cache.Get(ctx, cacheKey, &cached); err == nil {
		logger.Log.Info("cache hit for menus", "duration", time.Since(start))
		return cached.Menus, cached.Total, nil
	} else {
		logger.Log.Info("cache miss for menus", "duration", time.Since(start), "error", err.Error())
	}

	// cache miss, hit DB
	response, err := m.persistenceDB.ListMenus(ctx, db.ListMenusParams{
		TenantID: arg.TenantID,
		Limit:    int32(arg.Limit),
		Offset:   int32(arg.Offset),
	})

	if err != nil {
		return nil, 0, err
	}

	total, err := m.persistenceDB.CountMenus(ctx, arg.TenantID)
	if err != nil {
		return nil, 0, err
	}

	result := models.ConvertMenuToModles(response)
	// store in cache
	if err = cache.Set(ctx, cacheKey, struct {
		Menus []models.Menu `json:"menus"`
		Total int64         `json:"total"`
	}{result, total}, cache.TTLMenu); err != nil {
		logger.Log.Error("failed to set menus cache", err)
	} else {
		logger.Log.Info("menus cached successfully", "key", cacheKey)
	}

	return result, total, nil
}

func (m *WebModule) GetMenuByID(ctx context.Context, id, tenatID int64) (models.Menu, error) {
	response, err := m.persistenceDB.GetMenuByID(ctx, db.GetMenuByIDParams{
		ID:       id,
		TenantID: tenatID,
	})
	if err != nil {
		return models.Menu{}, err
	}

	result := models.ConvertMenuModel(response)
	return result, nil
}

func (m *WebModule) UpdateMenu(ctx context.Context, arg models.UpdateMenuReq) (models.Menu, error) {
	response, err := m.persistenceDB.UpdateMenu(ctx, db.UpdateMenuParams{
		ID:          arg.ID,
		TenantID:    arg.TenantID,
		Name:        arg.Name,
		Description: models.ToPGText(arg.Description),
		IsActive:    arg.IsActive,
	})
	if err != nil {
		return models.Menu{}, err
	}

	// Invalidate menus cache for this tenant
	_ = cache.Delete(ctx, cache.TenantMenusPattern(arg.TenantID))
	result := models.ConvertMenuModel(response)
	return result, nil
}

func (m *WebModule) DeleteMenu(ctx context.Context, id, tenantID int64) error {
	err := m.persistenceDB.DeleteMenu(ctx, db.DeleteMenuParams{
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		return err
	}

	// Invalidate menus cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantMenusPattern(tenantID))
	return nil
}

func (m *WebModule) RestoreMenu(ctx context.Context, id, tenantID int64) error {
	err := m.persistenceDB.RestoreMenu(ctx, db.RestoreMenuParams{
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		return err
	}

	// Invalidate menus cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantMenusPattern(tenantID))

	return nil
}
