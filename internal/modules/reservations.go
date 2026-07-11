package module

import (
	"context"
	"fmt"
	"restaurant-platform/internal/models"
	"time"
)

func (m *WebModule) CreateReservation(ctx context.Context, tenantID int64, req models.CreateReservationReq) (*models.Reservation, error) {
	// check table availability
	var count int64
	if err := m.persistenceDB.GormDB.WithContext(ctx).Model(&models.Reservation{}).Where("table_id = ? AND tenant_id = ? AND status IN ? AND reserved_at BETWEEN ? AND ?", req.TableID, tenantID, []string{"pending", "confirmed"}, req.ReservedAt.Add(-2*time.Hour), req.ReservedAt.Add(2*time.Hour)).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check table availabitlity: %w", err)
	}

	if count > 0 {
		return nil, fmt.Errorf("table is already reserved at this time")
	}

	// create reservation
	reservation := &models.Reservation{
		TenantID:      tenantID,
		TableID:       req.TableID,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		PartySize:     req.PartySize,
		ReservedAt:    req.ReservedAt,
		Status:        models.ReservationStatusPending,
		Notes:         req.Notes,
	}

	if err := m.persistenceDB.GormDB.WithContext(ctx).Create(reservation).Error; err != nil {
		return nil, err
	}

	// update table status to reserved
	if err := m.persistenceDB.GormDB.WithContext(ctx).Model(&models.Table{}).Where("id = ? AND tenant_id = ?", req.TableID, tenantID).Update("status", "reserved").Error; err != nil {
		return nil, err
	}

	return reservation, nil
}

func (m *WebModule) ListReservations(ctx context.Context, req models.ListReservationsReq) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	query := m.persistenceDB.GormDB.WithContext(ctx).Model(&models.Reservation{}).Where("tenant_id = ?", req.TenantID)

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)
	err := query.Order("reserved_at ASC").Limit(req.Limit).Offset(req.Offset).Find(&reservations).Error

	return reservations, total, err
}

func (m *WebModule) GetReservationByID(ctx context.Context, id, tenantID int64) (*models.Reservation, error) {
	var reservation models.Reservation

	err := m.persistenceDB.GormDB.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&reservation).Error
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (m *WebModule) UpdateReservationStatus(ctx context.Context, id, tenantID int64, status models.ReservationStatus) (*models.Reservation, error) {
	var reservation models.Reservation
	err := m.persistenceDB.GormDB.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&reservation).Error
	if err != nil {
		return nil, err
	}

	reservation.Status = status

	if err := m.persistenceDB.GormDB.WithContext(ctx).Save(&reservation).Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}
