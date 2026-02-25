package srv

import (
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/app/user/srv/config"
	gapp "Advanced_Shop/gnova/app"
	"Advanced_Shop/gnova/server/rpcserver"
	"Advanced_Shop/pkg/app"
	"Advanced_Shop/pkg/log"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"

	"Advanced_Shop/gnova/registry"
	"Advanced_Shop/gnova/registry/consul"
)

var ProviderSet = wire.NewSet(NewUserApp, NewRegistrar, NewUserRPCServer, NewNacosDataSource)

// NewApp user 服务端的总启动配置  返回所有服务共享的 app结构体 但只输入user服务的
func NewApp() *app.App {
	//这里这个new 做了很多事情 比如 log的初始化 你的rpc服务的端口，name 注册逻辑的前置参数初始化 等等
	cfg := config.New()
	appl := app.NewApp("user",
		"xshop",
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

// NewUserApp  这里返回的是 gnova 的 app
func NewUserApp(logOpts *log.Options, registry registry.Registrar,
	serverOpts *options.ServerOptions, rpcServer *rpcserver.Server) (*gapp.App, error) {
	//初始化log
	log.Init(logOpts)
	defer log.Flush()

	return gapp.New(
		gapp.WithName(serverOpts.Name),
		gapp.WithRPCServer(rpcServer),
		gapp.WithRegistrar(registry), // 在这注入   注册和注销 函数的实现
	), nil
}

func run(cfg *config.Config) app.RunFunc {
	return func(baseName string) error {
		userApp, err := initApp(cfg.Nacos, cfg.Log, cfg.Server, cfg.Registry, cfg.Telemetry, cfg.MySQLOptions)
		if err != nil {
			return err
		}

		//启动
		if err := userApp.Run(); err != nil {
			log.Errorf("run user app error: %s", err)
			return err
		}
		return nil
	}
}
