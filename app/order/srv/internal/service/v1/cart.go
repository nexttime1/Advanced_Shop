package service

import (
	proto "Advanced_Shop/api/goods/v1"
	proto1 "Advanced_Shop/api/inventory/v1"
	v12 "Advanced_Shop/app/order/srv/internal/data/v1"
	"Advanced_Shop/app/order/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"

	"gorm.io/gorm"
)

// CartSrv 购物车服务接口，对应数据层的所有核心操作
type CartSrv interface {
	// DeleteByGoodsIDs 批量删除指定用户的指定商品ID购物车记录（支持事务）
	DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID uint64, goodsIDs []int32) error
	// List 分页查询用户购物车列表
	List(ctx context.Context, userID uint64, checked bool, meta v1.ListMeta, orderby []string) (*do.ShoppingCartDOList, error)
	// Create 添加商品到购物车
	Create(ctx context.Context, cartItem *do.ShoppingCartDO) (int32, error)
	// Get 根据用户ID和商品ID查询购物车单品
	Get(ctx context.Context, userID, goodsID uint64) (*do.ShoppingCartDO, error)
	// UpdateNum 更新购物车商品数量和选中状态
	UpdateNum(ctx context.Context, cartItem *do.ShoppingCartDO) error
	// Delete 根据购物车ID删除单条记录
	Delete(ctx context.Context, userID uint64, goodID uint64) error
	// GetBatchByUser 获取用户选中的购物车商品批量信息
	GetBatchByUser(ctx context.Context, userID int32) (*do.GetShoppingBatchResponse, error)
	// ClearCheck 清空用户购物
	ClearCheck(ctx context.Context, userID uint64) error
}

// cartService 购物车服务实现结构体，依赖数据层工厂
type cartService struct {
	data v12.DataFactory // 数据层工厂
}

// NewCartService 创建购物车服务实例
func NewCartService(s *service) CartSrv {
	return &cartService{
		data: s.data,
	}
}

// DeleteByGoodsIDs 实现CartSrv接口的批量删除方法
func (cs *cartService) DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID uint64, goodsIDs []int32) error {
	if userID == 0 || len(goodsIDs) == 0 {
		return errors.WithCode(code.ErrInvalidParameter, "用户ID或商品ID列表不能为空")
	}

	err := cs.data.ShopCarts().DeleteByGoodsIDs(ctx, txn, userID, goodsIDs)
	if err != nil {
		log.Errorf("CartSrv DeleteByGoodsIDs failed: userID=%d, goodsIDs=%v, err=%v", userID, goodsIDs, err)
		return err
	}
	return nil
}

// List 实现CartSrv接口的列表查询方法
func (cs *cartService) List(ctx context.Context, userID uint64, checked bool, meta v1.ListMeta, orderby []string) (*do.ShoppingCartDOList, error) {
	if userID == 0 {
		return nil, errors.WithCode(code.ErrInvalidParameter, "用户ID不能为空")
	}

	// 调用数据层查询DO列表
	doList, err := cs.data.ShopCarts().List(ctx, userID, checked, meta, orderby)
	if err != nil {
		log.Errorf("CartSrv List failed: userID=%d, checked=%v, err=%v", userID, checked, err)
		return nil, err
	}

	return doList, nil
}

// Create 实现CartSrv接口的创建购物车方法
func (cs *cartService) Create(ctx context.Context, cartItem *do.ShoppingCartDO) (int32, error) {
	if cartItem == nil || cartItem.User == 0 || cartItem.Goods == 0 {
		return 0, errors.WithCode(code.ErrInvalidParameter, "购物车数据不能为空，用户ID和商品ID必须指定")
	}

	_, err := cs.data.Goods().GetGoodsDetail(ctx, &proto.GoodInfoRequest{Id: cartItem.Goods})
	if err != nil {
		return 0, err
	}

	// 检查库存
	detail, err := cs.data.Inventorys().InvDetail(ctx, &proto1.GoodsInvInfo{
		GoodsId: cartItem.Goods,
	})
	if err != nil {
		return 0, err
	}
	if cartItem.Nums > detail.Num {
		return 0, errors.WithCode(code.ErrInvNotEnough, "库存不足")
	}

	id, err := cs.data.ShopCarts().Create(ctx, cartItem)
	if err != nil {
		log.Errorf("CartSrv Create failed: userID=%d, goodsID=%d, err=%v", cartItem.User, cartItem.Goods, err)
		return 0, err
	}
	return id, nil
}

// Get 实现CartSrv接口的查询单条购物车方法
func (cs *cartService) Get(ctx context.Context, userID, goodsID uint64) (*do.ShoppingCartDO, error) {
	if userID == 0 || goodsID == 0 {
		return nil, errors.WithCode(code.ErrInvalidParameter, "用户ID和商品ID不能为空")
	}

	doItem, err := cs.data.ShopCarts().Get(ctx, userID, goodsID)
	if err != nil {
		log.Errorf("CartSrv Get failed: userID=%d, goodsID=%d, err=%v", userID, goodsID, err)
		return nil, err
	}
	return doItem, nil
}

// UpdateNum 实现CartSrv接口的更新数量和选中状态方法
func (cs *cartService) UpdateNum(ctx context.Context, cartItem *do.ShoppingCartDO) error {
	if cartItem == nil || cartItem.User == 0 || cartItem.Goods == 0 {
		return errors.WithCode(code.ErrInvalidParameter, "购物车数据不能为空，用户ID和商品ID必须指定")
	}

	err := cs.data.ShopCarts().UpdateNum(ctx, cartItem)
	if err != nil {
		log.Errorf("CartSrv UpdateNum failed: userID=%d, goodsID=%d, err=%v", cartItem.User, cartItem.Goods, err)
		return err
	}
	return nil
}

// Delete 实现CartSrv接口的删除单条购物车方法
func (cs *cartService) Delete(ctx context.Context, userID uint64, goodID uint64) error {
	if userID == 0 {
		return errors.WithCode(code.ErrInvalidParameter, "usrID不能为空")
	}

	err := cs.data.ShopCarts().Delete(ctx, userID, goodID)
	if err != nil {
		log.Errorf("CartSrv Delete failed: , err=%v", err)
		return err
	}
	return nil
}

// GetBatchByUser 实现CartSrv接口的批量获取选中商品方法
func (cs *cartService) GetBatchByUser(ctx context.Context, userID int32) (*do.GetShoppingBatchResponse, error) {
	if userID == 0 {
		return nil, errors.WithCode(code.ErrInvalidParameter, "用户ID不能为空")
	}

	doResp, err := cs.data.ShopCarts().GetBatchByUser(ctx, userID)
	if err != nil {
		log.Errorf("CartSrv GetBatchByUser failed: userID=%d, err=%v", userID, err)
		return nil, err
	}

	return doResp, nil
}

// ClearCheck 实现CartSrv接口的清空购物车方法
func (cs *cartService) ClearCheck(ctx context.Context, userID uint64) error {
	if userID == 0 {
		return errors.WithCode(code.ErrInvalidParameter, "用户ID不能为空")
	}

	err := cs.data.ShopCarts().ClearCheck(ctx, userID)
	if err != nil {
		log.Errorf("CartSrv ClearCheck failed: userID=%d, err=%v", userID, err)
		return err
	}
	return nil
}
