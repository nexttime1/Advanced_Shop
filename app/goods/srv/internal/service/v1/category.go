package v1

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"context"
)

// CategorySrv 分类业务服务接口
type CategorySrv interface {
	// Get 根据ID查询单个分类（包含子分类嵌套）
	Get(ctx context.Context, ID uint64) (*do.CategoryDO, error)

	// ListAll 查询所有一级分类（包含子分类嵌套），支持排序
	ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error)

	// Create 创建分类
	Create(ctx context.Context, category *do.CategoryDO) error

	// Update 更新分类
	Update(ctx context.Context, category *do.CategoryDO) error

	// Delete 删除分类
	Delete(ctx context.Context, ID uint64) error
}

// categoryService 分类业务服务具体实现
type categoryService struct {
	data v1.DataFactory
}

// newCategory 初始化分类业务服务实例
func newCategory(srv *serviceFactory) CategorySrv {
	return &categoryService{
		data: srv.data,
	}
}

// Get 根据ID查询分类
func (c *categoryService) Get(ctx context.Context, ID uint64) (*do.CategoryDO, error) {
	categoryDO, err := c.data.NewMysql().Categorys().Get(ctx, ID)
	return categoryDO, err
}

// ListAll 查询所有一级分类
func (c *categoryService) ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error) {
	categoryDOList, err := c.data.NewMysql().Categorys().ListAll(ctx, orderby)
	return categoryDOList, err
}

// Create 创建分类
func (c *categoryService) Create(ctx context.Context, category *do.CategoryDO) error {
	err := c.data.NewMysql().Categorys().Create(ctx, category)
	return err
}

// Update 更新分类
func (c *categoryService) Update(ctx context.Context, category *do.CategoryDO) error {

	err := c.data.NewMysql().Categorys().Update(ctx, category)
	return err
}

// Delete 删除分类
func (c *categoryService) Delete(ctx context.Context, ID uint64) error {
	err := c.data.NewMysql().Categorys().Delete(ctx, ID)
	return err
}

// 确保categoryService完全实现CategoryService接口
var _ CategorySrv = &categoryService{}
