//go:build wireinject
// +build wireinject

package srv

import (
	"github.com/google/wire"
	"mxshop/app/pkg/options"
	"mxshop/app/user/srv/controller/user"
	"mxshop/app/user/srv/data/v1/db"
	v1 "mxshop/app/user/srv/service/v1"
	gapp "mxshop/gmicro/app"
	"mxshop/pkg/log"
)

func initApp(*options.NacosOptions, *log.Options, *options.ServerOptions, *options.RegistryOptions, *options.TelemetryOptions, *options.MySQLOptions) (*gapp.App, error) {
	wire.Build(ProviderSet, v1.ProviderSet, db.ProviderSet, user.ProviderSet)
	return &gapp.App{}, nil
}
