package v1

import (
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"context"
)

// AddressStore 地址数据访问层接口
type AddressStore interface {
	// GetByID 根据ID和用户ID获取地址信息
	GetByID(ctx context.Context, ID uint, userID int32) (*do.AddressDO, error)

	// ListByUserID 根据用户ID获取地址列表
	ListByUserID(ctx context.Context, userID int32) ([]*do.AddressDO, error)

	// Create 创建新地址
	Create(ctx context.Context, address *do.AddressDO) error

	// Update 更新地址信息
	Update(ctx context.Context, address *do.AddressDO) error

	// Delete 根据ID和用户ID删除地址
	Delete(ctx context.Context, ID uint, userID int32) error
}
