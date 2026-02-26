package realize

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/data/v1/canal"
	"Advanced_Shop/app/goods/srv/internal/data/v1/db"
	"Advanced_Shop/app/goods/srv/internal/data/v1/mq"
	"Advanced_Shop/app/pkg/options"
	zlog "Advanced_Shop/pkg/log"
	"context"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"google.golang.org/protobuf/proto"
	"time"
)

type DataStore struct {
	mysqlOpts *options.MySQLOptions
	mqOpts    *options.RocketMQOptions
	canalOpts *options.CanalOptions
}

func NewDataStore(
	mysqlOpts *options.MySQLOptions,
	mqOpts *options.RocketMQOptions,
	canalOpts *options.CanalOptions) v1.DataFactory {
	return &DataStore{
		mysqlOpts: mysqlOpts,
		mqOpts:    mqOpts,
		canalOpts: canalOpts,
	}
}

func (store *DataStore) NewMysql() v1.MysqlFactory {
	factory, err := db.NewMySQLDataFactory(store.mysqlOpts)
	if err != nil {
		panic(err)
	}
	return factory
}
func (store *DataStore) NewCanal() v1.CanalFactory {
	factory, err := canal.NewCanalFactory(store.canalOpts)
	if err != nil {
		panic(err)
	}
	return factory
}
func (store *DataStore) NewMQ() v1.MQFactory {
	factory, err := mq.NewMQFactory(store.mqOpts)
	if err != nil {
		panic(err)
	}
	return factory
}

func (store *DataStore) StartCanalListener(ctx context.Context) {
	go func() {
		zlog.Info("Canal监听器启动成功，开始监听商品表binlog")
		//  构建MQ消息并发送
		mqData := store.NewMQ()
		for {
			select {
			case <-ctx.Done():
				zlog.Info("Canal监听器停止")
				return
			default:
				entries, err := store.NewCanal().ParseCanalMessage()
				if err != nil {
					zlog.Errorf("Canal获取消息失败, err=%v", err)
					time.Sleep(300 * time.Millisecond)
					continue
				}
				if len(entries) == 0 {
					time.Sleep(300 * time.Millisecond)
					continue
				}

				//  解析并处理商品表binlog
				for _, entry := range entries {
					// 过滤非商品表/非数据变更事件
					if entry.GetEntryType() != pbe.EntryType_ROWDATA {
						continue
					}
					header := entry.GetHeader()
					if header.GetTableName() != store.canalOpts.TableName { // 只处理商品表
						continue
					}

					// 解析RowChange（binlog内容）
					rowChange := &pbe.RowChange{}
					err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
					if err != nil {
						zlog.Errorf("解析binlog失败, table=%s, err=%v", header.GetTableName(), err)
						continue
					}

					// 只处理INSERT/UPDATE事件（同步到ES）
					eventType := rowChange.GetEventType()
					if eventType != pbe.EventType_INSERT && eventType != pbe.EventType_UPDATE {
						continue
					}

					for _, rowData := range rowChange.GetRowDatas() {
						mqMsg, err := mqData.BuildGoodsMQMessage(eventType, rowData, header)
						if err != nil {
							zlog.Errorf("构建MQ消息失败, goodsID=%v, err=%v", rowData.GetAfterColumns(), err)
							continue
						}
						goodsMap := make(map[string]interface{})
						for _, col := range rowData.GetAfterColumns() {
							goodsMap[col.GetName()] = col.GetValue()
						}

						// 发送MQ消息
						result, err := mqData.Send(ctx, mqMsg)
						if err != nil {
							zlog.Errorf("发送MQ消息失败, goods内容= %v, err=%v", goodsMap, err)
							continue
						}

						zlog.Infof("发送商品MQ消息成功, goods内容 =%v, msgID=%s", goodsMap, result.MsgID)
					}
				}
			}
		}
	}()
}
