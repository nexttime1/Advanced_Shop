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

// BannerList 轮播图列表
func (gs *goodsServer) BannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
	listMeta := v12.ListMeta{}
	bannerList, err := gs.srv.Banner().List(ctx, listMeta, []string{"sort asc"})
	if err != nil {
		log.Errorf("get banner list error: %v", err.Error())
		return nil, err
	}
	var ret proto.BannerListResponse
	for _, banner := range bannerList.Items {
		ret.Data = append(ret.Data, &proto.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Url:   banner.Url,
			Index: banner.Index,
		})
	}

	return &ret, nil
}

// CreateBanner 创建轮播图
func (gs *goodsServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {

	bannerDO := &do.BannerDO{
		Image: request.Image,
		Url:   request.Url,
		Index: request.Index,
	}

	err := gs.srv.Banner().Create(ctx, bannerDO)
	if err != nil {
		log.Errorf("create banner error: %v", err.Error())
		return nil, err
	}

	return &proto.BannerResponse{
		Id:    bannerDO.ID,
		Image: bannerDO.Image,
		Url:   bannerDO.Url,
		Index: bannerDO.Index,
	}, nil
}

// DeleteBanner 删除轮播图
func (gs *goodsServer) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {

	err := gs.srv.Banner().Delete(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("delete banner error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	// 2. 返回空响应
	return &emptypb.Empty{}, nil
}

// UpdateBanner 更新轮播图
func (gs *goodsServer) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {

	bannerDO := &do.BannerDO{
		Model: gorm.Model{ID: request.Id},
		Image: request.Image,
		Url:   request.Url,
		Index: request.Index,
	}

	err := gs.srv.Banner().Update(ctx, bannerDO)
	if err != nil {
		log.Errorf("update banner error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
