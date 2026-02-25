package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

// RocketMQOptions RocketMQ 配置
type RocketMQOptions struct {
	Host              string `mapstructure:"host" yaml:"host"`
	Port              int    `mapstructure:"port" yaml:"port"`
	GroupName         string `mapstructure:"group_name" yaml:"group_name"`
	Topic             string `mapstructure:"topic" yaml:"topic"`
	ConsumerGroupName string `mapstructure:"consumer_group_name" yaml:"consumer_group_name"`
	ConsumerSubscribe string `mapstructure:"consumer_subscribe" yaml:"consumer_subscribe"`
	ConsumerTopic     string `mapstructure:"consumer_topic" yaml:"consumer_topic"`
	MaxRetryTimes     int    `mapstructure:"max_retry_times" yaml:"max_retry_times"`
	BaseRetryDelay    int    `mapstructure:"base_retry_delay" yaml:"base_retry_delay"`
}

func NewRocketMQOptions() *RocketMQOptions {
	return &RocketMQOptions{
		Host:              "127.0.0.1",
		Port:              9876,
		GroupName:         "goods_group",
		Topic:             "goods_topic",
		ConsumerGroupName: "goods_consumer_group",
		ConsumerSubscribe: "*",
		ConsumerTopic:     "goods_topic",
		MaxRetryTimes:     3,
		BaseRetryDelay:    1000,
	}
}

func (o *RocketMQOptions) Addr() string {
	return fmt.Sprintf("%s:%d", o.Host, o.Port)
}

func (o *RocketMQOptions) Validate() []error {
	var errs []error
	// 可添加地址、Topic 等非空校验
	return errs
}

func (o *RocketMQOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "rocketmq.host", o.Host, "RocketMQ host")
	fs.IntVar(&o.Port, "rocketmq.port", o.Port, "RocketMQ port")
	fs.StringVar(&o.GroupName, "rocketmq.group_name", o.GroupName, "RocketMQ producer group name")
	fs.StringVar(&o.Topic, "rocketmq.topic", o.Topic, "RocketMQ producer topic")
	fs.StringVar(&o.ConsumerGroupName, "rocketmq.consumer_group_name", o.ConsumerGroupName, "RocketMQ consumer group name")
	fs.StringVar(&o.ConsumerSubscribe, "rocketmq.consumer_subscribe", o.ConsumerSubscribe, "RocketMQ consumer subscribe expression")
	fs.StringVar(&o.ConsumerTopic, "rocketmq.consumer_topic", o.ConsumerTopic, "RocketMQ consumer topic")
	fs.IntVar(&o.MaxRetryTimes, "rocketmq.max_retry_times", o.MaxRetryTimes, "RocketMQ max retry times")
	fs.IntVar(&o.BaseRetryDelay, "rocketmq.base_retry_delay", o.BaseRetryDelay, "RocketMQ base retry delay (ms)")
}
