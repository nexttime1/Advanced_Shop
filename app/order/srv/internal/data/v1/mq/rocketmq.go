package mq

import (
	v1 "Advanced_Shop/app/order/srv/internal/data/v1"
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
	"go.uber.org/zap"
	"sync"
	"time"
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
	zlog.Info("mq 开始发送给库存服务  超时的消息")
	result, err := mf.producer.SendSync(ctx, mqMsg)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SendDelayMsgWithRetry 发送延迟消息  带重试
func (mf *RocketMqFactory) SendDelayMsgWithRetry(ctx context.Context, msg *primitive.Message) (*primitive.SendResult, error) {
	var (
		sendResult *primitive.SendResult
		sendErr    error
		// 发送延时消息 就用 Topic  跨服务给 库存服务 就用 CrossTopic
		delayMsg = primitive.NewMessage(mf.mqOpts.Topic, msg.Body)
	)
	// 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
	delayMsg.WithDelayTimeLevel(7) // 30分钟延时

	var retryIdx int
	// 执行重试逻辑
	for retryIdx = 0; retryIdx < mf.mqOpts.MaxRetryTimes; retryIdx++ {
		// 发送消息
		sendResult, sendErr = mf.producer.SendSync(context.Background(), delayMsg)

		if sendErr != nil {
			zlog.Errorf("SendDelayMsgWithRetry send error: %v", sendErr)
		}

		// 发送成功：记录日志并返回
		if sendResult != nil && sendErr == nil && sendResult.Status == primitive.SendOK {
			zlog.Info("延时消息发送成功")
			return sendResult, nil
		}

		// 发送失败：记录日志，判断是否继续重试
		if retryIdx == mf.mqOpts.MaxRetryTimes-1 {
			// 最后一次重试失败：记录错误日志（标记最终失败）
			zlog.Errorf("延时消息3次发送均失败")
		} else {
			// 非最后一次：记录警告日志，等待后重试
			zap.S().Warn("延时消息发送失败，准备重试")
			// 递增间隔重试（避免高频重试）：第1次等500ms，第2次等1000ms
			baseRetryDelay := time.Duration(mf.mqOpts.BaseRetryDelay) * time.Millisecond

			// 2. 再计算递增的重试间隔
			retryDelay := baseRetryDelay * time.Duration(retryIdx+1)
			time.Sleep(retryDelay)
		}
	}

	// 3次重试均失败：返回最终结果
	return sendResult, sendErr
}
