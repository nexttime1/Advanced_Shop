package mysql

import (
	"Advanced_Shop/app/inventory/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/code"
	zlog "Advanced_Shop/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
	"os"
	"sync"
	"time"

	v12 "Advanced_Shop/app/inventory/srv/internal/data/v1"
	"Advanced_Shop/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mysqlStore struct {
	db     *gorm.DB
	mqOpts *options.RocketMQOptions
}

func (m *mysqlStore) Inventorys() v12.InventoryStore {
	return newInventorys(m)
}

var _ v12.DataFactory = &mysqlStore{}

var (
	dbFactory v12.DataFactory
	once      sync.Once
)

// GetDBFactoryOr 对于复杂的初始化过程，使用工厂模式
func GetDBFactoryOr(mysqlOpts *options.MySQLOptions, mqOpts *options.RocketMQOptions) (v12.DataFactory, error) {
	if mysqlOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlOpts.Username,
			mysqlOpts.Password,
			mysqlOpts.Host,
			mysqlOpts.Port,
			mysqlOpts.Database)

		// 封装logger
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second,                         // 慢 SQL 阈值
				LogLevel:                  logger.LogLevel(mysqlOpts.LogLevel), // 日志级别
				IgnoreRecordNotFoundError: true,                                // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,                               // 禁用彩色打印
			},
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return
		}

		sqlDB, _ := db.DB()
		dbFactory = &mysqlStore{
			db:     db,
			mqOpts: mqOpts,
		}

		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)
	})

	if dbFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}
	return dbFactory, nil
}

func (ds *mysqlStore) Begin() *gorm.DB {
	return ds.db.Begin()
}

func (m *mysqlStore) Listen(ctx context.Context) {

	mqConsumer, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{m.mqOpts.Addr()}),
		consumer.WithGroupName(m.mqOpts.ConsumerGroupName),
		// 最大重试次数
		// -1表示使用默认值16次，这里显式设置为3次
		consumer.WithMaxReconsumeTimes(int32(m.mqOpts.MaxRetryTimes)),
		// 重试延迟时间
		// 每次重试时，当前队列暂停拉取的时间（对应原DelayLevelWhenNextConsume）
		// 级别2对应5秒，这里直接设置为5*time.Second
		consumer.WithSuspendCurrentQueueTimeMillis(time.Duration(m.mqOpts.BaseRetryDelay)),
		// 消费超时时间（可选，防止长耗时消费）
		consumer.WithConsumeTimeout(30*time.Second),
		// 并发消费协程数（可选，根据业务调整）
		consumer.WithConsumeGoroutineNums(10),
		// 批量消费大小（可选，每次最多消费1条，保证幂等）
		consumer.WithConsumeMessageBatchMaxSize(1),
	)
	if err != nil {
		panic(err)
	}

	// 订阅 Topi 普通归还消息（订单超时的库存归还）
	err = mqConsumer.Subscribe(m.mqOpts.ConsumerTopic, consumer.MessageSelector{}, m.AutoReBack)
	if err != nil {
		zlog.Errorf("订阅普通归还Topic失败 %v", err)
		panic(err)
	}

	// 启动消费者
	if err = mqConsumer.Start(); err != nil {
		zlog.Errorf("启动RocketMQ消费者失败 %v", err)
		panic(err)
	}

	// 优雅退出
	go func() {
		<-ctx.Done()
		// 用超时机制包装Shutdown，避免无限阻塞
		shutdownDone := make(chan error, 1)
		go func() {
			shutdownDone <- mqConsumer.Shutdown()
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
	zlog.Infof("库存服务MQ消费者启动成功，监听Topic： %v", m.mqOpts.ConsumerTopic)

}

func (m *mysqlStore) AutoReBack(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	zlog.Infof("开始处理库存归还消息，本次处理消息数： %d", len(msg))
	needRetry := false
	for i := range msg {
		var orderInfo do.OrderMQMessageRequest
		err := json.Unmarshal(msg[i].Body, &orderInfo)
		if err != nil {
			zlog.Error("解析订单超时消息体失败, 跳过该消息")
			// 解析失败：标记无需重试（消息本身有问题，重试也没用），继续处理下一条
			continue
		}
		// 开启事务
		txn := m.Begin()
		rollbackFlag := true // 标记当前消息是否需要回滚
		// 延迟处理：panic/未提交时回滚事务
		defer func() {
			if r := recover(); r != nil {
				zlog.Errorf("处理库存归还消息失败，OrderSn: %s, err: %v", orderInfo.OrderSns, r)
				txn.Rollback()
			} else if rollbackFlag {
				txn.Rollback()
			}
		}()
		number := m.Inventorys().AutoReback(ctx, txn, orderInfo.OrderSns)
		if number == do.OptionFail {
			needRetry = true
			continue
		}
		// 提交当前消息的事务
		err = txn.Commit().Error
		if err != nil {
			zlog.Error("提交库存归还事务失败")
			needRetry = true
			continue
		}
		rollbackFlag = false
		zlog.Info("库存归还处理成功")

	}
	if needRetry {
		// 有消息处理失败，返回重试（客户端会重新推送未处理成功的消息）
		return consumer.ConsumeRetryLater, nil
	}
	// 所有消息处理完成（成功/无需处理）
	return consumer.ConsumeSuccess, nil

}
