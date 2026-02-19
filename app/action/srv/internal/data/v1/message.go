package v1

import (
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"context"
)

// MessageStore 留言数据访问层接口
type MessageStore interface {
	// ListByUserID 根据用户ID获取留言列表
	ListByUserID(ctx context.Context, userID int32) ([]*do.LeavingMessageDO, int64, error)

	// Create 创建留言
	Create(ctx context.Context, message *do.LeavingMessageDO) error
}
