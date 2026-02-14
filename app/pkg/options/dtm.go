package options

import "github.com/spf13/pflag"

type DtmOptions struct {
	GrpcServer string `mapstructure:"grpc" json:"grpc,omitempty"`
	HttpServer string `mapstructure:"http" json:"http,omitempty"`
}

func NewDtmOptions() *DtmOptions {
	return &DtmOptions{
		HttpServer: "http://127.0.0.1:36789/api/dtmsvr",
		GrpcServer: "127.0.0.1:36790",
	}
}

func (o *DtmOptions) Validate() []error {
	errs := []error{}
	return errs
}

func (o *DtmOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.GrpcServer, "dtm.grpc", o.GrpcServer, "")
	fs.StringVar(&o.HttpServer, "dtm.http", o.HttpServer, "")
}
