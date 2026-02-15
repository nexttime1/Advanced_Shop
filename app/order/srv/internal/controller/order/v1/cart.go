package order

import (
	pb "Advanced_Shop/api/order/v1"
	"Advanced_Shop/app/order/srv/internal/domain/do"
	v1 "Advanced_Shop/pkg/common/meta/v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (os *orderServer) CartItemList(ctx context.Context, info *pb.UserInfo) (*pb.CartItemListResponse, error) {

	list, err := os.srv.Cart().List(ctx, uint64(info.Id), false, v1.ListMeta{}, []string{})
	if err != nil {
		return nil, err
	}
	var response pb.CartItemListResponse
	var modelsInfo []*pb.ShopCartInfoResponse
	for _, model := range list.Items {
		modelsInfo = append(modelsInfo, &pb.ShopCartInfoResponse{
			Id:      model.ID,
			UserId:  model.User,
			GoodsId: model.Goods,
			Nums:    model.Nums,
			Checked: *model.Checked,
		})

	}
	response.Total = int32(list.TotalCount)
	response.Data = modelsInfo
	return &response, nil

}

func (os *orderServer) CreateCartItem(ctx context.Context, request *pb.CartItemRequest) (*pb.ShopCartInfoResponse, error) {
	info := do.ShoppingCartDO{
		User:    request.UserId,
		Goods:   request.GoodsId,
		Nums:    request.Nums,
		Checked: request.Checked,
	}
	id, err := os.srv.Cart().Create(ctx, &info)
	if err != nil {
		return nil, err
	}

	return &pb.ShopCartInfoResponse{
		Id: id,
	}, nil

}

func (os *orderServer) UpdateCartItem(ctx context.Context, request *pb.CartItemRequest) (*emptypb.Empty, error) {
	info := do.ShoppingCartDO{
		Goods:   request.GoodsId,
		Nums:    request.Nums,
		Checked: request.Checked,
		User:    request.UserId,
	}
	err := os.srv.Cart().UpdateNum(ctx, &info)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (os *orderServer) DeleteCartItem(ctx context.Context, request *pb.CartItemRequest) (*emptypb.Empty, error) {
	err := os.srv.Cart().Delete(ctx, uint64(request.UserId), uint64(request.GoodsId))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
