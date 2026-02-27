package v1

import (
	v1 "Advanced_Shop/app/inventory/srv/internal/data/v1"
	"Advanced_Shop/app/pkg/options"
	zlog "Advanced_Shop/pkg/log"
	"Advanced_Shop/pkg/storage"
	"context"
	"time"

	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type ServiceFactory interface {
	Inventories() InventorySrv
}

type service struct {
	data v1.DataFactory

	redisOptions *options.RedisOptions
	pool         redsyncredis.Pool
}

func (s *service) Inventories() InventorySrv {
	return newInventoryService(s)
}

func NewService(store v1.DataFactory, redisOptions *options.RedisOptions) ServiceFactory {
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

	// 返回service实例，只带这一个pool
	return &service{
		data:         store,
		redisOptions: redisOptions,
		pool:         pool, // 只用这一个复用的pool
	}
}

var _ ServiceFactory = &service{}
