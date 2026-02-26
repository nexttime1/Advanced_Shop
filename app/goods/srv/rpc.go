package srv

import (
	gpb "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/goods/srv/config"
	v12 "Advanced_Shop/app/goods/srv/internal/controller/v1"
	data "Advanced_Shop/app/goods/srv/internal/data/v1/realize"
	"Advanced_Shop/app/goods/srv/internal/data_search/v1/es"
	v1 "Advanced_Shop/app/goods/srv/internal/service/v1"
	"Advanced_Shop/gnova/core/trace"
	"Advanced_Shop/gnova/server/rpcserver"
	"context"
	"fmt"
	"time"

	"Advanced_Shop/pkg/log"
)

func NewGoodsRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	//初始化open-telemetry的exporter
	trace.InitAgent(trace.Options{
		cfg.Telemetry.Name,
		cfg.Telemetry.Endpoint,
		cfg.Telemetry.Sampler,
		cfg.Telemetry.Batcher,
	})

	//有点繁琐，wire， ioc-golang
	dataFactory := data.NewDataStore(cfg.MySQLOptions, cfg.MqOpts, cfg.CanalOpts)
	//构建，繁琐 - 工厂模式
	searchFactory, err := es.GetSearchFactoryOr(cfg.EsOptions, cfg.MqOpts, cfg.CanalOpts)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	// Canal监听器
	/*
		客户端调用Create → 参数校验 → 校验品牌/分类 → 开启MySQL事务 → 写入商品数据 → 提交事务 →
		Canal监听binlog → 解析商品表变更 → 发送RocketMQ消息 → ES消费者消费消息并写入ES
	*/
	dataFactory.StartCanalListener(context.Background())
	time.Sleep(2 * time.Second)
	err = searchFactory.Listen(context.Background())
	if err != nil {
		return nil, err
	}
	srvFactory := v1.NewService(dataFactory, searchFactory)
	goodsServer := v12.NewGoodsServer(srvFactory)
	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	grpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))

	gpb.RegisterGoodsServer(grpcServer.Server, goodsServer)

	return grpcServer, nil
}
