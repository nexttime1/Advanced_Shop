package options

import "github.com/spf13/pflag"

type SmsOptions struct {
	APIKey    string `mapstructure:"key" json:"key"`
	APISecret string `mapstructure:"secret" json:"secret"`
}

func NewSmsOptions() *SmsOptions {
	return &SmsOptions{
		APIKey:    "",
		APISecret: "",
	}
}

func (s *SmsOptions) Validate() []error {
	errs := []error{}
	return errs
}

func (o *SmsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.APIKey, "sms.apikey", o.APIKey, ""+
		"sms apikey")

	fs.StringVar(&o.APISecret, "sms.secret", o.APISecret, ""+
		"sms api secret")
}
