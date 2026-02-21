package canal

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/pkg/options"
	code2 "Advanced_Shop/gnova/code"
	errors2 "Advanced_Shop/pkg/errors"
	zlog "Advanced_Shop/pkg/log"
	"fmt"
	canalClient "github.com/withlin/canal-go/client"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"sync"
)

var (
	canalFactory v1.CanalFactory
	once         sync.Once
)

type CanFactory struct {
	canalOpts   *options.CanalOptions
	canalClient *canalClient.SimpleCanalConnector // Canal 客户端
}

func NewCanalFactory(canalOpts *options.CanalOptions) (v1.CanalFactory, error) {
	if canalOpts == nil {
		return nil, fmt.Errorf("canal配置不能为空")
	}

	var initErr error

	// 初始化Canal客户端
	once.Do(func() {
		// NewSimpleCanalConnector 参数说明：
		// 参数1: Canal服务IP → canalOpts.Addr
		// 参数2: Canal服务端口 → canalOpts.Port
		// 参数3: Canal用户名 → canalOpts.Username
		// 参数4: Canal密码 → canalOpts.Password
		// 参数5: Canal实例名 → canalOpts.Destination
		// 参数6: 超时时间（毫秒） → canalOpts.TimeoutMs
		// 参数7: 心跳间隔（毫秒） → canalOpts.HeartbeatIntervalMs
		canalConn := canalClient.NewSimpleCanalConnector(
			canalOpts.Addr,
			canalOpts.Port,
			canalOpts.Username,
			canalOpts.Password,
			canalOpts.Destination,
			canalOpts.TimeoutMs,
			canalOpts.HeartbeatIntervalMs,
		)

		// 连接Canal服务
		if err := canalConn.Connect(); err != nil {
			initErr = errors2.WithCode(code2.ErrConnectCanal, fmt.Sprintf("canal服务连接失败: %v", err))
			return
		}

		// 订阅binlog
		// 建议只订阅商品库表，避免全量订阅：如 "advanced_shop\\.goods"（库名.表名）  // good-srv.good_models
		if err := canalConn.Subscribe(canalOpts.SubscribeRegex); err != nil {
			initErr = errors2.WithCode(code2.ErrCanalSubscribe, fmt.Sprintf("canal binlog订阅失败: %v", err))
			return
		}
		canalFactory = &CanFactory{
			canalOpts:   canalOpts,
			canalClient: canalConn,
		}
		zlog.Infof("Canal客户端初始化成功 destination: %v", canalOpts.Destination)
	})
	// 初始化结果校验
	if canalFactory == nil || initErr != nil {
		return nil, initErr
	}
	return canalFactory, nil
}

func (mf *CanFactory) ParseCanalMessage() ([]*pbe.Entry, error) {
	if mf.canalClient == nil {
		return nil, errors2.WithCode(code2.ErrConnectCanal, "canal客户端未初始化")
	}

	// 批量获取大小使用配置项，不再硬编码100
	message, err := mf.canalClient.Get(mf.canalOpts.BatchSize, nil, nil)
	if err != nil {
		return nil, errors2.WithCode(code2.ErrCanalGetData, fmt.Sprintf("canal获取binlog数据失败: %v", err))
	}

	// 无数据时返回空切片
	batchId := message.Id
	if batchId == -1 || len(message.Entries) <= 0 {
		return []*pbe.Entry{}, nil
	}

	// 过滤事务开始/结束事件，返回有效数据
	var validEntries []*pbe.Entry
	for _, entry := range message.Entries {
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		validEntries = append(validEntries, &entry)
	}

	return validEntries, nil
}
