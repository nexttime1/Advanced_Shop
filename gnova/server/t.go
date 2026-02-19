package server

import (
	v1 "Advanced_Shop/api/user/v1"
	"Advanced_Shop/gnova/server/rpcserver"
	"context"
)

func main() {
	clientConn, err := rpcserver.DialInsecure(context.Background())
	if err != nil {
		panic(err)
	}
	client := v1.NewUserClient(clientConn)
	resp, err := client.GetUserList(ctx, req)
}
