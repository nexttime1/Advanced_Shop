package srv

import (
	"Advanced_Shop/app/action/srv/config"
	"Advanced_Shop/app/pkg/options"
	gapp "Advanced_Shop/gnova/app"
	"Advanced_Shop/pkg/app"
	"Advanced_Shop/pkg/log"
	"context"
	"github.com/hashicorp/consul/api"

	"Advanced_Shop/gnova/registry"
	"Advanced_Shop/gnova/registry/consul"
)

func NewApp(basename string) *app.App {
	cfg := config.New()
	appl := app.NewApp("option",
		basename,
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		//app.WithNoConfig(), //设置不读取配置文件
	)
	return appl
}

func NewRegistrar(registry *options.RegistryOptions) registry.Registrar {
	c := api.DefaultConfig()
	c.Address = registry.Address
	c.Scheme = registry.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(true))
	return r
}

func NewGoodsApp(cfg *config.Config) (*gapp.App, error) {
	//初始化log
	log.Init(cfg.Log)
	defer log.Flush()

	//服务注册
	register := NewRegistrar(cfg.Registry)

	//生成rpc服务
	rpcServer, err := NewActionRPCServer(cfg)
	if err != nil {
		return nil, err
	}

	return gapp.New(
		gapp.WithName(cfg.Server.Name),
		gapp.WithRPCServer(rpcServer),
		gapp.WithRegistrar(register),
	), nil
}

func run(cfg *config.Config) app.RunFunc {
	return func(baseName string, ctx context.Context) error {
		goodsApp, err := NewGoodsApp(cfg)
		if err != nil {
			return err
		}

		//启动
		if err := goodsApp.Run(ctx); err != nil {
			log.Errorf("run user app error: %s", err)
			return err
		}
		return nil
	}
}
