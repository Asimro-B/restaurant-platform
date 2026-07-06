package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/models"
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
	result := models.ConvertMenuModel(response)
	return result, nil
}

func (m *WebModule) ListMenus(ctx context.Context, arg models.ListMenusReq) ([]models.Menu, int64, error) {
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
	return nil
}
