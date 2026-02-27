package v1

import (
	proto "Advanced_Shop/api/goods/v1"
	proto2 "Advanced_Shop/api/inventory/v1"
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"gorm.io/gorm"
)

type DBFactory interface {
	Orders() OrderStore
	ShopCarts() ShopCartStore
	Goods() proto.GoodsClient
	Inventorys() proto2.InventoryClient

	Begin() *gorm.DB
}

type MQFactory interface {
	BuildGoodsMQMessage(eventType pbe.EventType, rowData *pbe.RowData, header *pbe.Header) (*primitive.Message, error)
	Send(ctx context.Context, mqMsg *primitive.Message) (*primitive.SendResult, error)
	SendDelayMsgWithRetry(ctx context.Context, msg *primitive.Message) (*primitive.SendResult, error)
	//Listen()
}

type DataFactory interface {
	NewDB() DBFactory
	NewMQ() MQFactory
	Listen(ctx context.Context)
}
