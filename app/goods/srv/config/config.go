package config

import (
	"Advanced_Shop/app/pkg/options"
	cliflag "Advanced_Shop/pkg/common/cli/flag"
	"Advanced_Shop/pkg/log"
)

type Config struct {
	Log       *log.Options       `json:"log" mapstructure:"log"`
	EsOptions *options.EsOptions `json:"es" mapstructure:"es"`

	Server       *options.ServerOptions    `json:"server" mapstructure:"server"`
	Registry     *options.RegistryOptions  `json:"registry" mapstructure:"registry"`
	Telemetry    *options.TelemetryOptions `json:"telemetry" mapstructure:"telemetry"`
	MySQLOptions *options.MySQLOptions     `json:"mysql" mapstructure:"mysql"`
	CanalOpts    *options.CanalOptions     `json:"canal" mapstructure:"canal"`
	MqOpts       *options.RocketMQOptions  `json:"mq" mapstructure:"mq"`
}

func (c *Config) Validate() []error {
	var errors []error
	errors = append(errors, c.Log.Validate()...)
	errors = append(errors, c.Server.Validate()...)
	errors = append(errors, c.Registry.Validate()...)
	errors = append(errors, c.Telemetry.Validate()...)
	errors = append(errors, c.MySQLOptions.Validate()...)
	errors = append(errors, c.EsOptions.Validate()...)
	errors = append(errors, c.CanalOpts.Validate()...)
	errors = append(errors, c.MqOpts.Validate()...)
	return errors
}

func (c *Config) Flags() (fss cliflag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	c.Telemetry.AddFlags(fss.FlagSet("telemetry"))
	c.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	c.EsOptions.AddFlags(fss.FlagSet("es"))
	c.CanalOpts.AddFlags(fss.FlagSet("canal"))
	c.MqOpts.AddFlags(fss.FlagSet("rabbitmq"))
	return fss
}

func New() *Config {
	//配置默认初始化
	return &Config{
		Log:          log.NewOptions(),
		Server:       options.NewServerOptions(),
		Registry:     options.NewRegistryOptions(),
		Telemetry:    options.NewTelemetryOptions(),
		MySQLOptions: options.NewMySQLOptions(),
		EsOptions:    options.NewEsOptions(),
		CanalOpts:    options.NewCanalOptions(),
		MqOpts:       options.NewRocketMQOptions(),
	}
}
