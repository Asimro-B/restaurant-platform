package module

import (
	"context"
	"fmt"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/models"

	"github.com/shopspring/decimal"
)

func (m *WebModule) ProcessPayment(ctx context.Context, tenantID, orderID int64, amount decimal.Decimal, arg models.ProcessPaymentReq) (models.Payment, error) {
	response, err := m.persistenceDB.CreatePayment(ctx, db.CreatePaymentParams{
		TenantID:      tenantID,
		OrderID:       orderID,
		Amount:        amount,
		PaymentMethod: arg.PaymentMethod,
		Reference:     models.ToPGText(arg.Reference),
		Notes:         models.ToPGText(arg.Notes),
	})
	if err != nil {
		return models.Payment{}, err
	}

	// update payment status to completed
	completed, err := m.persistenceDB.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		Status:   string(models.PaymentStatusCompleted),
		ID:       response.ID,
		TenantID: tenantID,
	})
	if err != nil {
		return models.Payment{}, err
	}

	result := models.ConvertPaymentModel(completed)
	return result, err
}

func (m *WebModule) GetBill(ctx context.Context, referenceID string, tenantID int64) (*models.BillResponse, error) {
	rows, err := m.persistenceDB.GetOrderWithItems(ctx, db.GetOrderWithItemsParams{
		ReferenceID: models.ToPGText(referenceID),
		TenantID:    tenantID,
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("order not found")
	}

	// build bill
	bill := &models.BillResponse{
		OrderID:     rows[0].OrderID,
		ReferenceID: rows[0].ReferenceID.String,
		TableID:     rows[0].TableID,
		Status:      rows[0].Status,
		TotalAmount: rows[0].TotalAmount,
		CreatedAt:   rows[0].CreatedAt.Time,
		Items:       make([]models.BillItem, 0, len(rows)),
	}

	for _, row := range rows {
		bill.Items = append(bill.Items, models.BillItem{
			OrderItemID:  row.OrderItemID,
			MenuItemName: row.MenuItemName,
			Quantity:     int(row.Quantity),
			UnitPrice:    row.UnitPrice,
			TotalPrice:   row.TotalPrice,
			Notes:        row.ItemNotes.String,
		})
	}

	return bill, nil
}
