package models

import (
	db "restaurant-platform/database/sqlc/gen"
	"time"
)

type Tenant struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Slug   string `json:"slug" db:"slug"`
	Status string `json:"status" db:"status"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"                     db:"created_at"`
	UpdatedAt time.Time `json:"updated_at"                     db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}

func ConvertTenantModel(tenant db.Tenant) Tenant {
	return Tenant{
		ID:        int(tenant.ID),
		Name:      tenant.Name,
		Slug:      tenant.Slug,
		Status:    tenant.Status,
		CreatedAt: tenant.CreatedAt.Time,
		UpdatedAt: tenant.UpdatedAt.Time,
		DeletedAt: tenant.DeletedAt.Time,
	}
}

func ConvertTenantModels(tenants []db.Tenant) []Tenant {
	result := make([]Tenant, 0, len(tenants))
	for _, tenant := range tenants {
		result = append(result, ConvertTenantModel(tenant))
	}
	return result
}

type CreateTenantReq struct {
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status string `json:"status"`
}

type ListTenantsReq struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UpdateTenantReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status string `json:"status"`
}
