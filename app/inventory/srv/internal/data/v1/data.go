package v1

import "gorm.io/gorm"

type DataFactory interface {
	Inventorys() InventoryStore

	Begin() *gorm.DB
}
