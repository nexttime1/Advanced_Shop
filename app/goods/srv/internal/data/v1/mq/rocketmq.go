package mq

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/pkg/options"
	code2 "Advanced_Shop/gnova/code"
	errors2 "Advanced_Shop/pkg/errors"
	zlog "Advanced_Shop/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"sync"
)

var (
	mqFactory v1.MQFactory
	once      sync.Once
)

type RocketMqFactory struct {
	mqOpts   *options.RocketMQOptions
	producer rocketmq.Producer // RocketMQ 生产者
}

func NewMQFactory(mqOpts *options.RocketMQOptions) (v1.MQFactory, error) {
	zlog.Debug("NewMQFactory")
	if mqOpts == nil {
		return nil, fmt.Errorf("rocketmq配置不能为空")
	}

	var initErr error

	// 初始化 RocketMQ 生产者
	once.Do(func() {
		addr := fmt.Sprintf("%s:%d", mqOpts.Host, mqOpts.Port)
		zlog.Infof("rocketmq addr:%s", addr)
		producerIns, err := rocketmq.NewProducer(
			producer.WithNameServer([]string{addr}),
			producer.WithGroupName(mqOpts.GroupName),
		)
		if err != nil {
			initErr = errors2.WithCode(code2.ErrConnectMQ, fmt.Sprintf("rocketmq生产者创建失败: %v", err))
			return
		}

		// 启动生产者
		if err = producerIns.Start(); err != nil {
			initErr = errors2.WithCode(code2.ErrConnectMQ, fmt.Sprintf("rocketmq生产者启动失败: %v", err))
			return
		}
		mqFactory = &RocketMqFactory{
			mqOpts:   mqOpts,
			producer: producerIns,
		}

		zlog.Infof("RocketMQ生产者初始化成功 topic: %v", mqOpts.Topic)
	})
	// 初始化结果校验
	if mqFactory == nil || initErr != nil {
		return nil, initErr
	}
	return mqFactory, nil
}

// BuildGoodsMQMessage 构建商品MQ消息体
func (mf *RocketMqFactory) BuildGoodsMQMessage(
	eventType pbe.EventType,
	rowData *pbe.RowData,
	header *pbe.Header,
) (*primitive.Message, error) {
	// 解析商品字段（从AfterColumns获取最新数据）
	goodsMap := make(map[string]interface{})
	for _, col := range rowData.GetAfterColumns() {
		goodsMap[col.GetName()] = col.GetValue()
	}

	// 构建MQ消息体（JSON格式）
	msgBody, err := json.Marshal(map[string]interface{}{
		"event_type": eventType.String(),
		"table":      header.GetTableName(),
		"schema":     header.GetSchemaName(),
		"goods":      goodsMap,
		"timestamp":  header.GetExecuteTime(),
	})
	if err != nil {
		return nil, errors2.WithCode(code2.ErrEncodingJSON, fmt.Sprintf("序列化MQ消息失败: %v", err))
	}

	// 创建RocketMQ消息
	msg := primitive.NewMessage(
		mf.mqOpts.Topic, // MQ Topic
		msgBody,
	)

	// 设置消息Key（去重/追溯）
	goodsID := goodsMap["id"]
	if goodsID != nil {
		msg.WithKeys([]string{fmt.Sprintf("goods_%v", goodsID)})
	}

	// 设置重试策略
	msg.WithDelayTimeLevel(mf.mqOpts.MaxRetryTimes)

	return msg, nil
}

func (mf *RocketMqFactory) Send(ctx context.Context, mqMsg *primitive.Message) (*primitive.SendResult, error) {
	zlog.Info("mq 开始发送关于es的消息")
	result, err := mf.producer.SendSync(ctx, mqMsg)
	if err != nil {
		return nil, err
	}
	return result, nil
}
