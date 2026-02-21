package v1

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"context"
)

// BrandsSrv 品牌业务服务接口
type BrandsSrv interface {
	Get(ctx context.Context, ID uint64) (*do.BrandsDO, error)

	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error)

	Create(ctx context.Context, brands *do.BrandsDO) error

	Update(ctx context.Context, brands *do.BrandsDO) error

	Delete(ctx context.Context, ID uint64) error
}
type brandService struct {
	//工厂
	data v1.DataFactory
}

func newBrand(srv *serviceFactory) BrandsSrv {
	return &brandService{
		data: srv.data,
	}
}

// Get 根据ID查询品牌
func (b *brandService) Get(ctx context.Context, ID uint64) (*do.BrandsDO, error) {
	brandDO, err := b.data.NewMysql().Brands().Get(ctx, ID)
	// 错误直接向上层抛出
	return brandDO, err
}

// List 分页查询品牌列表
func (b *brandService) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error) {
	brandDOList, err := b.data.NewMysql().Brands().List(ctx, opts, orderby)
	return brandDOList, err
}

// Create 创建品牌
func (b *brandService) Create(ctx context.Context, brands *do.BrandsDO) error {
	err := b.data.NewMysql().Brands().Create(ctx, nil, brands)
	return err
}

// Update 更新品牌
func (b *brandService) Update(ctx context.Context, brands *do.BrandsDO) error {
	err := b.data.NewMysql().Brands().Update(ctx, nil, brands)
	return err
}

// Delete 删除品牌，直接调用data层方法，错误向上抛
func (b *brandService) Delete(ctx context.Context, ID uint64) error {
	err := b.data.NewMysql().Brands().Delete(ctx, ID)
	return err
}

var _ BrandsSrv = &brandService{}
