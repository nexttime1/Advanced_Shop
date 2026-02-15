package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// AliyunOptions 阿里云相关配置，核心包含支付宝支付所需的所有配置项
type AliyunOptions struct {
	// AlipayAppId 支付宝应用ID
	AlipayAppId string `mapstructure:"alipay-app-id" json:"alipayAppId,omitempty"`
	// AlipayPrivateKey 应用私钥
	AlipayPrivateKey string `mapstructure:"alipay-private-key" json:"alipayPrivateKey,omitempty"`
	// AlipayPublicKey 支付宝公钥（从支付宝开放平台获取的公钥）
	AlipayPublicKey string `mapstructure:"alipay-public-key" json:"alipayPublicKey,omitempty"`
	// AlipaySubject 前缀
	AlipaySubject string `mapstructure:"alipay-subject" json:"alipaySubject,omitempty"`
	// AlipayNotifyUrl 支付宝异步通知回调地址
	AlipayNotifyUrl string `mapstructure:"alipay-notify-url" json:"alipayNotifyUrl,omitempty"`
	// AlipayReturnUrl 支付宝同步跳转地址
	AlipayReturnUrl string `mapstructure:"alipay-return-url" json:"alipayReturnUrl,omitempty"`
	// AlipayProductCode 支付宝产品码，固定值FAST_INSTANT_TRADE_PAY（
	AlipayProductCode string `mapstructure:"alipay-product-code" json:"alipayProductCode,omitempty"`
	// AlipayTimeoutExpress 支付链接失效时间，默认30分钟
	AlipayTimeoutExpress string `mapstructure:"alipay-timeout-express" json:"alipayTimeoutExpress,omitempty"`
}

// NewAliyunOptions 创建AliyunOptions实例并设置默认值
func NewAliyunOptions() *AliyunOptions {
	return &AliyunOptions{
		AlipayAppId:          "9021000158665836",
		AlipayPrivateKey:     "MIIEowIBAAKCAQEAmvDGaakWIw8oeWBgWxtOi4ieWetLQoBsll1J8ei6CH7FfHWni1L05UW+/ZOPy83MysrGAupo9F3cWK1NHYz3Kfafza7Vi0brPbR1qrSCLL+O61TPUFzc+fXbNW+/xRpcLMfQvl+eSwVsUBovW7WQU6zwZhG5WqiaxrVvEQRg59yYo1rh2ZytgtXv+LdWSSTLVIETC2i6vstlXLW6AGsj15t7ZPtmieZt4VZx1/LYufSTa/7o03ObEa/wMqE/sf6Md3SujXn0YP2KFOlRYnF8gvPYblhWDBSLE1yD6fV+KyfOPh/P/ggBdYeFbfpJW+pThzXgY7mKnDqiMFGPxWoRGwIDAQABAoIBAFYblclG5Tyawf2iqCo55M77IDYM5AiTYsW2FtBQbIMoIQzoPjLZ6aw5tMksZu/28eeKBb29FJMqTrkhpwfTPdGedHVUwuzifv4N+o7iPq4rz3vN6GFbGpv4HNl3v3YFDlD8w2/pqAk9fFKQGt226/z00a2IECDoLwxb7NviORDhzzlSiSbzBd7sUy6JoaK2fVInG5P9LykDejbwolhE1Kqn+nrWIdFLtKybYQpGNRH0wx087sBH+6uIMfMjNU/1GCGpdK/paR96xe6tJQ4J3iz/Nu0but6lqXiXfMFcJM56zWNr6q3Awr+7oDWuHFSJ6EkvYLpU+FSqVK578Nx6btECgYEA+NV+0eKPbIfLRCIzbVZNvWBkYO9Iwlr9gzEfps/fRyjPDuetHEWpqgcZT0Gz6p21E2vEr8SGMcAS5iI03wYQrHKGwtEBhQ91MnmFMwnAvIbifRt9BnebGML8Z34rrIdJNrg0TdJGg0EaHhCTZ6+7YEXnCZzd+1vbv3gNHY2eOh8CgYEAn2cPK//LVS55+jAjyLLLIjEXGk9YRrsGsewVM468qbM42tlD6T8C56P6Cpt7AI10KmtREPjWhrdgTP8N5PjVp0rMOQOvSzLcKCOO07FxTqXCaSPNRIOdpMOq9RiwQuWBpSEQDp8ArXn7HgEDnzRXkga0cl3WL1qq3uP0GTQMQYUCgYEAoNSlJpVwLC7M85nDcZ0BnDCMUJb4iR50kvISSig7YWwAANs/aXGhStNRyYdm+XK7kfTq6Mx2C/vgezyKvcfWyQ8xCQQ8HjuyfVBMBoP8Ph5Uj5ZPxflSlruYlm/XXKkIakS/Ebmid72BWwNNswvDaWNlBDKOy6NAsk2u9HYPWfMCgYBXIDCFvxl3ZKDdI+TbNQacmLJk+gtpFZ6yLzTjalgqdUBVNj3NRlijHdh0ZclUYvykluXHXgt7tM1ZKGuCxJObDeIUI7RzaMg21ECj6q/g6e8aIqx2j23h+eT+dFEbL3CuPiUVqMjpCOw92RYOtcBLm4iTnkCMv4T3sSbhg7ZTNQKBgAjrfm12WT1kzE7P9SVaoiNh3nqqyG86r6CW4tGfUdwl0PwoH/LYS8rMNBXpsGY61JqvZtsY7DlYyly+W5bezaTBeGTXzdNaLECs87i2KtNKbZgClK6S3r/6s+nRpal0OJVIgPo6sr5pEKsMPPTzeOce7huyz0ObRMnTdupH/IPC",
		AlipayPublicKey:      "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjHvLLIMpdbGr3xvT1rk+LXpkPbwyMgumQVTfbEC6h04mWdy2PnwargmzxJW6baKA8jJkG7nBGxwH8Py150vwai0jWPZnDnVTYOo3o1T0Sp1bNSJelOsMBok0+ELGK3Cc3A7xqnKL5Qz6gtrJBoqdSThXyLhh5H8b+0iuO4xgO/4B6w+th2/m3jEQMmVnhbr0e7M0zCYHNOOhNOsJ1QOm+jQSWNUB+hxXVHkRitcQV7UBVv2ebuqP0kDmEfbCb5mNKdz7MqhX6smPZECkX9iBWg17xy0VauSUIxPsVpvS7fyeAXm4lbDg5jcxa71MU2RQlln+kssZYWwWNeY/lRH7FQIDAQAB",
		AlipaySubject:        "",
		AlipayNotifyUrl:      "http://tvxw11009075.vicp.fun/o/v1/",
		AlipayReturnUrl:      "http://tvxw11009075.vicp.fun/o/v1/pay/callback",
		AlipayProductCode:    "FAST_INSTANT_TRADE_PAY", // 电脑网站支付默认产品码
		AlipayTimeoutExpress: "30m",                    // 默认30分钟链接失效
	}
}

