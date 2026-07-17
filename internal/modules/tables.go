package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/cache"
	"restaurant-platform/internal/models"
)

func (m *WebModule) CreateTable(ctx context.Context, arg models.CreateTableReq) (models.Table, error) {
	response, err := m.persistenceDB.CreateTable(ctx, db.CreateTableParams{
		TenantID: arg.TenantID,
		Name:     arg.Name,
		Capacity: arg.Capacity,
	})
	if err != nil {
		return models.Table{}, err
	}

	// Invalidate tables cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantTablesPattern(arg.TenantID))
	result := models.ConvertTableModel(response)

	return result, nil
}

func (m *WebModule) ListTables(ctx context.Context, arg models.ListTablesReq) ([]models.Table, int64, error) {
	// Try cache first
	cachedKey := cache.TablesKey(arg.TenantID, arg.Column2)
	var cached struct {
		Tables []models.Table `json:"tables"`
		Total  int64          `json:"total"`
	}

	if err := cache.Get(ctx, cachedKey, cached); err == nil {
		return cached.Tables, cached.Total, nil
	}

	// cache miss, hit db
	response, err := m.persistenceDB.ListTables(ctx, db.ListTablesParams{
		TenantID: arg.TenantID,
		Column2:  arg.Column2,
		Limit:    int32(arg.Limit),
		Offset:   int32(arg.Offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := m.persistenceDB.CountTables(ctx, db.CountTablesParams{
		TenantID: arg.TenantID,
		Column2:  arg.Column2,
	})
	if err != nil {
		return nil, 0, err
	}

	result := models.ConvertTableModels(response)

	// set cache for the next fetch
	_ = cache.Set(ctx, cachedKey, struct {
		Tables []models.Table `json:"tables"`
		Total  int64          `json:"total"`
	}{result, total}, cache.TTLTable)

	return result, total, nil
}

func (m *WebModule) GetTableByID(ctx context.Context, id, tenantID int64) (models.Table, error) {
	response, err := m.persistenceDB.GetTableByID(ctx, db.GetTableByIDParams{
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		return models.Table{}, err
	}

	result := models.ConvertTableModel(response)

	return result, nil
}

func (m *WebModule) UpdateTable(ctx context.Context, arg models.UpdateTableReq) (models.Table, error) {
	response, err := m.persistenceDB.UpdateTable(ctx, db.UpdateTableParams{
		Name:     arg.Name,
		Capacity: arg.Capacity,
		ID:       arg.ID,
		TenantID: arg.TenantID,
	})

	if err != nil {
		return models.Table{}, err
	}

	// Invalidate tables cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantTablesPattern(arg.TenantID))
	result := models.ConvertTableModel(response)

	return result, nil
}

func (m *WebModule) DeleteTable(ctx context.Context, id, tenantID int64) error {
	err := m.persistenceDB.DeleteTable(ctx, db.DeleteTableParams{
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		return err
	}

	// Invalidate tables cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantTablesPattern(tenantID))
	return nil
}

func (m *WebModule) UpdateTableStatus(ctx context.Context, id, tenantID int64, status string) (models.Table, error) {
	response, err := m.persistenceDB.UpdateTableStatus(ctx, db.UpdateTableStatusParams{
		Status:   status,
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		return models.Table{}, err
	}

	// Invalidate tables cache for this tenant
	_ = cache.DeleteByPattern(ctx, cache.TenantTablesPattern(tenantID))
	result := models.ConvertTableModel(response)

	return result, nil
}
