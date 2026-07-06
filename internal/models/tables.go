package models

import (
	db "restaurant-platform/database/sqlc/gen"
	"time"
)

type TableStatus string

const (
	TableStatusAvailable TableStatus = "available"
	TableStatusOccupied  TableStatus = "occupied"
	TableStatusReserved  TableStatus = "reserved"
	TableStatusInactive  TableStatus = "inactive"
)

func (s TableStatus) IsValid() bool {
	switch s {
	case TableStatusAvailable, TableStatusOccupied, TableStatusReserved, TableStatusInactive:
		return true
	}
	return false
}

type Table struct {
	ID        int64     `json:"id" db:"id"`
	TenantID  int64     `json:"tenant_id" db:"tenant_id"`
	Name      string    `json:"name" db:"name"`
	Capacity  int32     `json:"capacity" db:"capacity"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTableReq struct {
	TenantID int64  `json:"tenant_id"`
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
}

func ConvertTableModel(table db.Table) Table {
	return Table{
		ID:        table.ID,
		TenantID:  table.TenantID,
		Name:      table.Name,
		Capacity:  table.Capacity,
		Status:    table.Status,
		CreatedAt: table.CreatedAt.Time,
		UpdatedAt: table.UpdatedAt.Time,
	}
}

func ConvertTableModels(tables []db.Table) []Table {
	result := make([]Table, 0, len(tables))
	for _, table := range tables {
		result = append(result, ConvertTableModel(table))
	}
	return result
}

type ListTablesReq struct {
	TenantID int64  `json:"tenant_id"`
	Column2  string `json:"column_2"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

type UpdateTableReq struct {
	Name     string `json:"name"`
	Capacity int32  `json:"capacity"`
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
}
