package v1

import (
	invpb "Advanced_Shop/api/inventory/v1"
	"Advanced_Shop/app/inventory/srv/internal/domain/do"
	"Advanced_Shop/app/inventory/srv/internal/domain/dto"
	"Advanced_Shop/app/inventory/srv/internal/service/v1"
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TODO DTM (Distributed Transactions Manager)

type inventoryServer struct {
	invpb.UnimplementedInventoryServer
	srv v1.ServiceFactory
}

// SetInv 设置库存
func (is *inventoryServer) SetInv(ctx context.Context, info *invpb.GoodsInvInfo) (*emptypb.Empty, error) {
	invDTO := &dto.InventoryDTO{}
	invDTO.Goods = info.GoodsId
	invDTO.Stock = info.Num
	err := is.srv.Inventories().Create(ctx, invDTO)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (is *inventoryServer) InvDetail(ctx context.Context, info *invpb.GoodsInvInfo) (*invpb.GoodsInvInfo, error) {
	inv, err := is.srv.Inventories().Get(ctx, uint64(info.GoodsId))
	if err != nil {
		return nil, err
	}
	return &invpb.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stock,
	}, nil
}

func (is *inventoryServer) Sell(ctx context.Context, info *invpb.SellInfo) (*emptypb.Empty, error) {
	var detail []do.GoodsDetail
	for _, value := range info.GoodsInfo {
		detail = append(detail, do.GoodsDetail{GoodId: value.GoodsId, Num: value.Num})
	}
	err := is.srv.Inventories().Sell(ctx, info.OrderSn, detail)
	if err != nil {
		if errors.IsCode(err, code.ErrInvNotEnough) {
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (is *inventoryServer) Reback(ctx context.Context, info *invpb.SellInfo) (*emptypb.Empty, error) {
	log.Infof("订单%s归还库存", info.OrderSn)
	var detail []do.GoodsDetail
	for _, v := range info.GoodsInfo {
		detail = append(detail, do.GoodsDetail{GoodId: v.GoodsId, Num: v.Num})
	}
	err := is.srv.Inventories().Reback(ctx, info.OrderSn, detail)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func NewInventoryServer(srv v1.ServiceFactory) *inventoryServer {
	return &inventoryServer{srv: srv}
}

var (
	_ invpb.InventoryServer = &inventoryServer{}
)
