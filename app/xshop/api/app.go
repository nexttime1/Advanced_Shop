package admin

import (
	"context"

	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/app/xshop/api/config"
	gapp "Advanced_Shop/gnova/app"
	"Advanced_Shop/pkg/app"
	"Advanced_Shop/pkg/log"
	"github.com/hashicorp/consul/api"

	"Advanced_Shop/gnova/registry"
	"Advanced_Shop/gnova/registry/consul"
	"Advanced_Shop/pkg/storage"
)

func NewApp(basename string) *app.App {
	cfg := config.New()
	appl := app.NewApp("api",
		"xshop",
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
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

func NewAPIApp(cfg *config.Config) (*gapp.App, error) {
	// 初始化log
	log.Init(cfg.Log)
	defer log.Flush()

	// 服务注册
	register := NewRegistrar(cfg.Registry)

	//  连接redis
	redisConfig := &storage.Config{
		Host:                  cfg.Redis.Host,
		Port:                  cfg.Redis.Port,
		Addrs:                 cfg.Redis.Addrs,
		MasterName:            cfg.Redis.MasterName,
		Username:              cfg.Redis.Username,
		Password:              cfg.Redis.Password,
		Database:              cfg.Redis.Database,
		MaxIdle:               cfg.Redis.MaxIdle,
		MaxActive:             cfg.Redis.MaxActive,
		Timeout:               cfg.Redis.Timeout,
		EnableCluster:         cfg.Redis.EnableCluster,
		UseSSL:                cfg.Redis.UseSSL,
		SSLInsecureSkipVerify: cfg.Redis.SSLInsecureSkipVerify,
		EnableTracing:         cfg.Redis.EnableTracing,
	}

	go storage.ConnectToRedis(context.Background(), redisConfig)

	//生成http服务
	restServer, err := NewAPIHTTPServer(cfg)
	if err != nil {
		return nil, err
	}

	return gapp.New(
		gapp.WithName(cfg.Server.Name),
		gapp.WithRestServer(restServer),
		gapp.WithRegistrar(register),
	), nil
}

func run(cfg *config.Config) app.RunFunc {
	return func(baseName string) error {
		apiApp, err := NewAPIApp(cfg)
		if err != nil {
			return err
		}

		//启动
		if err := apiApp.Run(); err != nil {
			log.Errorf("run api app error: %s", err)
			return err
		}
		return nil
	}
}
