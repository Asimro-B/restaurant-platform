package models

import (
	db "restaurant-platform/database/sqlc/gen"
	"time"

	"github.com/shopspring/decimal"
)

type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodCard     PaymentMethod = "card"
	PaymentMethodTeleBirr PaymentMethod = "telebirr"
)

func (p PaymentMethod) IsValid() bool {
	switch p {
	case PaymentMethodCash, PaymentMethodCard, PaymentMethodTeleBirr:
		return true
	}
	return false
}

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID            int64           `json:"id" db:"id"`
	TenantID      int64           `json:"tenant_id" db:"tenant_id"`
	OrderID       int64           `json:"order_id" db:"order_id"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
	Status        string          `json:"status" db:"status"`
	PaymentMethod string          `json:"payment_method" db:"payment_method"`
	Reference     string          `json:"reference"` // TeleBirr transaction ID, card ref
	Notes         string          `json:"notes" db:"notes"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

func ConvertPaymentModel(payment db.Payment) Payment {
	return Payment{
		ID:            payment.ID,
		TenantID:      payment.TenantID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status,
		PaymentMethod: payment.PaymentMethod,
		Reference:     payment.Reference.String,
		Notes:         payment.Notes.String,
		CreatedAt:     payment.CreatedAt.Time,
		UpdatedAt:     payment.UpdatedAt.Time,
	}
}

func ConvertPaymentModels(payments []db.Payment) []Payment {
	result := make([]Payment, 0, len(payments))
	for _, payment := range payments {
		result = append(result, ConvertPaymentModel(payment))
	}
	return result
}

type BillResponse struct {
	OrderID     int64           `json:"order_id"`
	ReferenceID string          `json:"reference_id"`
	TableID     int64           `json:"table_id"`
	Status      string          `json:"status"`
	Items       []BillItem      `json:"items"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	CreatedAt   time.Time       `json:"created_at"`
}

type BillItem struct {
	OrderItemID  int64           `json:"order_item_id"`
	MenuItemName string          `json:"menu_item_name"`
	Quantity     int             `json:"quantity"`
	UnitPrice    decimal.Decimal `json:"unit_price"`
	TotalPrice   decimal.Decimal `json:"total_price"`
	Notes        string          `json:"notes"`
}

type ProcessPaymentReq struct {
	TableID       int64  `json:"table_id"`
	PaymentMethod string `json:"payment_method"` // cash | card | telebirr
	Reference     string `json:"reference"`      // TeleBirr transaction ID, card ref
	Notes         string `json:"notes"`
}
