package order

import (
	opbv1 "Advanced_Shop/api/order/v1"
	"Advanced_Shop/gnova/server/rpcserver"
	"Advanced_Shop/gnova/server/rpcserver/clientinterceptors"
	"context"

	"Advanced_Shop/gnova/registry"
)

const goodsserviceName = "discovery:///xshop-order-srv"

type Order struct {
	gc opbv1.OrderClient
}

func NewOrder(g opbv1.OrderClient) *Order {
	return &Order{g}
}

func NewOrderServiceClient(r registry.Discovery) opbv1.OrderClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(goodsserviceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := opbv1.NewOrderClient(conn)
	return c
}
