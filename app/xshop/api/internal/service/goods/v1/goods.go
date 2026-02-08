package v1

import (
	gpb "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/xshop/api/internal/data"
	"context"
)

type GoodsSrv interface {
	List(ctx context.Context, request *gpb.GoodsFilterRequest) (*gpb.GoodsListResponse, error)
}

type goodsService struct {
	data data.DataFactory
}

func (gs *goodsService) List(ctx context.Context, request *gpb.GoodsFilterRequest) (*gpb.GoodsListResponse, error) {
	return gs.data.Goods().GoodsList(ctx, request)
}

func NewGoods(data data.DataFactory) *goodsService {
	return &goodsService{data: data}
}

var _ GoodsSrv = &goodsService{}
