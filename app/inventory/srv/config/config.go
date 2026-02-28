package config

import (
	"Advanced_Shop/app/pkg/options"
	cliflag "Advanced_Shop/pkg/common/cli/flag"
	"Advanced_Shop/pkg/log"
	"encoding/json"
)

type Config struct {
	MySQLOptions *options.MySQLOptions     `json:"mysql"     mapstructure:"mysql"`
	Log          *log.Options              `json:"log"     mapstructure:"log"`
	Server       *options.ServerOptions    `json:"server"     mapstructure:"server"`
	Telemetry    *options.TelemetryOptions `json:"telemetry" mapstructure:"telemetry"`
	Registry     *options.RegistryOptions  `json:"registry" mapstructure:"registry"`
	RedisOptions *options.RedisOptions     `json:"redis" mapstructure:"redis"`
	Mq           *options.RocketMQOptions  `json:"mq" mapstructure:"mq"`
}

func New() *Config {
	//配置默认初始化
	return &Config{
		MySQLOptions: options.NewMySQLOptions(),
		Log:          log.NewOptions(),
		Server:       options.NewServerOptions(),
		Telemetry:    options.NewTelemetryOptions(),
		Registry:     options.NewRegistryOptions(),
		RedisOptions: options.NewRedisOptions(),
		Mq:           options.NewRocketMQOptions(),
	}
}

// Flags returns flags for a specific APIServer by section name.
func (o *Config) Flags() (fss cliflag.NamedFlagSets) {
	o.Server.AddFlags(fss.FlagSet("server"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.Telemetry.AddFlags(fss.FlagSet("telemetry"))
	o.Registry.AddFlags(fss.FlagSet("registry"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.Mq.AddFlags(fss.FlagSet("mq"))
	return fss
}

func (o *Config) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func (o *Config) Validate() []error {
	var errs []error

	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.Server.Validate()...)
	errs = append(errs, o.Telemetry.Validate()...)
	errs = append(errs, o.Registry.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.Mq.Validate()...)
	return errs
}
