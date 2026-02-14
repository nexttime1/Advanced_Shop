package order

import (
	pb "Advanced_Shop/api/order/v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (os *orderServer) CartItemList(ctx context.Context, info *pb.UserInfo) (*pb.CartItemListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *orderServer) CreateCartItem(ctx context.Context, request *pb.CartItemRequest) (*pb.ShopCartInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *orderServer) UpdateCartItem(ctx context.Context, request *pb.CartItemRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (os *orderServer) DeleteCartItem(ctx context.Context, request *pb.CartItemRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
