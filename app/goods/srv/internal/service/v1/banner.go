package v1

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"context"
)

// BannerSrv 轮播图业务服务接口
type BannerSrv interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BannerList, error)

	Create(ctx context.Context, banner *do.BannerDO) error

	Update(ctx context.Context, banner *do.BannerDO) error

	Delete(ctx context.Context, ID uint64) error
}

// bannerService 轮播图业务服务具体实现
type bannerService struct {
	data v1.DataFactory
}

// newBanner 初始化轮播图业务服务实例
func newBanner(srv *serviceFactory) BannerSrv {
	return &bannerService{
		data: srv.data,
	}
}

// List 分页查询轮播图列表
func (b *bannerService) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BannerList, error) {

	bannerList, err := b.data.NewMysql().Banners().List(ctx, opts, orderby)
	return bannerList, err
}

// Create 创建轮播图
func (b *bannerService) Create(ctx context.Context, banner *do.BannerDO) error {
	err := b.data.NewMysql().Banners().Create(ctx, nil, banner)
	return err
}

// Update 更新轮播图
func (b *bannerService) Update(ctx context.Context, banner *do.BannerDO) error {

	err := b.data.NewMysql().Banners().Update(ctx, nil, banner)
	return err
}

// Delete 删除轮播图
func (b *bannerService) Delete(ctx context.Context, ID uint64) error {

	err := b.data.NewMysql().Banners().Delete(ctx, ID)
	return err
}

var _ BannerSrv = &bannerService{}
