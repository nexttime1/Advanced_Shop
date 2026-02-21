package db

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
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
	dbFactory v1.MysqlFactory
	once      sync.Once
)

type mysqlFactory struct {
	db *gorm.DB
}

func (mf *mysqlFactory) Begin() *gorm.DB {
	return mf.db.Begin()
}

func (mf *mysqlFactory) Goods() v1.GoodsStore {
	return newGoods(mf)
}

func (mf *mysqlFactory) Categorys() v1.CategoryStore {
	return newCategories(mf)
}

func (mf *mysqlFactory) Brands() v1.BrandsStore {
	return newBrands(mf)
}

func (mf *mysqlFactory) Banners() v1.BannerStore {
	return newBanner(mf)
}

func (m *mysqlFactory) CategoryBrands() v1.GoodsCategoryBrandStore {
	return NewCategoryBrands(m)
}

var _ v1.MysqlFactory = &mysqlFactory{}

// NewMySQLDataFactory 这个方法会返回gorm连接
func NewMySQLDataFactory(mysqlOpts *options.MySQLOptions) (v1.MysqlFactory, error) {
	// 入参校验
	if mysqlOpts == nil {
		return nil, fmt.Errorf("mysql配置不能为空")
	}

	var initErr error
	once.Do(func() {
		// 初始化MySQL连接
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlOpts.Username,
			mysqlOpts.Password,
			mysqlOpts.Host,
			mysqlOpts.Port,
			mysqlOpts.Database)

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.LogLevel(mysqlOpts.LogLevel),
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			initErr = errors2.WithCode(code2.ErrConnectDB, fmt.Sprintf("mysql连接失败: %v", err))
			return
		}

		// 设置MySQL连接池参数
		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)

		dbFactory = &mysqlFactory{
			db: db,
		}
	})

	// 初始化结果校验
	if dbFactory == nil || initErr != nil {
		return nil, initErr
	}
	return dbFactory, nil
}
