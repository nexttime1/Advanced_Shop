package srv

import (
	gpb "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/goods/srv/config"
	v12 "Advanced_Shop/app/goods/srv/internal/controller/v1"
	db2 "Advanced_Shop/app/goods/srv/internal/data/v1/db"
	"Advanced_Shop/app/goods/srv/internal/data_search/v1/es"
	v1 "Advanced_Shop/app/goods/srv/internal/service/v1"
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

	dataFactory, err := db2.GetDBFactoryOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	//构建，繁琐 - 工厂模式
	searchFactory, err := es.GetSearchFactoryOr(cfg.EsOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	srvFactory := v1.NewService(dataFactory, searchFactory)
	goodsServer := v12.NewGoodsServer(srvFactory)
	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))

	gpb.RegisterGoodsServer(grpcServer.Server, goodsServer)

	return grpcServer, nil
}
