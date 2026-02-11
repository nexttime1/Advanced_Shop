package v1

import (
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"context"
	"gorm.io/gorm"
)

type BannerStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BannerList, error)
	Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error
	Update(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error
	Delete(ctx context.Context, ID uint64) error
}
