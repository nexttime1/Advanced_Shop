package rpc

import (
	gpb "Advanced_Shop/api/goods/v1"
	upb "Advanced_Shop/api/user/v1"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/app/xshop/api/internal/data"
	"Advanced_Shop/app/xshop/api/internal/data/rpc/good"
	code2 "Advanced_Shop/gnova/code"
	"Advanced_Shop/gnova/registry"
	"Advanced_Shop/gnova/registry/consul"
	errors2 "Advanced_Shop/pkg/errors"
	"fmt"
	cosulAPI "github.com/hashicorp/consul/api"
	"sync"
)

// grpcData 实现工厂接口
type grpcData struct {
	gc gpb.GoodsClient
	uc upb.UserClient
}

func (g grpcData) Goods() gpb.GoodsClient {
	return g.gc
}

func (g grpcData) Users() data.UserData {
	return NewUsers(g.uc)
}

func NewDiscovery(opts *options.RegistryOptions) registry.Discovery {
	c := cosulAPI.DefaultConfig()
	c.Address = opts.Address
	c.Scheme = opts.Scheme
	cli, err := cosulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}

var (
	dbFactory data.DataFactory
	once      sync.Once
)

// GetDataFactoryOr rpc的连接， 基于服务发现
func GetDataFactoryOr(options *options.RegistryOptions) (data.DataFactory, error) {
	if options == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get grpc store fatory")
	}

	//这里负责依赖的所有的rpc连接
	once.Do(func() {
		discovery := NewDiscovery(options)
		userClient := NewUserServiceClient(discovery)
		goodsClient := good.NewGoodsServiceClient(discovery)
		dbFactory = &grpcData{
			gc: goodsClient,
			uc: userClient,
		}
	})

	if dbFactory == nil {
		return nil, errors2.WithCode(code2.ErrConnectGRPC, "failed to get grpc store factory")
	}
	return dbFactory, nil
}
