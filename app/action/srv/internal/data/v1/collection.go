package v1

import (
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"context"
)

// CollectionStore 用户收藏数据访问层接口
type CollectionStore interface {
	// ListByUserID 根据用户ID获取收藏列表（支持商品ID筛选）
	ListByUserID(ctx context.Context, userID int32, goodID int32) ([]*do.UserCollectionDO, int64, error)

	// Create 创建用户收藏
	Create(ctx context.Context, collection *do.UserCollectionDO) error

	// Delete 根据用户ID和商品ID删除收藏
	Delete(ctx context.Context, userID int32, goodID int32) (int64, error)

	// GetByUserAndGoodID 检查用户是否收藏了某个商品
	GetByUserAndGoodID(ctx context.Context, userID int32, goodID int32) (*do.UserCollectionDO, error)
}
