package admin

import (
	"Advanced_Shop/app/xshop/api/config"
	"Advanced_Shop/gnova/server/restserver"
)

func NewAPIHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	aRestServer := restserver.NewServer(restserver.WithPort(cfg.Server.HttpPort),
		restserver.WithMiddlewares(cfg.Server.Middlewares),
		restserver.WithMetrics(true),
	)

	//配置好路由
	initRouter(aRestServer, cfg)

	return aRestServer, nil
}
