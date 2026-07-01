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
