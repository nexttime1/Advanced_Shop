package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

// CanalOptions Canal 客户端配置（全参数配置化，无硬编码）
type CanalOptions struct {
	Addr                string `mapstructure:"addr" json:"addr,omitempty"`               // Canal服务IP
	Port                int    `mapstructure:"port" json:"port,omitempty"`               // Canal服务端口
	Username            string `mapstructure:"username" json:"username,omitempty"`       // Canal用户名（无则空）
	Password            string `mapstructure:"password" json:"password,omitempty"`       // Canal密码（无则空）
	Destination         string `mapstructure:"destination" json:"destination,omitempty"` // Canal实例名（如example）
	TableName           string `mapstructure:"table_name" json:"table_name,omitempty"`
	SubscribeRegex      string `mapstructure:"subscribe_regex" json:"subscribe_regex,omitempty"`             // binlog订阅正则（默认.*\\..*）
	BatchSize           int32  `mapstructure:"batch_size" json:"batch_size,omitempty"`                       // 批量获取消息大小（默认100）
	TimeoutMs           int32  `mapstructure:"timeout_ms" json:"timeout_ms,omitempty"`                       // 连接超时时间（毫秒，默认60000）
	HeartbeatIntervalMs int32  `mapstructure:"heartbeat_interval_ms" json:"heartbeat_interval_ms,omitempty"` // 心跳间隔（毫秒，默认3600000=1小时）
}

// NewCanalOptions 创建Canal配置实例（默认值与官方示例保持一致）
func NewCanalOptions() *CanalOptions {
	return &CanalOptions{
		Addr:                "127.0.0.1",    // 默认IP
		Port:                11111,          // 官方默认端口
		Username:            "",             // 官方示例无用户名
		Password:            "",             // 官方示例无密码
		Destination:         "example",      // 官方示例的destination
		SubscribeRegex:      ".*\\..*",      // 默认订阅所有库表
		BatchSize:           100,            // 官方示例批量获取100条
		TimeoutMs:           60000,          // 官方示例超时60秒
		HeartbeatIntervalMs: 60 * 60 * 1000, // 官方示例心跳1小时（3600000毫秒）
	}
}

// Validate 配置校验（可扩展，如端口范围、非空校验）
func (o *CanalOptions) Validate() []error {
	var errs []error
	// 示例：端口范围校验
	if o.Port <= 0 || o.Port > 65535 {
		errs = append(errs, fmt.Errorf("canal port %d is invalid (must 1-65535)", o.Port))
	}
	// 示例：超时时间校验
	if o.TimeoutMs <= 0 {
		errs = append(errs, fmt.Errorf("canal timeout_ms %d must be positive", o.TimeoutMs))
	}
	return errs
}

// AddFlags 将配置绑定到命令行参数（全参数可通过命令行覆盖）
func (o *CanalOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Addr, "canal.addr", o.Addr, "Canal server IP address (default 127.0.0.1)")
	fs.IntVar(&o.Port, "canal.port", o.Port, "Canal server port (default 11111)")
	fs.StringVar(&o.Username, "canal.username", o.Username, "Canal username (empty if not set)")
	fs.StringVar(&o.Password, "canal.password", o.Password, "Canal password (empty if not set)")
	fs.StringVar(&o.Destination, "canal.destination", o.Destination, "Canal destination name (instance name, default example)")
	fs.StringVar(&o.SubscribeRegex, "canal.subscribe_regex", o.SubscribeRegex, "Canal binlog subscribe regex (Perl style, default .*\\..*)")
	fs.Int32Var(&o.BatchSize, "canal.batch_size", o.BatchSize, "Canal batch size for fetching messages (default 100)")
	fs.Int32Var(&o.TimeoutMs, "canal.timeout_ms", o.TimeoutMs, "Canal connection timeout in milliseconds (default 60000)")
	fs.Int32Var(&o.HeartbeatIntervalMs, "canal.heartbeat_interval_ms", o.HeartbeatIntervalMs, "Canal heartbeat interval in milliseconds (default 3600000)")
}
