package v1

import (
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"context"
	"gorm.io/gorm"
)

type GoodsInfo struct {
	GoodsDO         do.GoodsDO
	Images          []string `json:"images,omitempty"`
	DescImages      []string `json:"descImages,omitempty"`
	GoodsFrontImage string   `json:"goodsFrontImage,omitempty"`
}

type GoodsStore interface {
	Get(ctx context.Context, ID uint64) (*do.GoodsDO, error)
	ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error)
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.GoodsDOList, error)
	Create(ctx context.Context, goods *GoodsInfo) error
	CreateInTxn(ctx context.Context, txn *gorm.DB, goods *GoodsInfo) error
	Update(ctx context.Context, goods *GoodsInfo) error
	UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *GoodsInfo) error
	Delete(ctx context.Context, ID uint64) error
	DeleteInTxn(ctx context.Context, txn *gorm.DB, ID uint64) error

	Begin() *gorm.DB
}
