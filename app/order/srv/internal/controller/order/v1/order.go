package order

import (
	pb "Advanced_Shop/api/order/v1"
	"Advanced_Shop/app/order/srv/internal/domain/do"
	"Advanced_Shop/app/order/srv/internal/domain/dto"
	"Advanced_Shop/app/order/srv/internal/service/v1"
	v1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/log"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderServer struct {
	pb.UnimplementedOrderServer
	srv service.ServiceFactory
}

func NewOrderServer(srv service.ServiceFactory) *orderServer {
	return &orderServer{srv: srv}
}

// CreateOrder 这个是给分布式事务saga调用的，没为api提供的需求
func (os *orderServer) CreateOrder(ctx context.Context, request *pb.CreateRequest) (*emptypb.Empty, error) {
	orderGoods := make([]*do.OrderGoodsModel, len(request.OrderItems))
	for i, item := range request.OrderItems {
		orderGoods[i] = &do.OrderGoodsModel{
			Goods:      item.GoodsId,
			GoodsName:  item.GoodsName,
			GoodImages: item.GoodsImage,
			GoodsPrice: item.GoodsPrice,
			Nums:       item.Nums,
		}
	}

	err := os.srv.Orders().Create(ctx, &dto.OrderInfoResponse{
		OrderInfoDO: do.OrderInfoDO{
			OrderMount:   request.PriceSum,
			User:         request.UserId,
			Address:      request.Address,
			SignerName:   request.Name,
			SignerMobile: request.Mobile,
			Post:         request.Post,
			OrderSn:      request.OrderSn,
		},
		OrderGoods: orderGoods,
		GoodIds:    request.GoodsId,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (os *orderServer) CreateOrderCom(ctx context.Context, request *pb.OrderRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (os *orderServer) SubmitOrder(ctx context.Context, request *pb.OrderRequest) (*pb.SubmitResponse, error) {
	//从购物车中得到选中的商品
	orderDTO := dto.OrderDTO{
		OrderInfoDO: do.OrderInfoDO{
			User:         request.UserId,
			Address:      request.Address,
			SignerName:   request.Name,
			SignerMobile: request.Mobile,
			Post:         request.Post,
			OrderSn:      request.OrderSn,
		},
	}
	total, err := os.srv.Orders().Submit(ctx, &orderDTO)
	if err != nil {
		log.Errorf("新建订单失败: %v", err)
		return nil, err
	}

	//另外一款解决ioc的库，wire

	return &pb.SubmitResponse{PriceSum: total}, nil
}

func (os *orderServer) OrderList(ctx context.Context, request *pb.OrderFilterRequest) (*pb.OrderListResponse, error) {
	response := &pb.OrderListResponse{}
	pageInfo := v1.ListMeta{
		Page:     int(request.Pages),
		PageSize: int(request.PagePerNums),
	}
	list, err := os.srv.Orders().List(ctx, uint64(request.UserId), pageInfo, []string{})
	if err != nil {
		return nil, err
	}
	response.Total = int32(list.TotalCount)
	var modelsInfo []*pb.OrderInfoResponse
	for _, item := range list.Items {
		modelsInfo = append(modelsInfo, &pb.OrderInfoResponse{
			Id:      item.ID,
			UserId:  item.User,
			OrderSn: item.OrderSn,
			PayType: item.PayType,
			Status:  item.Status,
			Post:    item.Post,
			Total:   item.OrderMount,
			Address: item.Address,
			Name:    item.SignerName,
			Mobile:  item.SignerMobile,
		})
	}
	response.Data = modelsInfo

	return response, nil
}

func (os *orderServer) OrderDetail(ctx context.Context, request *pb.OrderRequest) (*pb.OrderInfoDetailResponse, error) {
	var response *pb.OrderInfoDetailResponse
	detail := dto.OrderDetailRequest{
		UserID:  request.UserId,
		OrderID: request.Id,
	}
	resp, err := os.srv.Orders().Get(ctx, detail)
	if err != nil {
		return nil, err
	}
	// 构建返回
	response.OrderInfo = &pb.OrderInfoResponse{
		Id:      resp.ID,
		UserId:  resp.User,
		OrderSn: resp.OrderSn,
		PayType: resp.PayType,
		Status:  resp.Status,
		Post:    resp.Post,
		Total:   resp.OrderMount,
		Address: resp.Address,
		Name:    resp.SignerName,
		Mobile:  resp.SignerMobile,
	}
	var Goods []*pb.OrderItemResponse
	for _, item := range resp.OrderGoods {
		Goods = append(Goods, &pb.OrderItemResponse{
			Id:         item.ID,
			OrderId:    item.Order,
			GoodsId:    item.Goods,
			GoodsName:  item.GoodsName,
			GoodsPrice: item.GoodsPrice,
			Nums:       item.Nums,
		})
	}
	response.Goods = Goods

	return response, nil
}

func (os *orderServer) UpdateOrderStatus(ctx context.Context, status *pb.OrderStatus) (*emptypb.Empty, error) {
	err := os.srv.Orders().UpdateStatus(ctx, status.OrderSn, status.Status)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil

}

func (os *orderServer) OrderDetailByOrderSn(ctx context.Context, request *pb.AlipayOrderSnRequest) (*pb.OrderInfoDetailResponse, error) {
	resp, err := os.srv.Orders().GetByOrderSn(ctx, request.OrderSn)
	if err != nil {
		return nil, err
	}
	var response *pb.OrderInfoDetailResponse
	// 构建返回
	response.OrderInfo = &pb.OrderInfoResponse{
		Id:      resp.ID,
		UserId:  resp.User,
		OrderSn: resp.OrderSn,
		PayType: resp.PayType,
		Status:  resp.Status,
		Post:    resp.Post,
		Total:   resp.OrderMount,
		Address: resp.Address,
		Name:    resp.SignerName,
		Mobile:  resp.SignerMobile,
	}
	var Goods []*pb.OrderItemResponse
	for _, item := range resp.OrderGoods {
		Goods = append(Goods, &pb.OrderItemResponse{
			Id:         item.ID,
			OrderId:    item.Order,
			GoodsId:    item.Goods,
			GoodsName:  item.GoodsName,
			GoodsPrice: item.GoodsPrice,
			Nums:       item.Nums,
		})
	}
	response.Goods = Goods

	return response, nil

}

var _ pb.OrderServer = &orderServer{}
