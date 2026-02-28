package realize

import (
	v1 "Advanced_Shop/app/order/srv/internal/data/v1"
	"Advanced_Shop/app/order/srv/internal/data/v1/db"
	"Advanced_Shop/app/order/srv/internal/data/v1/mq"
	"Advanced_Shop/app/order/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/options"
	zlog "Advanced_Shop/pkg/log"
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"time"
)

type dataFactory struct {
	mqOpts    *options.RocketMQOptions
	mysqlOpts *options.MySQLOptions
	registry  *options.RegistryOptions
}

func NewDataFactory(mysqlOpts *options.MySQLOptions, registry *options.RegistryOptions, mqOpts *options.RocketMQOptions) v1.DataFactory {
	d := &dataFactory{
		mqOpts:    mqOpts,
		mysqlOpts: mysqlOpts,
		registry:  registry,
	}
	// 初始化一下
	d.NewDB()
	d.NewMQ()
	return d
}

func (d *dataFactory) NewDB() v1.DBFactory {
	factory, err := db.NewDBFactoryOr(d.mysqlOpts, d.registry)
	if err != nil {
		panic(err)
	}
	return factory
}

func (d *dataFactory) NewMQ() v1.MQFactory {
	factory, err := mq.NewMQFactory(d.mqOpts)
	if err != nil {
		panic(err)
	}
	return factory
}

func (d *dataFactory) Listen(ctx context.Context) {
	messgaes, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{d.mqOpts.Addr()}),
		consumer.WithGroupName(d.mqOpts.ConsumerGroupName),
		// 最大重试次数
		consumer.WithMaxReconsumeTimes(int32(d.mqOpts.MaxRetryTimes)),
		// 重试延迟时间
		consumer.WithSuspendCurrentQueueTimeMillis(time.Duration(d.mqOpts.BaseRetryDelay)),
	)
	if err != nil {
		panic(err)
	}
	// 监听普通消息 订单时间超时   Topic 订阅这个
	messgaes.Subscribe(d.mqOpts.Topic, consumer.MessageSelector{}, d.Timeout)
	messgaes.Start()

	go func() {
		<-ctx.Done()
		// 用超时机制包装Shutdown，避免无限阻塞
		shutdownDone := make(chan error, 1)
		go func() {
			shutdownDone <- messgaes.Shutdown()
		}()

		// 设置10秒超时，避免Shutdown卡住
		select {
		case err := <-shutdownDone:
			if err != nil {
				// 替换为你的日志组件
				zlog.Errorf("RocketMQ消费者关闭失败: %v\n", err)
			} else {
				zlog.Info("RocketMQ消费者已优雅关闭")
			}
		case <-time.After(10 * time.Second):
			zlog.Error("RocketMQ消费者关闭超时（10秒），强制退出")
		}
	}()

	zlog.Info("RocketMQ消费者启动成功，开始监听消息")

}

func (d *dataFactory) Timeout(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	zap.S().Info("订单超时")
	// 只要一个需要重试就得重试
	needRetry := false
	for i := range msg {
		var orderInfo do.OrderMQMessageRequest
		err := json.Unmarshal(msg[i].Body, &orderInfo)
		// 需要查看 订单存不存在
		if err != nil {
			zlog.Error("解析订单超时消息体失败, 跳过该消息")
			// 解析失败：标记无需重试（消息本身有问题，重试也没用），继续处理下一条
			continue
		}
		tx := d.NewDB().Begin()
		// 标记当前消息是否处理成功
		success := true
		defer func() {
			if !success {
				tx.Rollback()
			}
		}()
		txn := d.NewDB().Begin()

		number := d.NewDB().Orders().TimeoutHandler(ctx, txn, orderInfo.OrderSns)
		// 0 代表没找到 说明没有 这样就不需要管了 因为都没创建订单 所以不需要归还库存
		// 返回1是有错需要重试保存状态失败
		//返回2 是继续向下走 说明要发消息 进行回收
		if number == do.DirectPass {
			zlog.Warn("未进行创建 直接跳过")
			continue
		} else if number == do.OptionFail {
			needRetry = true
			continue
		}
		// 下面就是没啥问题 正常流程 进行库存归还    这个时候要发到  CrossTopic 这个 topic
		_, err = d.NewMQ().Send(ctx, primitive.NewMessage(d.mqOpts.CrossTopic, msg[i].Body))
		if err != nil {
			// 这个时候就需要回滚了
			success = false
			needRetry = true
			continue
		}
		err = txn.Commit().Error
		if err != nil {
			zlog.Error("提交订单超时事务失败")
			success = false
			needRetry = true
			continue
		}
		zlog.Info("订单超时处理成功")
	}
	// 根据是否有失败消息，返回整体消费结果
	if needRetry {
		// 有消息处理失败，返回重试（客户端会重新推送所有未处理成功的消息）
		return consumer.ConsumeRetryLater, nil
	}
	// 所有消息处理完成（成功/无需处理），返回成功
	return consumer.ConsumeSuccess, nil
}
