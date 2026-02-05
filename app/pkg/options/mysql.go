package options

import (
	"github.com/spf13/pflag"
	"time"
)

type MySQLOptions struct {
	Host                  string        `mapstructure:"host" json:"host,omitempty"`
	Port                  string        `mapstructure:"port" json:"port,omitempty"`
	Username              string        `mapstructure:"username" json:"username,omitempty"`
	Password              string        `mapstructure:"password" json:"password,omitempty"`
	Database              string        `mapstructure:"database" json:"database"`
	MaxIdleConnections    int           `mapstructure:"max-idle-connections" json:"max-idle-connections,omitempty"`
	MaxOpenConnections    int           `mapstructure:"max-open-connections" json:"max-open-connections,omitempty"`
	MaxConnectionLifetime time.Duration `mapstructure:"ax-connection-life-time" json:"max-connection-life-time,omitempty"`
	LogLevel              int           `mapstructure:"log-level" json:"log-level"`
}

// NewMySQLOptions create a `zero` value instance.
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1",
		Port:                  "3306",
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

// Validate verifies flags passed to MySQLOptions.
func (o *MySQLOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (mo *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&mo.Host, "mysql.host", mo.Host, ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

	fs.StringVar(&mo.Port, "mysql.port", mo.Port, ""+
		"MySQL service port")

	fs.StringVar(&mo.Username, "mysql.username", mo.Username, ""+
		"Username for access to mysql service.")

	fs.StringVar(&mo.Password, "mysql.password", mo.Password, ""+
		"Password for access to mysql, should be used pair with password.")

	fs.StringVar(&mo.Database, "mysql.database", mo.Database, ""+
		"Database name for the server to use.")

	fs.IntVar(&mo.MaxIdleConnections, "mysql.max-idle-connections", mo.MaxOpenConnections, ""+
		"Maximum idle connections allowed to connect to mysql.")

	fs.IntVar(&mo.MaxOpenConnections, "mysql.max-open-connections", mo.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to mysql.")

	fs.DurationVar(&mo.MaxConnectionLifetime, "mysql.max-connection-life-time", mo.MaxConnectionLifetime, ""+
		"Maximum connection life time allowed to connecto to mysql.")

	fs.IntVar(&mo.LogLevel, "mysql.log-mode", mo.LogLevel, ""+
		"Specify gorm log level.")
}
