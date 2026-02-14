package main

import (
	v1 "Advanced_Shop/api/order/v1"
	"Advanced_Shop/gnova/registry/consul"
	rpc "Advanced_Shop/gnova/server/rpcserver"
	_ "Advanced_Shop/gnova/server/rpcserver/resolver/direct"
	"Advanced_Shop/gnova/server/rpcserver/selector"
	"Advanced_Shop/gnova/server/rpcserver/selector/random"
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"time"
)

func generateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

func main() {
	//设置全局的负载均衡策略
	selector.SetGlobalSelector(random.NewBuilder())
	rpc.InitBuilder()

	conf := api.DefaultConfig()
	conf.Address = "127.0.0.1:8500"
	conf.Scheme = "http"
	cli, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))

	conn, err := rpc.DialInsecure(context.Background(),
		rpc.WithBalancerName("selector"),
		rpc.WithDiscovery(r),
		rpc.WithClientTimeout(time.Second*5000),
		rpc.WithEndpoint("discovery:///Advanced_Shop-order-srv"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := v1.NewOrderClient(conn)

	_, err = uc.SubmitOrder(context.Background(), &v1.OrderRequest{
		UserId:  1,
		Address: "慕课网",
		OrderSn: generateOrderSn(1),
		Name:    "bobby",
		Post:    "尽快发货",
		Mobile:  "18787878787",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("订单创建成功")
}
