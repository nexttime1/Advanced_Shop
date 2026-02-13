package v1

import (
	proto "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/gorm"
	v12 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/log"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CategoryBrandList 分类-品牌关联列表（分页）
func (gs *goodsServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	// 1. 构造分页参数（匹配ListMeta的Page/PageSize字段，而非Limit/Offset）
	listMeta := v12.ListMeta{
		Page:     int(request.Pages),
		PageSize: int(request.PagePerNums), // 修正：原代码用Limit，改为PageSize匹配proto
	}

	// 2. 调用service层获取分类-品牌关联列表
	cbList, err := gs.srv.CategoryBrands().List(ctx, listMeta, []string{"id asc"})
	if err != nil {
		log.Errorf("get category-brand list error: %v", err.Error())
		return nil, err
	}

	// 3. 转换为Proto响应（匹配新的CategoryBrandResponse结构体）
	var ret proto.CategoryBrandListResponse
	ret.Total = int32(cbList.TotalCount)

	for _, cb := range cbList.Items {
		// 初始化响应对象
		cbResp := &proto.CategoryBrandResponse{
			Id: int32(cb.ID), // uint64转int32（匹配proto的int32类型）
		}

		// 获取品牌信息并转换为BrandInfoResponse
		if brand, err := gs.srv.Brands().Get(ctx, uint64(cb.BrandsID)); err == nil {
			cbResp.Brand = &proto.BrandInfoResponse{
				Id:   int32(brand.ID),
				Name: brand.Name,
				Logo: brand.Logo,
			}
		}

		// 获取分类信息并转换为CategoryInfoResponse
		if category, err := gs.srv.Category().Get(ctx, uint64(cb.CategoryID)); err == nil {
			cbResp.Category = &proto.CategoryInfoResponse{
				Id:               category.ID,
				Name:             category.Name,
				ParentCategoryID: category.ParentCategoryID,
				Level:            category.Level,
				IsTab:            category.IsTab,
				SubCategorys:     []*proto.CategoryInfoResponse{}, // 子分类按需填充
			}
		}

		ret.Data = append(ret.Data, cbResp)
	}

	return &ret, nil
}

// GetCategoryBrandList 根据分类ID获取关联品牌列表
func (gs *goodsServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	// 1. 调用service层新增的ListByCategoryID方法
	brandList, err := gs.srv.CategoryBrands().ListByCategoryID(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("get category brand list error, category id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	// 2. 转换为Proto响应（匹配BrandListResponse）
	var ret proto.BrandListResponse
	ret.Total = int32(len(brandList))

	for _, brand := range brandList {
		ret.Data = append(ret.Data, &proto.BrandInfoResponse{
			Id:   int32(brand.ID), // uint64转int32
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	return &ret, nil
}

// CreateCategoryBrand 创建分类-品牌关联
func (gs *goodsServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	// 1. Proto请求转换为DO（注意类型转换：int32转uint64）
	cbDO := &do.GoodsCategoryBrandDO{
		CategoryID: request.CategoryId,
		BrandsID:   request.BrandId,
	}

	// 2. 调用service层创建方法
	err := gs.srv.CategoryBrands().Create(ctx, cbDO)
	if err != nil {
		log.Errorf("create category-brand error, category id: %d, brand id: %d, err: %v",
			request.CategoryId, request.BrandId, err.Error())
		return nil, err
	}

	// 3. 构造Proto响应（匹配新的CategoryBrandResponse结构体）
	cbResp := &proto.CategoryBrandResponse{
		Id: int32(cbDO.ID), // 自增ID转int32
	}

	// 填充品牌信息
	if brand, err := gs.srv.Brands().Get(ctx, uint64(cbDO.BrandsID)); err == nil {
		cbResp.Brand = &proto.BrandInfoResponse{
			Id:   int32(brand.ID),
			Name: brand.Name,
			Logo: brand.Logo,
		}
	}

	// 填充分类信息
	if category, err := gs.srv.Category().Get(ctx, uint64(cbDO.CategoryID)); err == nil {
		cbResp.Category = &proto.CategoryInfoResponse{
			Id:               category.ID,
			Name:             category.Name,
			ParentCategoryID: category.ParentCategoryID,
			Level:            category.Level,
			IsTab:            category.IsTab,
		}
	}

	return cbResp, nil
}

// DeleteCategoryBrand 删除分类-品牌关联
func (gs *goodsServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	// 1. 调用service层删除方法
	err := gs.srv.CategoryBrands().Delete(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("delete category-brand error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	// 2. 返回空响应
	return &emptypb.Empty{}, nil
}

// UpdateCategoryBrand 更新分类-品牌关联
func (gs *goodsServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	// 1. Proto请求转换为DO
	cbDO := &do.GoodsCategoryBrandDO{
		Model:      gorm.Model{ID: request.Id},
		CategoryID: request.CategoryId,
		BrandsID:   request.BrandId,
	}

	// 2. 调用service层更新方法
	err := gs.srv.CategoryBrands().Update(ctx, cbDO)
	if err != nil {
		log.Errorf("update category-brand error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	// 3. 返回空响应
	return &emptypb.Empty{}, nil
}
