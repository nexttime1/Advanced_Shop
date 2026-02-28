package srv

import (
	gpb "Advanced_Shop/api/inventory/v1"
	"Advanced_Shop/app/inventory/srv/config"
	v12 "Advanced_Shop/app/inventory/srv/internal/controller/v1"
	db2 "Advanced_Shop/app/inventory/srv/internal/data/v1/db"
	v13 "Advanced_Shop/app/inventory/srv/internal/service/v1"
	"Advanced_Shop/gnova/core/trace"
	"Advanced_Shop/gnova/server/rpcserver"
	"context"
	"fmt"

	"Advanced_Shop/pkg/log"
)

func NewInventoryRPCServer(cfg *config.Config, ctx context.Context) (*rpcserver.Server, error) {
	//初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		cfg.Telemetry.Name,
		cfg.Telemetry.Endpoint,
		cfg.Telemetry.Sampler,
		cfg.Telemetry.Batcher,
	})

	//有点繁琐，wire， ioc-golang
	dataFactory, err := db2.GetDBFactoryOr(cfg.MySQLOptions, cfg.Mq)
	dataFactory.Listen(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	invService := v13.NewService(dataFactory, cfg.RedisOptions)

	invServer := v12.NewInventoryServer(invService)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))

	gpb.RegisterInventoryServer(grpcServer.Server, invServer)

	return grpcServer, nil
}
