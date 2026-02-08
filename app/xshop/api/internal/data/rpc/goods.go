package rpc

import (
	gpbv1 "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/gnova/server/rpcserver"
	"Advanced_Shop/gnova/server/rpcserver/clientinterceptors"
	"context"

	"Advanced_Shop/gnova/registry"
)

const goodsserviceName = "discovery:///xshop-goods-srv"

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
