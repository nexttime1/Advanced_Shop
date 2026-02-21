package v1

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"gorm.io/gorm"
)

type MysqlFactory interface {
	Goods() GoodsStore
	Categorys() CategoryStore
	Brands() BrandsStore
	Banners() BannerStore
	CategoryBrands() GoodsCategoryBrandStore
	Begin() *gorm.DB
}

type CanalFactory interface {
	ParseCanalMessage() ([]*pbe.Entry, error)
}
type MQFactory interface {
	BuildGoodsMQMessage(eventType pbe.EventType, rowData *pbe.RowData, header *pbe.Header) (*primitive.Message, error)
	Send(ctx context.Context, mqMsg *primitive.Message) (*primitive.SendResult, error)
}

type DataFactory interface {
	NewMysql() MysqlFactory
	NewCanal() CanalFactory
	NewMQ() MQFactory
	StartCanalListener(context.Context)
}
