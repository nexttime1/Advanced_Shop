package v1

import (
	pb "Advanced_Shop/api/order/v1"
	"Advanced_Shop/app/xshop/api/internal/data"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderSrv interface {
	CartItemList(context.Context, *pb.UserInfo) (*pb.CartItemListResponse, error)
	CreateCartItem(context.Context, *pb.CartItemRequest) (*pb.ShopCartInfoResponse, error)
	UpdateCartItem(context.Context, *pb.CartItemRequest) (*emptypb.Empty, error)
	DeleteCartItem(context.Context, *pb.CartItemRequest) (*emptypb.Empty, error)
	// 订单
	CreateOrder(context.Context, *pb.CreateRequest) (*emptypb.Empty, error)
	CreateOrderCom(context.Context, *pb.OrderRequest) (*emptypb.Empty, error)
	SubmitOrder(context.Context, *pb.OrderRequest) (*emptypb.Empty, error)
	OrderList(context.Context, *pb.OrderFilterRequest) (*pb.OrderListResponse, error)
	OrderDetail(context.Context, *pb.OrderRequest) (*pb.OrderInfoDetailResponse, error)
	UpdateOrderStatus(context.Context, *pb.OrderStatus) (*emptypb.Empty, error)
}

type orderService struct {
	data data.DataFactory
}

var _ OrderSrv = (*orderService)(nil)

func NewOrderService(data data.DataFactory) OrderSrv {
	return &orderService{
		data: data,
	}
}

func (o orderService) CartItemList(ctx context.Context, info *pb.UserInfo) (*pb.CartItemListResponse, error) {
	return o.data.Order().CartItemList(ctx, info)
}

func (o orderService) CreateCartItem(ctx context.Context, request *pb.CartItemRequest) (*pb.ShopCartInfoResponse, error) {
	return o.data.Order().CreateCartItem(ctx, request)
}

func (o orderService) UpdateCartItem(ctx context.Context, request *pb.CartItemRequest) (*emptypb.Empty, error) {
	return o.data.Order().UpdateCartItem(ctx, request)
}

func (o orderService) DeleteCartItem(ctx context.Context, request *pb.CartItemRequest) (*emptypb.Empty, error) {
	return o.data.Order().DeleteCartItem(ctx, request)
}

func (o orderService) CreateOrder(ctx context.Context, request *pb.CreateRequest) (*emptypb.Empty, error) {
	return o.data.Order().CreateOrder(ctx, request)
}

func (o orderService) CreateOrderCom(ctx context.Context, request *pb.OrderRequest) (*emptypb.Empty, error) {
	return o.data.Order().CreateOrderCom(ctx, request)
}

func (o orderService) SubmitOrder(ctx context.Context, request *pb.OrderRequest) (*emptypb.Empty, error) {
	return o.data.Order().SubmitOrder(ctx, request)
}

func (o orderService) OrderList(ctx context.Context, request *pb.OrderFilterRequest) (*pb.OrderListResponse, error) {
	return o.data.Order().OrderList(ctx, request)
}

func (o orderService) OrderDetail(ctx context.Context, request *pb.OrderRequest) (*pb.OrderInfoDetailResponse, error) {
	return o.data.Order().OrderDetail(ctx, request)
}

func (o orderService) UpdateOrderStatus(ctx context.Context, status *pb.OrderStatus) (*emptypb.Empty, error) {
	return o.data.Order().UpdateOrderStatus(ctx, status)
}
