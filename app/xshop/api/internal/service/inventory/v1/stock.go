package v1

import (
	pb "Advanced_Shop/api/inventory/v1"
	"Advanced_Shop/app/xshop/api/internal/data"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventorySrv interface {
	SetInv(ctx context.Context, in *pb.GoodsInvInfo) (*emptypb.Empty, error)
	InvDetail(ctx context.Context, in *pb.GoodsInvInfo) (*pb.GoodsInvInfo, error)
	Sell(ctx context.Context, in *pb.SellInfo) (*emptypb.Empty, error)
	Reback(ctx context.Context, in *pb.SellInfo) (*emptypb.Empty, error)
}

type inventorySrv struct {
	data data.DataFactory
}

func NewInventoryService(data data.DataFactory) InventorySrv {
	return &inventorySrv{
		data: data,
	}
}

func (is *inventorySrv) SetInv(ctx context.Context, in *pb.GoodsInvInfo) (*emptypb.Empty, error) {
	return is.data.Inventory().SetInv(ctx, in)
}

func (is *inventorySrv) InvDetail(ctx context.Context, in *pb.GoodsInvInfo) (*pb.GoodsInvInfo, error) {
	return is.data.Inventory().InvDetail(ctx, in)
}

func (is *inventorySrv) Sell(ctx context.Context, in *pb.SellInfo) (*emptypb.Empty, error) {
	return is.data.Inventory().Sell(ctx, in)
}

func (is *inventorySrv) Reback(ctx context.Context, in *pb.SellInfo) (*emptypb.Empty, error) {
	return is.data.Inventory().Reback(ctx, in)
}

var _ InventorySrv = &inventorySrv{}
