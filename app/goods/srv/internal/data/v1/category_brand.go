package v1

import (
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"context"

	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type GoodsCategoryBrandStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandList, error)
	Create(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error
	Update(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error
	Delete(ctx context.Context, ID uint64) error
}