// Validate 校验配置项的合法性
func (o *AliyunOptions) Validate() []error {
	errs := []error{}
	// 校验核心必填配置项，避免空值导致支付初始化失败
	if o.AlipayAppId == "" {
		errs = append(errs, fmt.Errorf("aliyun alipay app id cannot be empty"))
	}
	if o.AlipayPrivateKey == "" {
		errs = append(errs, fmt.Errorf("aliyun alipay private key cannot be empty"))
	}
	if o.AlipayPublicKey == "" {
		errs = append(errs, fmt.Errorf("aliyun alipay public key cannot be empty"))
	}
	return errs
}

// AddFlags 为配置项添加命令行参数（pflag）
func (o *AliyunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.AlipayAppId, "aliyun.alipay.app-id", o.AlipayAppId, "AliPay application ID (required)")
	fs.StringVar(&o.AlipayPrivateKey, "aliyun.alipay.private-key", o.AlipayPrivateKey, "AliPay application private key (required)")
	fs.StringVar(&o.AlipayPublicKey, "aliyun.alipay.public-key", o.AlipayPublicKey, "AliPay platform public key (required)")
	fs.StringVar(&o.AlipayNotifyUrl, "aliyun.alipay.notify-url", o.AlipayNotifyUrl, "AliPay asynchronous notify callback URL")
	fs.StringVar(&o.AlipayReturnUrl, "aliyun.alipay.return-url", o.AlipayReturnUrl, "AliPay synchronous redirect URL")
	fs.StringVar(&o.AlipayProductCode, "aliyun.alipay.product-code", o.AlipayProductCode, "AliPay product code (default: FAST_INSTANT_TRADE_PAY)")
	fs.StringVar(&o.AlipayTimeoutExpress, "aliyun.alipay.timeout-express", o.AlipayTimeoutExpress, "AliPay payment link timeout (default: 30m)")
}
