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

func (m *WebModule) GetUserByEmail(ctx context.Context, email string, tenantID int64) (models.User, error) {
	response, err := m.persistenceDB.GetUserByEmail(ctx, db.GetUserByEmailParams{
		Email:    email,
		TenantID: tenantID,
	})
	if err != nil {
		return models.User{}, err
	}

	result := models.ConvertUserModel(response)
	return result, nil
}

func (m *WebModule) ListUsers(ctx context.Context, req models.ListUsersReq) ([]models.User, error) {
	response, err := m.persistenceDB.ListUsers(ctx, db.ListUsersParams{
		TenantID:    req.TenantID,
		Search:      req.Search,
		Role:        req.Role,
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
		OffsetCount: req.OffsetCount,
		LimitCount:  req.LimitCount,
	})
	if err != nil {
		return nil, err
	}

	result := models.ConvertUserModels(response)
	return result, nil
}

func (m *WebModule) GetUserByID(ctx context.Context, id int64, tenantID int64) (models.User, error) {
	response, err := m.persistenceDB.GetUserByID(ctx, db.GetUserByIDParams{
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		return models.User{}, err
	}

	result := models.ConvertUserModel(response)
	return result, nil
}

func (m *WebModule) UpdateUser(ctx context.Context, req models.UpdateUserReq) (models.User, error) {
	response, err := m.persistenceDB.UpdateUser(ctx, db.UpdateUserParams{
		ID:           req.ID,
		TenantID:     req.TenantID,
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

func (m *WebModule) DeleteUser(ctx context.Context, id int64, tenantID int64) error {
	err := m.persistenceDB.DeleteUser(ctx, db.DeleteUserParams{
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		return err
	}
	return nil
}
