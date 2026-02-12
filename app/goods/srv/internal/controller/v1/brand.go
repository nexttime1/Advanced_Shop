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

// BrandList 品牌列表（分页）
func (gs *goodsServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	// 1. 构造分页参数
	listMeta := v12.ListMeta{
		Page:     int(request.Pages),
		PageSize: int(request.PagePerNums),
	}

	// 2. 调用service层获取品牌列表
	brandList, err := gs.srv.Brands().List(ctx, listMeta, []string{"id asc"})
	if err != nil {
		log.Errorf("get brand list error: %v", err.Error())
		return nil, err
	}

	var ret proto.BrandListResponse
	ret.Total = int32(brandList.TotalCount)
	for _, brand := range brandList.Items {
		ret.Data = append(ret.Data, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	return &ret, nil
}

// CreateBrand 创建品牌
func (gs *goodsServer) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	// 1. Proto请求转换为DO
	brandDO := &do.BrandsDO{
		Name: request.Name,
		Logo: request.Logo,
	}

	// 2. 调用service层创建方法
	err := gs.srv.Brands().Create(ctx, brandDO)
	if err != nil {
		log.Errorf("create brand error: %v", err.Error())
		return nil, err
	}

	// 3. 转换为Proto响应
	return &proto.BrandInfoResponse{
		Id:   brandDO.ID,
		Name: brandDO.Name,
		Logo: brandDO.Logo,
	}, nil
}

// DeleteBrand 删除品牌
func (gs *goodsServer) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	// 1. 调用service层删除方法
	err := gs.srv.Brands().Delete(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("delete brand error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	// 2. 返回空响应
	return &emptypb.Empty{}, nil
}

// UpdateBrand 更新品牌
func (gs *goodsServer) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*emptypb.Empty, error) {
	// 1. Proto请求转换为DO
	brandDO := &do.BrandsDO{
		Model: gorm.Model{ID: request.Id},
		Name:  request.Name,
		Logo:  request.Logo,
	}

	// 2. 调用service层更新方法
	err := gs.srv.Brands().Update(ctx, brandDO)
	if err != nil {
		log.Errorf("update brand error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	// 3. 返回空响应
	return &emptypb.Empty{}, nil
}
