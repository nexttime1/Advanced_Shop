package v1

import (
	v1 "Advanced_Shop/app/inventory/srv/internal/data/v1"
	"Advanced_Shop/app/pkg/options"
)

type ServiceFactory interface {
	Inventories() InventorySrv
}

type service struct {
	data v1.DataFactory

	redisOptions *options.RedisOptions
}

func (s *service) Inventories() InventorySrv {
	return newInventoryService(s)
}

func NewService(store v1.DataFactory, redisOptions *options.RedisOptions) ServiceFactory {

	return &service{
		data:         store,
		redisOptions: redisOptions,
	}
}

var _ ServiceFactory = &service{}
