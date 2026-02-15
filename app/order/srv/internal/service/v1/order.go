package service

import (
	proto3 "Advanced_Shop/api/goods/v1"
	proto2 "Advanced_Shop/api/inventory/v1"
	proto "Advanced_Shop/api/order/v1"
	v12 "Advanced_Shop/app/order/srv/internal/data/v1"
	"Advanced_Shop/app/order/srv/internal/domain/do"
	"Advanced_Shop/app/order/srv/internal/domain/dto"
	code2 "Advanced_Shop/app/pkg/code"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/code"
	v1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"github.com/dtm-labs/client/dtmgrpc"
)

type OrderSrv interface {
	Get(ctx context.Context, orderSn dto.OrderDetailRequest) (*dto.OrderInfoResponse, error)
	List(ctx context.Context, userID uint64, meta v1.ListMeta, orderby []string) (*dto.OrderDTOList, error)
	Submit(ctx context.Context, order *dto.OrderDTO) (float32, error)
	Create(ctx context.Context, order *dto.OrderInfoResponse) error
	CreateCom(ctx context.Context, order *dto.OrderDTO) error //这是create的补偿
	UpdateStatus(ctx context.Context, orderSn string, status string) error
	GetByOrderSn(ctx context.Context, orderSn string) (*dto.OrderInfoResponse, error)
}

type orderService struct {
	data    v12.DataFactory
	dtmOpts *options.DtmOptions
}

func (os *orderService) CreateCom(ctx context.Context, order *dto.OrderDTO) error {
	/*
		1. 删除orderinfo表
		2. 删除ordergoods表
		3. 删除order找到对应的购物车条目，删除购物车条目
	*/
	//其实不用回滚
	//你应该先查询订单是否已经存在，如果已经存在删除相关记录即可， 同时删除购物车记录
	return nil
}

func (os *orderService) Create(ctx context.Context, order *dto.OrderInfoResponse) error {
	/*
		1. 生成orderinfo表
		2. 生成ordergoods表
		3. 根据order找到对应的购物车条目，删除购物车条目
	*/
	txn := os.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			_ = txn.Rollback()
			log.Error("新建订单事务进行中出现异常，回滚")
			return
		}
	}()
	// 所有的创建在这里
	err := os.data.Orders().Create(ctx, txn, order)
	if err != nil {
		txn.Rollback()
		log.Errorf("创建订单失败，err:%v", err)
		return err // 这个不是abort 也就是说会不停的重试
	}

	err = os.data.ShopCarts().DeleteByGoodsIDs(ctx, txn, uint64(order.User), order.GoodIds)
	if err != nil {
		txn.Rollback()
		log.Errorf("删除购物车失败，goodids:%v, err:%v", order.GoodIds, err)
		return err
	}

	txn.Commit()

	return nil
}

func (os *orderService) Get(ctx context.Context, detail dto.OrderDetailRequest) (*dto.OrderInfoResponse, error) {
	orderInfo, err := os.data.Orders().Get(ctx, detail)
	if err != nil {
		return nil, err
	}

	return orderInfo, nil
}
func (os *orderService) GetByOrderSn(ctx context.Context, orderSn string) (*dto.OrderInfoResponse, error) {
	orderInfo, err := os.data.Orders().GetByOrderSn(ctx, orderSn)
	if err != nil {
		return nil, err
	}

	return orderInfo, nil
}

func (os *orderService) List(ctx context.Context, userID uint64, meta v1.ListMeta, orderby []string) (*dto.OrderDTOList, error) {
	orders, err := os.data.Orders().List(ctx, userID, meta, orderby)
	if err != nil {
		return nil, err
	}
	var ret dto.OrderDTOList
	ret.TotalCount = orders.TotalCount
	for _, value := range orders.Items {
		ret.Items = append(ret.Items, &dto.OrderDTO{
			*value,
		})
	}
	return &ret, nil
}

func (os *orderService) Submit(ctx context.Context, order *dto.OrderDTO) (float32, error) {

	//先拿到 选中的 good ID
	response, err := os.data.ShopCarts().GetBatchByUser(ctx, order.User)
	if err != nil {
		log.Errorf("购物车中没有商品，无法下单")
		return 0, err
	}

	goods, err := os.data.Goods().BatchGetGoods(ctx, &proto3.BatchGoodsIdInfo{
		Id: response.GoodsId,
	})
	if err != nil {
		log.Errorf("批量获取商品信息失败，goodids: %v, err:%v", response.GoodsId, err)
		return 0, err
	}

	var PriceSum float32
	// 生成表用的
	var orderGoods []*do.OrderGoodsModel
	// 库存微服务用的
	var goodsInfo []*proto2.GoodsInvInfo
	for _, goodModel := range goods.Data {
		PriceSum += goodModel.ShopPrice * float32(response.GoodNumMap[goodModel.Id])
		orderGoods = append(orderGoods, &do.OrderGoodsModel{
			Goods:      goodModel.Id,
			GoodsName:  goodModel.Name,
			GoodsPrice: goodModel.ShopPrice,
			GoodImages: goodModel.GoodsFrontImage,
			Nums:       response.GoodNumMap[goodModel.Id],
		})
		// 库存服务接收参数
		goodsInfo = append(goodsInfo, &proto2.GoodsInvInfo{
			GoodsId: goodModel.Id,
			Num:     response.GoodNumMap[goodModel.Id],
		})
	}

	// 构建 请求体
	var orderItems []*proto.OrderItemResponse
	for _, good := range orderGoods {
		orderItems = append(orderItems, &proto.OrderItemResponse{
			GoodsId:    good.Goods,
			GoodsName:  good.GoodsName,
			GoodsImage: good.GoodImages,
			GoodsPrice: good.GoodsPrice,
			Nums:       good.Nums,
		})
	}
	// 库存服务
	req := &proto2.SellInfo{
		GoodsInfo: goodsInfo,
		OrderSn:   order.OrderSn,
	}
	// 订单服务
	oReq := &proto.CreateRequest{
		PriceSum:   PriceSum,
		OrderSn:    order.OrderSn,
		UserId:     order.User,
		Address:    order.Address,
		Name:       order.SignerName,
		Mobile:     order.SignerMobile,
		Post:       order.Post,
		OrderItems: orderItems,
		GoodsId:    response.GoodsId, // 用于删除 购物车的
	}

	qsBusi := "discovery:///Shop-inventory-srv"
	gBusi := "discovery:///Shop-order-srv"
	saga := dtmgrpc.NewSagaGrpc(os.dtmOpts.GrpcServer, order.OrderSn).
		Add(qsBusi+"/Inventory/Sell", qsBusi+"/Inventory/Reback", req).
		Add(gBusi+"/Order/CreateOrder", gBusi+"/Order/CreateOrderCom", oReq)
	saga.WaitResult = true
	err = saga.Submit()
	//通过OrderSn查询一下， 当前的状态如何状态一直值Submitted那么就你一直不要给前端返回， 如果是failed那么你提示给前端说下单失败，重新下单
	return PriceSum, err
}

func (os *orderService) UpdateStatus(ctx context.Context, orderSn string, status string) error {
	row, err := os.data.Orders().UpdateStatus(ctx, orderSn, status)
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	if row == 0 {
		return errors.WithCode(code2.ErrOrderStatus, "<UNK>")
	}

	return nil
}

func newOrderService(sv *service) *orderService {
	return &orderService{
		data:    sv.data,
		dtmOpts: sv.dtmopts,
	}
}

var _ OrderSrv = &orderService{}
