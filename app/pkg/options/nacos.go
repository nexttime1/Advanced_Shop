package options

import "github.com/spf13/pflag"

type NacosOptions struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" and:"group"`
}

func NewNacosOptions() *NacosOptions {
	return &NacosOptions{
		Host:      "127.0.0.1",
		Port:      8848,
		Namespace: "public",
		User:      "nacos",
		Password:  "nacos",
		DataId:    "flow",
		Group:     "sentinel-go",
	}
}

func (n *NacosOptions) Validate() []error {
	errs := []error{}

	return errs
}

func (n *NacosOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&n.Host, "nacos.host", n.Host, "nacos host")
	fs.Uint64Var(&n.Port, "nacos.port", n.Port, "nacos port")
	fs.StringVar(&n.Namespace, "nacos.namespace", n.Namespace, "nacos namespace")
	fs.StringVar(&n.User, "nacos.user", n.User, "nacos user")
	fs.StringVar(&n.Password, "nacos.password", n.Password, "nacos password")
	fs.StringVar(&n.DataId, "nacos.dataid", n.DataId, "nacos dataid")
	fs.StringVar(&n.Group, "nacos.group", n.Group, "nacos group")
}
