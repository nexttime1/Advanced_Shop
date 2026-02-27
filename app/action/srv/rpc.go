package srv

import (
	apb "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/action/srv/config"
	v12 "Advanced_Shop/app/action/srv/internal/controller/v1"
	db2 "Advanced_Shop/app/action/srv/internal/data/v1/db"
	v1 "Advanced_Shop/app/action/srv/internal/service/v1"

	"Advanced_Shop/gnova/core/trace"
	"Advanced_Shop/gnova/server/rpcserver"
	"fmt"

	"Advanced_Shop/pkg/log"
)

func NewActionRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		cfg.Telemetry.Name,
		cfg.Telemetry.Endpoint,
		cfg.Telemetry.Sampler,
		cfg.Telemetry.Batcher,
	})

	dataFactory, err := db2.GetDBFactoryOr(cfg.MySQLOptions, cfg.Registry)
	if err != nil {
		log.Fatal(err.Error())
	}

	srvFactory := v1.NewService(dataFactory)
	actionServer := v12.NewActionServer(srvFactory)
	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))

	apb.RegisterUserFavServer(grpcServer.Server, actionServer)
	apb.RegisterMessageServer(grpcServer.Server, actionServer)
	apb.RegisterAddressServer(grpcServer.Server, actionServer)

	return grpcServer, nil
}
