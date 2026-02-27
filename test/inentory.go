package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
	// 替换为你inventory-srv的proto生成的包路径
	inventoryPb "Advanced_Shop/api/inventory/v1"
)

func main() {
	// 连接inventory-srv的8022端口
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 加长超时
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		"10.121.30.224:8022",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		fmt.Printf("gRPC连接失败：%v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("gRPC连接成功，开始调用Sell方法")

	// 创建客户端并调用Sell方法（参数按你的proto定义调整）
	client := inventoryPb.NewInventoryClient(conn)
	req := &inventoryPb.SellInfo{
		OrderSn: "213123123123",
		GoodsInfo: []*inventoryPb.GoodsInvInfo{
			{
				GoodsId: 1,
				Num:     2,
			},
		},
	}
	resp, err := client.Sell(ctx, req)
	if err != nil {
		fmt.Printf("调用Sell方法失败：%v\n", err)
		return
	}
	fmt.Printf("Sell方法调用成功，响应：%+v\n", resp)
}
