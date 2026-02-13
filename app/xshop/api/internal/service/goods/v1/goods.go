package v1

import (
	gpb "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/xshop/api/internal/data"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsSrv interface {
	List(ctx context.Context, request *gpb.GoodsFilterRequest) (*gpb.GoodsListResponse, error)
	BatchGetGoods(ctx context.Context, in *gpb.BatchGoodsIdInfo, opts ...grpc.CallOption) (*gpb.GoodsListResponse, error)
	CreateGoods(ctx context.Context, in *gpb.CreateGoodsInfo, opts ...grpc.CallOption) (*gpb.GoodsInfoResponse, error)
	DeleteGoods(ctx context.Context, in *gpb.DeleteGoodsInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateGoods(ctx context.Context, in *gpb.CreateGoodsInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetGoodsDetail(ctx context.Context, in *gpb.GoodInfoRequest, opts ...grpc.CallOption) (*gpb.GoodsInfoResponse, error)
	// 商品分类
	GetAllCategorysList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*gpb.CategoryListResponse, error)
	GetSubCategory(ctx context.Context, in *gpb.CategoryListRequest, opts ...grpc.CallOption) (*gpb.SubCategoryListResponse, error)
	CreateCategory(ctx context.Context, in *gpb.CategoryInfoRequest, opts ...grpc.CallOption) (*gpb.CategoryInfoResponse, error)
	DeleteCategory(ctx context.Context, in *gpb.DeleteCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateCategory(ctx context.Context, in *gpb.CategoryInfoRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 品牌和轮播图（品牌）
	BrandList(ctx context.Context, in *gpb.BrandFilterRequest, opts ...grpc.CallOption) (*gpb.BrandListResponse, error)
	CreateBrand(ctx context.Context, in *gpb.BrandRequest, opts ...grpc.CallOption) (*gpb.BrandInfoResponse, error)
	DeleteBrand(ctx context.Context, in *gpb.BrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateBrand(ctx context.Context, in *gpb.BrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 轮播图
	BannerList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*gpb.BannerListResponse, error)
	CreateBanner(ctx context.Context, in *gpb.BannerRequest, opts ...grpc.CallOption) (*gpb.BannerResponse, error)
	DeleteBanner(ctx context.Context, in *gpb.BannerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateBanner(ctx context.Context, in *gpb.BannerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 品牌分类
	CategoryBrandList(ctx context.Context, in *gpb.CategoryBrandFilterRequest, opts ...grpc.CallOption) (*gpb.CategoryBrandListResponse, error)
	GetCategoryBrandList(ctx context.Context, in *gpb.CategoryInfoRequest, opts ...grpc.CallOption) (*gpb.BrandListResponse, error)
	CreateCategoryBrand(ctx context.Context, in *gpb.CategoryBrandRequest, opts ...grpc.CallOption) (*gpb.CategoryBrandResponse, error)
	DeleteCategoryBrand(ctx context.Context, in *gpb.CategoryBrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateCategoryBrand(ctx context.Context, in *gpb.CategoryBrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type goodsService struct {
	data data.DataFactory
}

// -------------------------- 商品相关方法 --------------------------

func (gs *goodsService) List(ctx context.Context, request *gpb.GoodsFilterRequest) (*gpb.GoodsListResponse, error) {
	return gs.data.Goods().GoodsList(ctx, request)
}

// BatchGetGoods 批量获取商品
func (gs *goodsService) BatchGetGoods(ctx context.Context, in *gpb.BatchGoodsIdInfo, opts ...grpc.CallOption) (*gpb.GoodsListResponse, error) {
	return gs.data.Goods().BatchGetGoods(ctx, in)
}

// CreateGoods 创建商品
func (gs *goodsService) CreateGoods(ctx context.Context, in *gpb.CreateGoodsInfo, opts ...grpc.CallOption) (*gpb.GoodsInfoResponse, error) {
	return gs.data.Goods().CreateGoods(ctx, in)
}

// DeleteGoods 删除商品
func (gs *goodsService) DeleteGoods(ctx context.Context, in *gpb.DeleteGoodsInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().DeleteGoods(ctx, in)
}

// UpdateGoods 更新商品
func (gs *goodsService) UpdateGoods(ctx context.Context, in *gpb.CreateGoodsInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().UpdateGoods(ctx, in)
}

// GetGoodsDetail 获取商品详情
func (gs *goodsService) GetGoodsDetail(ctx context.Context, in *gpb.GoodInfoRequest, opts ...grpc.CallOption) (*gpb.GoodsInfoResponse, error) {
	return gs.data.Goods().GetGoodsDetail(ctx, in)
}

// -------------------------- 商品分类相关方法 --------------------------

// GetAllCategorysList 获取所有分类列表
func (gs *goodsService) GetAllCategorysList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*gpb.CategoryListResponse, error) {
	return gs.data.Goods().GetAllCategorysList(ctx, in)
}

// GetSubCategory 获取子分类
func (gs *goodsService) GetSubCategory(ctx context.Context, in *gpb.CategoryListRequest, opts ...grpc.CallOption) (*gpb.SubCategoryListResponse, error) {
	return gs.data.Goods().GetSubCategory(ctx, in)
}

// CreateCategory 创建分类
func (gs *goodsService) CreateCategory(ctx context.Context, in *gpb.CategoryInfoRequest, opts ...grpc.CallOption) (*gpb.CategoryInfoResponse, error) {
	return gs.data.Goods().CreateCategory(ctx, in)
}

// DeleteCategory 删除分类
func (gs *goodsService) DeleteCategory(ctx context.Context, in *gpb.DeleteCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().DeleteCategory(ctx, in)
}

// UpdateCategory 更新分类
func (gs *goodsService) UpdateCategory(ctx context.Context, in *gpb.CategoryInfoRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().UpdateCategory(ctx, in)
}

// -------------------------- 品牌相关方法 --------------------------

// BrandList 品牌列表
func (gs *goodsService) BrandList(ctx context.Context, in *gpb.BrandFilterRequest, opts ...grpc.CallOption) (*gpb.BrandListResponse, error) {
	return gs.data.Goods().BrandList(ctx, in)
}

// CreateBrand 创建品牌
func (gs *goodsService) CreateBrand(ctx context.Context, in *gpb.BrandRequest, opts ...grpc.CallOption) (*gpb.BrandInfoResponse, error) {
	return gs.data.Goods().CreateBrand(ctx, in)
}

// DeleteBrand 删除品牌
func (gs *goodsService) DeleteBrand(ctx context.Context, in *gpb.BrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().DeleteBrand(ctx, in)
}

// UpdateBrand 更新品牌
func (gs *goodsService) UpdateBrand(ctx context.Context, in *gpb.BrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().UpdateBrand(ctx, in)
}

// -------------------------- 轮播图相关方法 --------------------------
// BannerList 轮播图列表
func (gs *goodsService) BannerList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*gpb.BannerListResponse, error) {
	return gs.data.Goods().BannerList(ctx, in)
}

// CreateBanner 创建轮播图
func (gs *goodsService) CreateBanner(ctx context.Context, in *gpb.BannerRequest, opts ...grpc.CallOption) (*gpb.BannerResponse, error) {
	return gs.data.Goods().CreateBanner(ctx, in)
}

// DeleteBanner 删除轮播图
func (gs *goodsService) DeleteBanner(ctx context.Context, in *gpb.BannerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().DeleteBanner(ctx, in)
}

// UpdateBanner 更新轮播图
func (gs *goodsService) UpdateBanner(ctx context.Context, in *gpb.BannerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().UpdateBanner(ctx, in)
}

// -------------------------- 品牌分类关联相关方法 --------------------------
// CategoryBrandList 品牌分类关联列表
func (gs *goodsService) CategoryBrandList(ctx context.Context, in *gpb.CategoryBrandFilterRequest, opts ...grpc.CallOption) (*gpb.CategoryBrandListResponse, error) {
	return gs.data.Goods().CategoryBrandList(ctx, in)
}

// GetCategoryBrandList 获取指定分类下的品牌列表
func (gs *goodsService) GetCategoryBrandList(ctx context.Context, in *gpb.CategoryInfoRequest, opts ...grpc.CallOption) (*gpb.BrandListResponse, error) {
	return gs.data.Goods().GetCategoryBrandList(ctx, in)
}

// CreateCategoryBrand 创建品牌分类关联
func (gs *goodsService) CreateCategoryBrand(ctx context.Context, in *gpb.CategoryBrandRequest, opts ...grpc.CallOption) (*gpb.CategoryBrandResponse, error) {
	return gs.data.Goods().CreateCategoryBrand(ctx, in)
}

// DeleteCategoryBrand 删除品牌分类关联
func (gs *goodsService) DeleteCategoryBrand(ctx context.Context, in *gpb.CategoryBrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().DeleteCategoryBrand(ctx, in)
}

// UpdateCategoryBrand 更新品牌分类关联
func (gs *goodsService) UpdateCategoryBrand(ctx context.Context, in *gpb.CategoryBrandRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return gs.data.Goods().UpdateCategoryBrand(ctx, in)
}

// NewGoods 创建goodsService实例
func NewGoods(data data.DataFactory) *goodsService {
	return &goodsService{data: data}
}

// 确保goodsService实现了GoodsSrv接口（编译期检查）
var _ GoodsSrv = &goodsService{}
