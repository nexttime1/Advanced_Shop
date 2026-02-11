package v1

import (
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"context"

	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type BrandsStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error)
	Create(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error
	Update(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error
	Delete(ctx context.Context, ID uint64) error
	Get(ctx context.Context, ID uint64) (*do.BrandsDO, error)
}
