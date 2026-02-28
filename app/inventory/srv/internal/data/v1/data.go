package v1

import (
	"context"
	"gorm.io/gorm"
)

type DataFactory interface {
	Inventorys() InventoryStore
	Listen(ctx context.Context)
	Begin() *gorm.DB
}
