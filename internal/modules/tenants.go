package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	persistencedb "restaurant-platform/database/sqlc/persistenceDB"
	"restaurant-platform/internal/models"
)

type WebModule struct {
	persistenceDB *persistencedb.PersistenceDB
}

func NewModule(persistenceDB *persistencedb.PersistenceDB) *WebModule {
	return &WebModule{
		persistenceDB: persistenceDB,
	}
}

func (m *WebModule) CreateTenant(ctx context.Context, tenant models.CreateTenantReq) (models.Tenant, error) {
	// create tenant on the persistence db
	response, err := m.persistenceDB.CreateTenant(ctx, db.CreateTenantParams{
		Name:   tenant.Name,
		Slug:   tenant.Slug,
		Status: tenant.Status,
	})
	if err != nil {
		return models.Tenant{}, err
	}

	result := models.ConvertTenantModel(response)
	return result, nil
}

func (m *WebModule) ListTenants(ctx context.Context, req models.ListTenantsReq) ([]models.Tenant, int64, error) {
	response, err := m.persistenceDB.ListTenants(ctx, db.ListTenantsParams{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := m.persistenceDB.CountTenants(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := models.ConvertTenantModels(response)
	return result, total, nil
}

func (m *WebModule) GetTenantByID(ctx context.Context, id int64) (models.Tenant, error) {
	response, err := m.persistenceDB.GetTenantByID(ctx, id)
	if err != nil {
		return models.Tenant{}, err
	}

	result := models.ConvertTenantModel(response)
	return result, nil
}

func (m *WebModule) GetTenantBySlug(ctx context.Context, slug string) (models.Tenant, error) {
	response, err := m.persistenceDB.GetTenantBySlug(ctx, slug)
	if err != nil {
		return models.Tenant{}, err
	}

	result := models.ConvertTenantModel(response)
	return result, nil
}

func (m *WebModule) UpdateTenant(ctx context.Context, req models.UpdateTenantReq) (models.Tenant, error) {
	response, err := m.persistenceDB.UpdateTenant(ctx, db.UpdateTenantParams{
		ID:     req.ID,
		Name:   req.Name,
		Slug:   req.Slug,
		Status: req.Status,
	})
	if err != nil {
		return models.Tenant{}, err
	}

	result := models.ConvertTenantModel(response)
	return result, nil
}

func (m *WebModule) DeleteTenant(ctx context.Context, id int64) (models.Tenant, error) {
	response, err := m.persistenceDB.DeleteTenant(ctx, id)
	if err != nil {
		return models.Tenant{}, err
	}

	result := models.ConvertTenantModel(response)
	return result, nil
}

func (m *WebModule) RestoreTenant(ctx context.Context, id int64) (models.Tenant, error) {
	response, err := m.persistenceDB.RestoreTenant(ctx, id)
	if err != nil {
		return models.Tenant{}, err
	}

	result := models.ConvertTenantModel(response)
	return result, nil
}
