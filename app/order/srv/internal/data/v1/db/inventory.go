package db

import (
	"context"

	proto "Advanced_Shop/api/inventory/v1"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/server/rpcserver"
	"Advanced_Shop/gnova/server/rpcserver/clientinterceptors"

	"Advanced_Shop/gnova/registry"
)

const ginvserviceName = "discovery:///xshop-inventory-srv"

func GetInventoryClient(opts *options.RegistryOptions) proto.InventoryClient {
	discovery := NewDiscovery(opts)
	invClient := NewInventoryServiceClient(discovery)
	return invClient
}

func NewInventoryServiceClient(r registry.Discovery) proto.InventoryClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(ginvserviceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := proto.NewInventoryClient(conn)
	return c
}
