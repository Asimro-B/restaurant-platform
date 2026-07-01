package module

import (
	"context"
	db "restaurant-platform/database/sqlc/gen"
	"restaurant-platform/internal/models"
)

func (m *WebModule) CreateUser(ctx context.Context, req models.CreateUserReq) (models.User, error) {
	response, err := m.persistenceDB.CreateUser(ctx, db.CreateUserParams{
		TenantID:     int64(req.TenantID),
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
		Role:         req.Role,
		FirstName:    models.ToPGText(req.FirstName),
		LastName:     models.ToPGText(req.LastName),
		Location:     models.ToPGText(req.Location),
		MobilePhone:  models.ToPGText(req.MobilePhone),
		Phone:        models.ToPGText(req.Phone),
	})
	if err != nil {
		return models.User{}, err
	}

	result := models.ConvertUserModel(response)
	return result, nil
}

func (m *WebModule) GetUserByEmail(ctx context.Context, req models.GetUserByEmailReq) (models.User, error) {
	response, err := m.persistenceDB.GetUserByEmail(ctx, db.GetUserByEmailParams{
		Email:    req.Email,
		TenantID: int64(req.TenantID),
	})
	if err != nil {
		return models.User{}, err
	}

	result := models.ConvertUserModel(response)
	return result, nil
}
