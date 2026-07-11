package models

import "time"

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusCompleted ReservationStatus = "completed"
)

func (s ReservationStatus) IsValid() bool {
	switch s {
	case ReservationStatusPending, ReservationStatusConfirmed, ReservationStatusCancelled, ReservationStatusCompleted:
		return true
	}
	return false
}

type Reservation struct {
	ID            int64             `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID      int64             `json:"tenant_id" gorm:"not null"`
	TableID       int64             `json:"table_id" gorm:"not null"`
	CustomerName  string            `json:"customer_name" gorm:"not null"`
	CustomerPhone string            `json:"customer_phone" gorm:"not null"`
	PartySize     int               `json:"party_size" gorm:"not null"`
	ReservedAt    time.Time         `json:"reserved_at" gorm:"not null"`
	Status        ReservationStatus `json:"status" gorm:"default:pending"`
	Notes         string            `json:"notes"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type CreateReservationReq struct {
	TableID       int64     `json:"table_id"`
	CustomerName  string    `json:"customer_name"`
	CustomerPhone string    `json:"customer_phone"`
	PartySize     int       `json:"party_size"`
	ReservedAt    time.Time `json:"reserved_at"`
	Notes         string    `json:"notes"`
}

type ListReservationsReq struct {
	TenantID int64
	Status   string
	Page     int
	Limit    int
	Offset   int
}
