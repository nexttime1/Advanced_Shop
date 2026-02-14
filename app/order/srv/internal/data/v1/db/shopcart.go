package db

import (
	"Advanced_Shop/app/pkg/code"
	code2 "Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"

	v1 "Advanced_Shop/app/order/srv/internal/data/v1"
	"Advanced_Shop/app/order/srv/internal/domain/do"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"gorm.io/gorm"
)

type shopCarts struct {
	db *gorm.DB
}

func newShopCarts(factory *dataFactory) *shopCarts {
	return &shopCarts{
		db: factory.db,
	}
}

// DeleteByGoodsIDs 这个在事务中执行，以后可以使用消息队列来实现
func (sc *shopCarts) DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID uint64, goodsIDs []int32) error {
	db := sc.db
	if txn != nil {
		db = txn
	}
	return db.Where("user = ? AND goods IN (?)", userID, goodsIDs).Delete(&do.ShoppingCartDO{}).Error
}

func (sc *shopCarts) List(ctx context.Context, userID uint64, checked bool, meta metav1.ListMeta, orderby []string) (*do.ShoppingCartDOList, error) {
	ret := &do.ShoppingCartDOList{}
	query := sc.db

	if userID > 0 {
		query = query.Where("user = ?", userID)
	}
	if checked {
		query = query.Where("checked = ?", true)
	}
	//分页
	limit := meta.GetLimit()
	offset := meta.GetOffset()

	//排序
	for _, value := range orderby {
		query = query.Order(value)
	}

	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

func (sc *shopCarts) Create(ctx context.Context, cartItem *do.ShoppingCartDO) error {
	tx := sc.db.Create(cartItem)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (sc *shopCarts) Get(ctx context.Context, userID, goodsID uint64) (*do.ShoppingCartDO, error) {
	var shopCart do.ShoppingCartDO
	err := sc.db.WithContext(ctx).Where("user = ? AND goods = ?", userID, goodsID).First(&shopCart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrShopCartItemNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &shopCart, nil
}

func (sc *shopCarts) UpdateNum(ctx context.Context, cartItem *do.ShoppingCartDO) error {
	return sc.db.Model(&do.ShoppingCartDO{}).Where("user = ? AND goods = ?", cartItem.User, cartItem.Goods).Update("nums", cartItem.Nums).Update("checked", cartItem.Checked).Error
}

func (sc *shopCarts) Delete(ctx context.Context, ID uint64) error {
	err := sc.db.Where("id = ?", ID).Delete(&do.ShoppingCartDO{}).Error
	if err != nil {
		log.Errorf(err.Error())
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

func (sc *shopCarts) GetBatchByUser(ctx context.Context, userID int32) (*do.GetShoppingBatchResponse, error) {
	check := true
	var shopModels []do.ShoppingCartDO
	sc.db.Where(do.ShoppingCartDO{
		User:    userID,
		Checked: &check,
	}).Find(&shopModels)

	if len(shopModels) == 0 {
		return nil, errors.WithCode(code.ErrNoGoodsSelect, "未选择商品")
	}
	var goodsId []int32
	goodNumMap := make(map[int32]int32)
	for _, shopModel := range shopModels {
		goodsId = append(goodsId, shopModel.Goods)
		goodNumMap[shopModel.Goods] = shopModel.Nums
	}
	response := &do.GetShoppingBatchResponse{
		GoodsId:    goodsId,
		GoodNumMap: goodNumMap,
	}
	return response, nil
}

// ClearCheck 清空check状态
func (sc *shopCarts) ClearCheck(ctx context.Context, userID uint64) error {
	err := sc.db.Where("user = ?", userID).Delete(&do.ShoppingCartDO{}).Error
	if err != nil {
		log.Errorf(err.Error())
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

// 删除选中商品的购物车记录， 下订单了
// 从架构上来讲，这种实现有两种方案
// 下单后， 直接执行删除购物车的记录，比较简单
// 下单后什么都不做，直接给rocketmq发送一个消息，然后由rocketmq来执行删除购物车的记录
var _ v1.ShopCartStore = &shopCarts{}
