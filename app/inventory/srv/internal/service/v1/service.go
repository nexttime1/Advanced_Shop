package v1

import (
	v1 "Advanced_Shop/app/inventory/srv/internal/data/v1"
	"Advanced_Shop/app/pkg/options"
	"fmt"

	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type ServiceFactory interface {
	Inventories() InventorySrv
}

type service struct {
	data v1.DataFactory

	redisOptions *options.RedisOptions
	pool         redsyncredis.Pool
}

func (s *service) Inventories() InventorySrv {
	return newInventoryService(s)
}

func NewService(store v1.DataFactory, redisOptions *options.RedisOptions) ServiceFactory {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", redisOptions.Host, redisOptions.Port),
	})
	pool := goredis.NewPool(client)

	return &service{data: store, redisOptions: redisOptions, pool: pool}
}

var _ ServiceFactory = &service{}
