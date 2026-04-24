package mysql

import (
	v12 "Advanced_Shop/app/inventory/srv/internal/data/v1"
	"Advanced_Shop/app/inventory/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	zlog "Advanced_Shop/pkg/log"
	"Advanced_Shop/pkg/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

type mysqlStore struct {
	db     *gorm.DB
	mqOpts *options.RocketMQOptions
	pool   redsyncredis.Pool
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

		// 第一步：复用storage层已经配置好的Redis客户端（带密码、超时等）
		// 注意：这里要从storage层获取已初始化的客户端，而不是重新创建
		redisCluster := &storage.RedisCluster{} // 实例化storage的RedisCluster
		redisClient := redisCluster.GetClient() // 获取已认证的Redis客户端
		if redisClient == nil {
			zlog.Fatal("无法从storage层获取Redis客户端（请确保storage.ConnectToRedis已执行）")
		}

		// 第二步：验证客户端是否能正常认证（可选，但建议保留）
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			zlog.Fatalf("Redis客户端认证失败：%v，请检查密码配置", err)
		}
		zlog.Info("✅ Redis客户端复用成功，认证通过")

		// 第三步：基于已认证的客户端创建redsync的Pool（仅此一个即可）
		pool := goredis.NewPool(redisClient.(redis.UniversalClient))

		sqlDB, _ := db.DB()
		dbFactory = &mysqlStore{
			db:     db,
			mqOpts: mqOpts,
			pool:   pool,
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

func (ds *mysqlStore) Pool() redsyncredis.Pool {
	return ds.pool
}

func (ds *mysqlStore) DB() *gorm.DB {
	return ds.db
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

// AutoReBack ==================== MQ消费入口 ====================
func (m *mysqlStore) AutoReBack(ctx context.Context, msg ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	zlog.Infof("开始处理库存归还消息，本次消息数: %d", len(msg))
	needRetry := false

	for i := range msg {
		// 每条消息独立处理，互不影响
		if err := m.processOneMessage(ctx, msg[i]); err != nil {
			zlog.Errorf("消息处理失败，msgId: %s, err: %v", msg[i].MsgId, err)
			needRetry = true
			// 继续处理下一条，不能因为一条失败影响其他
		}
	}

	if needRetry {
		return consumer.ConsumeRetryLater, nil
	}
	return consumer.ConsumeSuccess, nil
}

// ==================== 单条消息处理 ====================
func (m *mysqlStore) processOneMessage(ctx context.Context, msg *primitive.MessageExt) error {

	// 解析消息体
	var orderInfo do.OrderMQMessageRequest
	if err := json.Unmarshal(msg.Body, &orderInfo); err != nil {
		// 消息体损坏，记录告警，直接丢弃（重试无意义）
		zlog.Errorf("消息体解析失败，msgId: %s, body: %s, err: %v", msg.MsgId, string(msg.Body), err)
		return nil // 返回nil表示不重试，避免死信队列积压
	}

	zlog.Infof("开始处理订单库存归还，OrderSn: %s", orderInfo.OrderSns)

	// 开启事务
	txn := m.Begin()
	if txn.Error != nil {
		return fmt.Errorf("开启事务失败: %w", txn.Error)
	}

	// 事务闭包，确保一定被处理
	committed := false
	defer func() {
		if !committed {
			if err := txn.Rollback().Error; err != nil {
				zlog.Errorf("事务回滚失败，OrderSn: %s, err: %v",
					orderInfo.OrderSns, err)
			} else {
				zlog.Warnf("事务已回滚，OrderSn: %s", orderInfo.OrderSns)
			}
		}
	}()

	// 核心业务逻辑
	result, err := m.Inventorys().AutoReback(ctx, txn, orderInfo.OrderSns, m.pool)

	if err != nil {
		// 业务失败，defer会自动回滚
		return fmt.Errorf("库存归还业务失败，OrderSn: %s, err: %w", orderInfo.OrderSns, err)
	}

	// 幂等：已处理过，直接跳过，但要提交（避免没必要的回滚日志）
	if result == do.DirectPass {
		committed = true // 标记为"不需要回滚"，幂等成功
		zlog.Infof("订单库存归还已处理过，跳过，OrderSn: %s", orderInfo.OrderSns)
		return nil
	}

	// 6. 提交事务
	if err := txn.Commit().Error; err != nil {
		// defer会回滚
		return fmt.Errorf("提交事务失败，OrderSn: %s, err: %w",
			orderInfo.OrderSns, err)
	}
	committed = true
	zlog.Infof("库存归还成功，OrderSn: %s", orderInfo.OrderSns)
	return nil
}
