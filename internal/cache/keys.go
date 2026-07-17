package cache

import (
	"fmt"
	"time"
)

const (
	TTLMenu  = 10 * time.Minute
	TTLTable = 5 * time.Minute
	TTLBill  = 30 * time.Minute
)

func MenusKey(tenantID int64) string {
	return fmt.Sprintf("menus:tenant:%d", tenantID)
}

func MecnuCategoriesKey(tenantID, menuID int64) string {
	return fmt.Sprintf("menu_catagories:tenant:%d:menu:%d", tenantID, menuID)
}

func MenuItemsKey(tenantID, menuID, categoryID int64) string {
	return fmt.Sprintf("menu_items:tenant:%d:menu:%d:category:%d", tenantID, menuID, categoryID)
}

func TablesKey(tenantID int64, status string) string {
	return fmt.Sprintf("tables:tenant:%d:status:%s", tenantID, status)
}

func BillKey(referenceID string) string {
	return fmt.Sprintf("bill:%s", referenceID)
}

func TenantMenusPattern(tenantID int64) string {
	return fmt.Sprintf("menus:tenant:%d", tenantID)
}

func TenantMenuCategoriesPattern(tenantID int64) string {
	return fmt.Sprintf("menu_categories:tenant:%d", tenantID)
}

func TenantMenuItemsPattern(tenantID int64) string {
	return fmt.Sprintf("menu_items:tenant:%d", tenantID)
}

func TenantTablesPattern(tenantID int64) string {
	return fmt.Sprintf("tables:tenant:%d*", tenantID)
}
