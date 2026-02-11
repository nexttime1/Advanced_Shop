package v1

import (
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"context"
	"gorm.io/gorm"
)

type GoodsStore interface {
	Get(ctx context.Context, ID uint64) (*do.GoodsDO, error)
	ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error)
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.GoodsDOList, error)
	Create(ctx context.Context, goods *do.GoodsDO) error
	CreateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Update(ctx context.Context, goods *do.GoodsDO) error
	UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Delete(ctx context.Context, ID uint64) error
	DeleteInTxn(ctx context.Context, txn *gorm.DB, ID uint64) error

	Begin() *gorm.DB
}
