package v1

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"context"
)

// CategoryBrandSrv 分类-品牌关联业务服务接口
type CategoryBrandSrv interface {
	// List 分页查询分类-品牌关联列表，支持排序
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandList, error)

	// Create 创建分类-品牌关联关系
	// 备注：service层隐藏事务参数，默认传nil（如需事务可在service层扩展）
	Create(ctx context.Context, gcb *do.GoodsCategoryBrandDO) error

	// Update 更新分类-品牌关联关系
	Update(ctx context.Context, gcb *do.GoodsCategoryBrandDO) error

	// Delete 删除分类-品牌关联关系
	Delete(ctx context.Context, ID uint64) error
}

// categoryBrandService 分类-品牌关联业务服务具体实现
type categoryBrandService struct {
	data v1.DataFactory
}

// newCategoryBrand 初始化分类-品牌关联业务服务实例
func newCategoryBrand(srv *serviceFactory) CategoryBrandSrv {
	return &categoryBrandService{
		data: srv.data,
	}
}

// List 分页查询分类-品牌关联列表
func (cb *categoryBrandService) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandList, error) {

	gcbList, err := cb.data.CategoryBrands().List(ctx, opts, orderby)
	return gcbList, err
}

// Create 创建分类-品牌关联
func (cb *categoryBrandService) Create(ctx context.Context, gcb *do.GoodsCategoryBrandDO) error {
	err := cb.data.CategoryBrands().Create(ctx, nil, gcb)
	return err
}

// Update 更新分类-品牌关联
func (cb *categoryBrandService) Update(ctx context.Context, gcb *do.GoodsCategoryBrandDO) error {
	// data层Update需要txn参数，无事务场景传nil
	err := cb.data.CategoryBrands().Update(ctx, nil, gcb)
	return err
}

// Delete 删除分类-品牌关联
func (cb *categoryBrandService) Delete(ctx context.Context, ID uint64) error {

	err := cb.data.CategoryBrands().Delete(ctx, ID)
	return err
}

var _ CategoryBrandService = &categoryBrandService{}
