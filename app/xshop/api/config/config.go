package config

import (
	"Advanced_Shop/app/pkg/options"
	cliflag "Advanced_Shop/pkg/common/cli/flag"
	"Advanced_Shop/pkg/log"
)

type Config struct {
	Log *log.Options `json:"log" mapstructure:"log"`

	Server   *options.ServerOptions   `json:"server" mapstructure:"server"`
	Registry *options.RegistryOptions `json:"registry" mapstructure:"registry"`
	Jwt      *options.JwtOptions      `json:"jwt" mapstructure:"jwt"`
	Sms      *options.SmsOptions      `json:"sms" mapstructure:"sms"`
	Redis    *options.RedisOptions    `json:"redis" mapstructure:"redis"`
}

func (c *Config) Validate() []error {
	var errors []error
	errors = append(errors, c.Log.Validate()...)
	errors = append(errors, c.Server.Validate()...)
	errors = append(errors, c.Registry.Validate()...)
	errors = append(errors, c.Jwt.Validate()...)
	errors = append(errors, c.Sms.Validate()...)
	errors = append(errors, c.Redis.Validate()...)
	return errors
}

func (c *Config) Flags() (fss cliflag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	c.Jwt.AddFlags(fss.FlagSet("jwt"))
	c.Sms.AddFlags(fss.FlagSet("sms"))
	c.Redis.AddFlags(fss.FlagSet("redis"))
	return fss
}

func New() *Config {
	//配置默认初始化
	return &Config{
		Log:      log.NewOptions(),
		Server:   options.NewServerOptions(),
		Registry: options.NewRegistryOptions(),
		Jwt:      options.NewJwtOptions(),
		Sms:      options.NewSmsOptions(),
		Redis:    options.NewRedisOptions(),
	}
}
