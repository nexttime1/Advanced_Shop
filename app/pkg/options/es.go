package options

import "github.com/spf13/pflag"

type EsOptions struct {
	Host string `json:"host" mapstructure:"host"`
	Port string `json:"port" mapstructure:"port"`
}

func NewEsOptions() *EsOptions {
	return &EsOptions{
		Host: "127.0.0.1",
		Port: "9200",
	}
}

func (e *EsOptions) Validate() []error {
	errs := []error{}
	return errs
}

func (e *EsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&e.Host, "es.host", e.Host, ""+
		"es service host address. If left blank, the following related es options will be ignored.")

	fs.StringVar(&e.Port, "es.port", e.Port, ""+
		"es service port If left blank, the following related es options will be ignored..")
}
