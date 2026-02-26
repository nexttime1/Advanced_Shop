package db

import (
	gpbv1 "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/registry/consul"
	"Advanced_Shop/gnova/server/rpcserver"
	"Advanced_Shop/gnova/server/rpcserver/clientinterceptors"
	"context"
	cosulAPI "github.com/hashicorp/consul/api"

	"Advanced_Shop/gnova/registry"
)

const goodsserviceName = "discovery:///xshop-goods-srv"

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

func GetGoodsClient(opts *options.RegistryOptions) gpbv1.GoodsClient {
	discovery := NewDiscovery(opts)
	goodsClient := NewGoodsServiceClient(discovery)
	return goodsClient
}

func NewGoodsServiceClient(r registry.Discovery) gpbv1.GoodsClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(goodsserviceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := gpbv1.NewGoodsClient(conn)
	return c
}
