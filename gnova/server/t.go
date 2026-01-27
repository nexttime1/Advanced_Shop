package server

import (
	v1 "Advanced_Shop/api/user/v1"
	"context"
)

func main() {
	clientConn, err := DialInsecure(context.Background())
	if err != nil {
		panic(err)
	}
	client := v1.NewUserClient(clientConn)
	resp, err := client.GetUserList(ctx, req)
}
