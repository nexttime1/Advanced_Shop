package db

import (
	"context"

	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type brands struct {
	db *gorm.DB
}

func newBrands(factory *mysqlFactory) *brands {
	brands := &brands{
		db: factory.db,
	}
	return brands
}

//func NewBrands(db *gorm.DB) *brands {
//	return &brands{
//		db: db,
//	}
//}

func (b *brands) Get(ctx context.Context, ID uint64) (*do.BrandsDO, error) {
	//TODO implement me
	panic("implement me")
}

func (b *brands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error) {
	//TODO implement me
	panic("implement me")
}

func (b *brands) Create(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	//TODO implement me
	panic("implement me")
}

func (b *brands) Update(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	//TODO implement me
	panic("implement me")
}

func (b *brands) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

var _ v1.BrandsStore = &brands{}
