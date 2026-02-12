package v1

import (
	proto "Advanced_Shop/api/goods/v1"
	good "Advanced_Shop/app/goods/srv/internal/data/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"Advanced_Shop/app/goods/srv/internal/domain/dto"
	v1 "Advanced_Shop/app/goods/srv/internal/service/v1"
	"Advanced_Shop/app/pkg/gorm"
	v12 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/log"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type goodsServer struct {
	proto.UnimplementedGoodsServer
	srv v1.ServiceFactory
}

func GoodInfoFunction(goods *dto.GoodsDTO) *proto.GoodsInfoResponse {
	// 封面图
	firstImage := ""
	// 详情图 type = 2
	var descImages []string
	// 其他图 type = 3
	var otherImages []string
	for _, imageModel := range goods.Images {
		if imageModel.IsMain {
			firstImage = imageModel.ImageURL
		}
		if imageModel.ImageType == do.DetailImageType {
			descImages = append(descImages, imageModel.ImageURL)
		}
		if imageModel.ImageType == do.OtherImageType {
			otherImages = append(otherImages, imageModel.ImageURL)
		}
	}
	var response proto.GoodsInfoResponse
	if goods.ShipFree != nil {
		response.ShipFree = goods.ShipFree
	}
	if goods.IsNew != nil {
		response.IsNew = goods.IsNew
	}
	if goods.IsHot != nil {
		response.IsHot = goods.IsHot
	}
	if goods.OnSale != nil {
		response.OnSale = goods.OnSale
	}
	response.Id = goods.ID
	response.Name = goods.Name
	response.CategoryId = goods.CategoryID
	response.GoodsSn = goods.GoodsSn
	response.ClickNum = goods.ClickNum
	response.SoldNum = goods.SoldNum
	response.FavNum = goods.FavNum
	response.MarketPrice = goods.MarketPrice
	response.ShopPrice = goods.ShopPrice
	response.GoodsBrief = goods.GoodsBrief
	response.GoodsFrontImage = firstImage
	response.DescImages = descImages
	response.Images = otherImages
	response.Brand = &proto.BrandInfoResponse{
		Id:   goods.Brands.ID,
		Name: goods.Brands.Name,
		Logo: goods.Brands.Logo,
	}
	response.Category = &proto.CategoryBriefInfoResponse{
		Id:   goods.Category.ID,
		Name: goods.Category.Name,
	}
	return &response
}

func (gs *goodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	list, err := gs.srv.Goods().List(ctx, v12.ListMeta{int(request.Pages), int(request.PagePerNums)}, request, []string{})
	if err != nil {
		log.Errorf("get goods list error: %v", err.Error())
		return nil, err
	}
	var ret proto.GoodsListResponse
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, GoodInfoFunction(item))
	}
	return &ret, nil
}

func (gs *goodsServer) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var ids []uint64
	for _, id := range info.Id {
		ids = append(ids, uint64(id))
	}
	get, err := gs.srv.Goods().BatchGet(ctx, ids)
	if err != nil {
		return nil, err
	}
	var ret proto.GoodsListResponse
	for _, item := range get {
		ret.Data = append(ret.Data, GoodInfoFunction(item))
	}
	return &ret, nil
}

// CreateGoods 创建商品
func (gs *goodsServer) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {

	goodsDO := &do.GoodsDO{
		Name:        info.Name,
		GoodsSn:     info.GoodsSn,
		CategoryID:  info.CategoryId,
		BrandsID:    info.BrandId,
		MarketPrice: info.MarketPrice,
		ShopPrice:   info.ShopPrice,
		GoodsBrief:  info.GoodsBrief,
		ShipFree:    info.ShipFree,
		IsNew:       info.IsNew,
		IsHot:       info.IsHot,
		OnSale:      info.OnSale,
	}

	request := good.GoodsInfo{
		GoodsDO:         *goodsDO,
		Images:          info.Images,
		DescImages:      info.DescImages,
		GoodsFrontImage: info.GoodsFrontImage,
	}

	// 调用service层创建方法
	err := gs.srv.Goods().CreateInTxn(ctx, &request)
	if err != nil {
		log.Errorf("create goods error: %v", err.Error())
		return nil, err
	}

	// 查询创建后的商品详情
	goodsDTO, err := gs.srv.Goods().Get(ctx, uint64(goodsDO.ID))
	if err != nil {
		log.Errorf("get created goods detail error: %v", err.Error())
		return nil, err
	}

	return GoodInfoFunction(goodsDTO), nil
}

// DeleteGoods 删除商品
func (gs *goodsServer) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	// 调用service层删除方法
	err := gs.srv.Goods().Delete(ctx, uint64(info.Id))
	if err != nil {
		log.Errorf("delete goods error, id: %d, err: %v", info.Id, err.Error())
		return nil, err
	}

	// 返回空响应
	return &emptypb.Empty{}, nil
}

// UpdateGoods 更新商品
func (gs *goodsServer) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*emptypb.Empty, error) {

	goodsDO := &do.GoodsDO{
		Model:       gorm.Model{ID: info.Id},
		Name:        info.Name,
		GoodsSn:     info.GoodsSn,
		CategoryID:  info.CategoryId,
		BrandsID:    info.BrandId,
		MarketPrice: info.MarketPrice,
		ShopPrice:   info.ShopPrice,
		GoodsBrief:  info.GoodsBrief,
		ShipFree:    info.ShipFree,
		IsNew:       info.IsNew,
		IsHot:       info.IsHot,
		OnSale:      info.OnSale,
	}

	request := good.GoodsInfo{
		GoodsDO:         *goodsDO,
		Images:          info.Images,
		DescImages:      info.DescImages,
		GoodsFrontImage: info.GoodsFrontImage,
	}

	// 调用service层更新方法
	err := gs.srv.Goods().UpdateInTxn(ctx, &request)
	if err != nil {
		log.Errorf("update goods error, id: %d, err: %v", info.Id, err.Error())
		return nil, err
	}
	// 3. 返回空响应
	return &emptypb.Empty{}, nil
}

// GetGoodsDetail 获取商品详情
func (gs *goodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	// 调用service层获取商品详情
	goodsDTO, err := gs.srv.Goods().Get(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("get goods detail error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	return GoodInfoFunction(goodsDTO), nil
}

func (gs *goodsServer) BannerList(ctx context.Context, empty *emptypb.Empty) (*proto.BannerListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *goodsServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func NewGoodsServer(srv v1.ServiceFactory) *goodsServer {
	return &goodsServer{srv: srv}
}
