package action

import (
	apbv1 "Advanced_Shop/api/action/v1"
	"Advanced_Shop/gnova/registry"
	"Advanced_Shop/gnova/server/rpcserver"
	"Advanced_Shop/gnova/server/rpcserver/clientinterceptors"
	"context"
)

const optionserviceName = "discovery:///xshop-option-srv"

type Action struct {
	gc apbv1.AddressClient
	uc apbv1.UserFavClient
	mc apbv1.MessageClient
}

func NewActionServiceClient(r registry.Discovery) (apbv1.AddressClient, apbv1.UserFavClient, apbv1.MessageClient) {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(optionserviceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c1 := apbv1.NewAddressClient(conn)
	c2 := apbv1.NewUserFavClient(conn)
	c3 := apbv1.NewMessageClient(conn)
	return c1, c2, c3
}
