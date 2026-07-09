package models

import (
	db "restaurant-platform/database/sqlc/gen"
	"time"

	"github.com/shopspring/decimal"
)

type Menu struct {
	ID          int    `json:"id"            db:"id"`
	TenantID    int    `json:"tenant_id"     db:"tenant_id"`
	Name        string `json:"name"         db:"name"`
	Description string `json:"description"             db:"description"`
	IsActive    bool   `json:"is_active"          db:"is_active"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func ConvertMenuModel(menu db.Menu) Menu {
	return Menu{
		ID:          int(menu.ID),
		TenantID:    int(menu.TenantID),
		Name:        menu.Name,
		Description: menu.Description.String,
		IsActive:    menu.IsActive,
		CreatedAt:   menu.CreatedAt.Time,
		UpdatedAt:   menu.UpdatedAt.Time,
	}
}

func ConvertMenuToModles(menus []db.Menu) []Menu {
	result := make([]Menu, 0, len(menus))
	for _, menu := range menus {
		result = append(result, ConvertMenuModel(menu))
	}

	return result
}

type CreateMenuReq struct {
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type ListMenusReq struct {
	TenantID int64 `json:"tenant_id"`
	Page     int   `json:"page"`
	Limit    int   `json:"limit"`
	Offset   int   `json:"offset"`
}

type UpdateMenuReq struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type MenuCategory struct {
	ID          int    `json:"id"            db:"id"`
	MenuID      int64  `json:"menu_id" db:"menu_id"`
	TenantID    int    `json:"tenant_id"     db:"tenant_id"`
	Name        string `json:"name"         db:"name"`
	Description string `json:"description"             db:"description"`
	IsActive    bool   `json:"is_active"          db:"is_active"`
	SortOrder   int    `json:"sort_order" db:"sort_order"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateMenuCategoryReq struct {
	TenantID    int64  `json:"tenant_id"`
	MenuID      int64  `json:"menu_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

func ConvertMenuCategoryModel(menuCategory db.MenuCategory) MenuCategory {
	return MenuCategory{
		ID:          int(menuCategory.ID),
		TenantID:    int(menuCategory.TenantID),
		MenuID:      menuCategory.MenuID,
		Name:        menuCategory.Name,
		Description: menuCategory.Description.String,
		IsActive:    menuCategory.IsActive,
		CreatedAt:   menuCategory.CreatedAt.Time,
		UpdatedAt:   menuCategory.UpdatedAt.Time,
	}
}

func ConvertMenuCategoryCategoryToModles(MenuCategorys []db.MenuCategory) []MenuCategory {
	result := make([]MenuCategory, 0, len(MenuCategorys))
	for _, MenuCategory := range MenuCategorys {
		result = append(result, ConvertMenuCategoryModel(MenuCategory))
	}

	return result
}

type ListMenuCategoriesReq struct {
	TenantID int64 `json:"tenant_id"`
	MenuID   int64 `json:"menu_id"`
	Page     int   `json:"page"`
	Limit    int   `json:"limit"`
	Offset   int   `json:"offset"`
}

type UpdateMenuCategoryReq struct {
	ID          int64
	MenuID      int64
	TenantID    int64
	Name        string
	Description string
	SortOrder   int64
	IsActive    bool
}

type MenuItem struct {
	ID          int             `json:"id" db:"id"`
	CategoryID  int64           `json:"category_id" db:"category_id"`
	MenuID      int64           `json:"menu_id" db:"menu_id"`
	TenantID    int             `json:"tenant_id" db:"tenant_id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	Price       decimal.Decimal `json:"price" db:"price"`
	IsAvailable bool            `json:"is_available" db:"is_available"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateMenuItemReq struct {
	TenantID    int64           `json:"tenant_id"`
	MenuID      int64           `json:"menu_id"`
	CategoryID  int64           `json:"category_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	IsAvailable bool            `json:"is_available"`
}

func ConvertMenuItemModel(menuItem db.MenuItem) MenuItem {
	return MenuItem{
		ID:          int(menuItem.ID),
		TenantID:    int(menuItem.TenantID),
		MenuID:      menuItem.MenuID,
		CategoryID:  menuItem.CategoryID,
		Name:        menuItem.Name,
		Description: menuItem.Description.String,
		Price:       menuItem.Price,
		IsAvailable: menuItem.IsAvailable,
		CreatedAt:   menuItem.CreatedAt.Time,
		UpdatedAt:   menuItem.UpdatedAt.Time,
	}
}

func ConvertMenuItemsToModels(menuItems []db.MenuItem) []MenuItem {
	result := make([]MenuItem, 0, len(menuItems))
	for _, menuItem := range menuItems {
		result = append(result, ConvertMenuItemModel(menuItem))
	}

	return result
}

type ListMenuItemsReq struct {
	TenantID   int64 `json:"tenant_id"`
	MenuID     int64 `json:"menu_id"`
	CategoryID int64 `json:"category_id"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
}

type UpdateMenuItemReq struct {
	ID          int64           `json:"id"`
	CategoryID  int64           `json:"category_id"`
	MenuID      int64           `json:"menu_id"`
	TenantID    int64           `json:"tenant_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	IsAvailable bool            `json:"is_available"`
}
