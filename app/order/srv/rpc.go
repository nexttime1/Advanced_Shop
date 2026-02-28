package srv

import (
	gpb "Advanced_Shop/api/order/v1"
	"Advanced_Shop/app/order/srv/config"
	"Advanced_Shop/app/order/srv/internal/controller/order/v1"
	db2 "Advanced_Shop/app/order/srv/internal/data/v1/realize"
	v13 "Advanced_Shop/app/order/srv/internal/service/v1"
	"Advanced_Shop/gnova/core/trace"
	"Advanced_Shop/gnova/server/rpcserver"
	"context"
	"fmt"
)

func NewOrderRPCServer(cfg *config.Config, ctx context.Context) (*rpcserver.Server, error) {
	//初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		cfg.Telemetry.Name,
		cfg.Telemetry.Endpoint,
		cfg.Telemetry.Sampler,
		cfg.Telemetry.Batcher,
	})

	dataFactory := db2.NewDataFactory(cfg.MySQLOptions, cfg.Registry, cfg.MQOptions)
	// 监听延时消息
	dataFactory.Listen(ctx)
	orderSrvFactory := v13.NewService(dataFactory, cfg.Dtm, cfg.MQOptions)
	orderServer := order.NewOrderServer(orderSrvFactory)
	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))
	gpb.RegisterOrderServer(grpcServer.Server, orderServer)
	return grpcServer, nil
}
