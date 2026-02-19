package inventory

import (
	ipbv1 "Advanced_Shop/api/inventory/v1"
	"Advanced_Shop/gnova/server/rpcserver"
	"Advanced_Shop/gnova/server/rpcserver/clientinterceptors"
	"context"

	"Advanced_Shop/gnova/registry"
)

const inventoryserviceName = "discovery:///xshop-inventory-srv"

type inventory struct {
	gc ipbv1.InventoryClient
}

func NewInventoryServiceClient(r registry.Discovery) ipbv1.InventoryClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(inventoryserviceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := ipbv1.NewInventoryClient(conn)
	return c
}
