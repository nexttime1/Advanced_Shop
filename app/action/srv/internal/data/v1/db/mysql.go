package db

import (
	proto "Advanced_Shop/api/goods/v1"
	v1 "Advanced_Shop/app/action/srv/internal/data/v1"
	"Advanced_Shop/app/pkg/options"
	code2 "Advanced_Shop/gnova/code"
	errors2 "Advanced_Shop/pkg/errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
)

var (
	dbFactory v1.DataFactory
	once      sync.Once
)

type mysqlFactory struct {
	gc proto.GoodsClient
	db *gorm.DB
}

func (mf *mysqlFactory) Messages() v1.MessageStore {
	return newMessage(mf)
}

func (mf *mysqlFactory) Address() v1.AddressStore {
	return newAddress(mf)
}

func (mf *mysqlFactory) Collection() v1.CollectionStore {
	return newCollection(mf)
}

func (mf *mysqlFactory) Goods() proto.GoodsClient {
	return mf.gc
}

var _ v1.DataFactory = &mysqlFactory{}

// GetDBFactoryOr 这个方法会返回gorm连接
func GetDBFactoryOr(mysqlOpts *options.MySQLOptions, registry *options.RegistryOptions) (v1.DataFactory, error) {
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

		//希望大家自己可以去封装logger
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

		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)

		//服务发现
		goodsClient := GetGoodsClient(registry)
		dbFactory = &mysqlFactory{
			db: db,
			gc: goodsClient,
		}

	})

	if dbFactory == nil || err != nil {
		return nil, errors2.WithCode(code2.ErrConnectDB, "failed to get mysql store factory")
	}
	return dbFactory, nil
}
