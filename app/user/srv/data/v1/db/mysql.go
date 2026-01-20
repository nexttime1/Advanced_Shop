package db

import (
	"fmt"
	"gorm.io/gorm"
	"mxshop/app/pkg/code"
	"mxshop/app/pkg/options"
	errors2 "mxshop/pkg/errors"
	"sync"

	"gorm.io/driver/mysql"
)

var (
	dbFactory *gorm.DB
	once      sync.Once
)

// 这个方法会返回gorm连接
// 还不够
// 这个方法应该返回的是全局的一个变量，如果一开始的时候没有初始化好，那么就初始化一次，后续呢直接拿到这个变量
func GetDBFactoryOr(mysqlOpts *options.MySQLOptions) (*gorm.DB, error) {
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
		dbFactory, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}

		sqlDB, _ := dbFactory.DB()

		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)
	})

	if dbFactory == nil || err != nil {
		return nil, errors2.WithCode(code.ErrConnectDB, "failed to get mysql store factory")
	}
	return dbFactory, nil
}
