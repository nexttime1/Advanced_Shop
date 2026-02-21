package v1

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type SearchFactory interface {
	Goods() GoodsStore
	Listen(ctx context.Context) error
	SyncGoodsToES(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
	Close() error
}
