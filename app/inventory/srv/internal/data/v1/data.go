package v1

import (
	"context"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"gorm.io/gorm"
)

type DataFactory interface {
	Inventorys() InventoryStore
	Listen(ctx context.Context)
	Begin() *gorm.DB
	DB() *gorm.DB
	Pool() redsyncredis.Pool
}
