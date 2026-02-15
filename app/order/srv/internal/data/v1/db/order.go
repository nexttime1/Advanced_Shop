package db

import (
	v1 "Advanced_Shop/app/order/srv/internal/data/v1"
	"Advanced_Shop/app/order/srv/internal/domain/do"
	"Advanced_Shop/app/order/srv/internal/domain/dto"
	"Advanced_Shop/app/pkg/code"
	gorm2 "Advanced_Shop/app/pkg/gorm"
	code2 "Advanced_Shop/gnova/code"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"gorm.io/gorm"
)

type orders struct {
	db *gorm.DB
}

func newOrders(factory *dataFactory) *orders {
	return &orders{
		db: factory.db,
	}
}

func (o *orders) Get(ctx context.Context, detail dto.OrderDetailRequest) (*dto.OrderInfoResponse, error) {

	var model do.OrderInfoDO
	err := o.db.Where(do.OrderInfoDO{User: detail.UserID, Model: gorm2.Model{ID: detail.OrderID}}).Take(&model).Error
	if err != nil {
		log.Errorf("get order info error: %v", err)
		return nil, errors.WithCode(code.ErrOrderNotFound, "get order info error")
	}
	response := dto.OrderInfoResponse{
		OrderInfoDO: do.OrderInfoDO{
			Model:        gorm2.Model{ID: model.ID},
			User:         model.User,
			OrderSn:      model.OrderSn,
			PayType:      model.PayType,
			Status:       model.Status,
			TradeNo:      model.TradeNo,
			OrderMount:   model.OrderMount,
			PayTime:      model.PayTime,
			Address:      model.Address,
			SignerName:   model.SignerName,
			SignerMobile: model.SignerMobile,
			Post:         model.Post,
		},
	}
	// 找一下商品
	var goodModels []*do.OrderGoodsModel
	o.db.Where("`order` = ?", model.ID).Find(&goodModels)
	response.OrderGoods = goodModels

	return &response, nil
}
func (o *orders) GetByOrderSn(ctx context.Context, orderSn string) (*dto.OrderInfoResponse, error) {

	var model do.OrderInfoDO
	err := o.db.Where("order_sn = ?", orderSn).Take(&model).Error
	if err != nil {
		log.Errorf("get order info error: %v", err)
		return nil, errors.WithCode(code.ErrOrderNotFound, "get order info error")
	}
	response := dto.OrderInfoResponse{
		OrderInfoDO: do.OrderInfoDO{
			Model:        gorm2.Model{ID: model.ID},
			User:         model.User,
			OrderSn:      model.OrderSn,
			PayType:      model.PayType,
			Status:       model.Status,
			TradeNo:      model.TradeNo,
			OrderMount:   model.OrderMount,
			PayTime:      model.PayTime,
			Address:      model.Address,
			SignerName:   model.SignerName,
			SignerMobile: model.SignerMobile,
			Post:         model.Post,
		},
	}
	// 找一下商品
	var goodModels []*do.OrderGoodsModel
	o.db.Where("`order` = ?", model.ID).Find(&goodModels)
	response.OrderGoods = goodModels

	return &response, nil
}

func (o *orders) List(ctx context.Context, userID uint64, meta metav1.ListMeta, orderby []string) (*do.OrderInfoDOList, error) {
	ret := &do.OrderInfoDOList{}
	//分页
	limit := meta.GetLimit()
	offset := meta.GetOffset()
	//排序
	query := o.db.Model(do.OrderInfoDO{User: int32(userID)})
	for _, value := range orderby {
		query = query.Order(value)
	}

	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

// Create 创建订单之后要删除对应的购物车记录
func (o *orders) Create(ctx context.Context, txn *gorm.DB, order *dto.OrderInfoResponse) error {
	db := o.db
	if txn != nil {
		db = txn
	}
	orderModel := &do.OrderInfoDO{
		User:         order.User,
		OrderSn:      order.OrderSn,
		OrderMount:   order.OrderMount,
		Address:      order.Address,
		SignerName:   order.SignerName,
		SignerMobile: order.SignerMobile,
		Post:         order.Post,
	}

	err := db.Create(&orderModel).Error
	if err != nil {
		log.Errorf("create order failed, error: %v", err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}

	var orderGoodsModels []do.OrderGoodsModel
	for _, goodsDTO := range order.OrderGoods {
		orderGoodsModels = append(orderGoodsModels, do.OrderGoodsModel{
			Order:      orderModel.ID,
			Goods:      goodsDTO.Goods,
			GoodsName:  goodsDTO.GoodsName,
			GoodsPrice: goodsDTO.GoodsPrice,
			GoodImages: goodsDTO.GoodImages,
			Nums:       goodsDTO.Nums,
		})
	}

	// 生成 OrderGoodsModel 表数据
	err = db.CreateInBatches(&orderGoodsModels, 100).Error
	if err != nil {
		log.Errorf("create order failed, error: %v", err)
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}

	return nil
}

func (o *orders) UpdateStatus(ctx context.Context, orderSn string, status string) (int64, error) {
	result := o.db.Model(&do.OrderInfoDO{}).
		Where("order_sn = ?", orderSn).
		Update("status", status)

	return result.RowsAffected, result.Error
}

var _ v1.OrderStore = &orders{}
