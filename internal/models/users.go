package models

import (
	db "restaurant-platform/database/sqlc/gen"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           int    `json:"id"            db:"id"`
	TenantID     int    `json:"tenant_id"     db:"tenant_id"`
	Email        string `json:"email"         db:"email"`
	PasswordHash string `json:"-"             db:"password_hash"`
	Role         string `json:"role"          db:"role"`
	FirstName    string `json:"first_name"    db:"first_name"`
	LastName     string `json:"last_name"     db:"last_name"`
	Location     string `json:"location"      db:"location"`
	MobilePhone  string `json:"mobile_phone"  db:"mobile_phone"`
	Phone        string `json:"phone"         db:"phone"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func ConvertUserModel(user db.User) User {
	return User{
		ID:           int(user.ID),
		TenantID:     int(user.TenantID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		FirstName:    user.FirstName.String,
		LastName:     user.LastName.String,
		Location:     user.Location.String,
		MobilePhone:  user.MobilePhone.String,
		Phone:        user.Phone.String,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}
}

type CreateUserReq struct {
	TenantID     int64  `json:"tenant_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Location     string `json:"location"`
	MobilePhone  string `json:"mobile_phone"`
	Phone        string `json:"phone"`
}

func ToPGText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  value != "",
	}
}

type GetUserByEmailReq struct {
	Email    string `json:"email"`
	TenantID int    `json:"tenant_id"`
}
