package v1

import (
	"Advanced_Shop/app/order/srv/internal/domain/do"
	"context"
	"gorm.io/gorm"

	metav1 "Advanced_Shop/pkg/common/meta/v1"
)

type ShopCartStore interface {
	List(ctx context.Context, userID uint64, checked bool, meta metav1.ListMeta, orderby []string) (*do.ShoppingCartDOList, error)
	Create(ctx context.Context, cartItem *do.ShoppingCartDO) error
	Get(ctx context.Context, userID, goodsID uint64) (*do.ShoppingCartDO, error)
	UpdateNum(ctx context.Context, cartItem *do.ShoppingCartDO) error
	Delete(ctx context.Context, ID uint64) error
	ClearCheck(ctx context.Context, userID uint64) error
	GetBatchByUser(ctx context.Context, userID int32) (*do.GetShoppingBatchResponse, error)
	DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID uint64, goodsIDs []int32) error
}
