package mysql

import (
	"Advanced_Shop/app/inventory/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/code"
	code2 "Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"context"
	"time"

	"Advanced_Shop/app/inventory/srv/internal/data/v1"
	"Advanced_Shop/pkg/log"
	"gorm.io/gorm"
)

type inventorys struct {
	db *gorm.DB
}

func (i *inventorys) UpdateStockSellDetailStatus(ctx context.Context, txn *gorm.DB, ordersn string, status int32) error {
	db := i.db
	if txn != nil {
		db = txn
	}

	// update 语句如果没有更新的话那么不会报错，但是他会返回一个影响的行数，所以我们可以根据影响的行数来判断是否更新成功
	result := db.Model(do.StockSellDetailDO{}).Where("order_sn = ?", ordersn).Update("status", status)
	if result.Error != nil {
		return errors.WithCode(code2.ErrDatabase, result.Error.Error())
	}

	return nil
}

func (i *inventorys) GetSellDetail(ctx context.Context, txn *gorm.DB, ordersn string) (*do.StockSellDetailDO, error) {
	db := i.db
	if txn != nil {
		db = txn
	}
	var orderSellDetail do.StockSellDetailDO
	err := db.Where("order_sn = ?", ordersn).First(&orderSellDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrInvSellDetailNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &orderSellDetail, err
}

func (i *inventorys) Reduce(ctx context.Context, txn *gorm.DB, goodsID uint64, num int) error {
	db := i.db
	if txn != nil {
		db = txn
	}
	return db.Model(&do.InventoryDO{}).Where("goods=?", goodsID).Where("stock >= ?", num).UpdateColumn("stock", gorm.Expr("stock - ?", num)).Error
}

func (i *inventorys) Increase(ctx context.Context, txn *gorm.DB, inventory *do.InventoryDO) (int64, error) {
	db := i.db
	if txn != nil {
		db = txn
	}
	tx := db.Model(do.InventoryDO{}).
		Where("goods = ? and version = ?", inventory.Goods, inventory.Version).
		Select("stock", "version").
		Updates(map[string]interface{}{"stock": inventory.Stock, "version": inventory.Version + 1})
	// 查不到 也不报错
	return tx.RowsAffected, tx.Error

}

func (i *inventorys) CreateStockSellDetail(ctx context.Context, txn *gorm.DB, detail *do.StockSellDetailDO) error {
	db := i.db
	if txn != nil {
		db = txn
	}

	tx := db.Create(&detail)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (i *inventorys) Create(ctx context.Context, inv *do.InventoryDO) error {
	//设置库存， 如果我要更新库存
	tx := i.db.Create(&inv)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (i *inventorys) Get(ctx context.Context, goodsID uint64) (*do.InventoryDO, error) {
	inv := do.InventoryDO{}
	err := i.db.Where("goods = ?", goodsID).First(&inv).Error
	if err != nil {
		log.Errorf("get inv err: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrInventoryNotFound, err.Error())
		}

		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	return &inv, nil
}

func newInventorys(data *mysqlStore) *inventorys {
	return &inventorys{db: data.db}
}

func (i *inventorys) AutoReback(ctx context.Context, txn *gorm.DB, OrderSns string) do.MQMessageType {
	var history do.StockSellDetailDO
	err := txn.Where(do.StockSellDetailDO{OrderSn: OrderSns, Status: 1}).Take(&history).Error
	if err != nil {
		// 没找到 说明没有  或者 防止重复消费  这个消息直接跳过 所以
		return do.DirectPass
	}
	// 找到了  进行归还库存 并且 改历史记录 状态 变成2
	// 构造 Reback 函数需要的参数
	var info do.RebackInfo
	var list []*do.GoodsInvInfo
	for _, inv := range history.Detail {
		list = append(list, &do.GoodsInvInfo{
			GoodsId: inv.GoodId,
			Num:     inv.Num,
		})
	}
	info.GoodsInfo = list
	info.OrderSn = OrderSns
	number := Reback(txn, &info)
	if number == do.OptionFail {
		return do.OptionFail
	}
	// 下面一定是 继续了

	// 改历史记录 状态 变成2  并且要必须是 我们拿到的版本 如果不是 正常有个for循环 咱这个环境下 我们就没必要for了 因为一会就延迟再来一次
	err = txn.Model(&history).
		Where("id = ? and version = ?", history.ID, history.Version). // 基于 ID 和 version 做乐观锁
		Updates(map[string]interface{}{
			"status":  2,
			"version": history.Version + 1,
		}).Error
	if err != nil {
		return do.OptionFail
	}
	return do.Continuing
}

func Reback(tx *gorm.DB, info *do.RebackInfo) do.MQMessageType {

	// 传递事务，保证操作原子性
	for _, invInfo := range info.GoodsInfo {
		// 乐观锁保证 高并发情况下 不会发生错误  比如两个请求同一个商品进行归还  读取的值都是100 都加50  防止最后是150
		retryCount := 0
		maxOptimisticRetry := 10
		for retryCount < maxOptimisticRetry {
			var model do.InventoryDO
			err := tx.Where("goods = ?", invInfo.GoodsId).Take(&model).Error
			if err != nil {
				log.Errorf("商品库存不存在 err: %v", err)
				break
			}
			// 库存 +
			model.Stock += invInfo.Num

			err = tx.Model(do.InventoryDO{}).Where("goods = ? and version = ?", model.Goods, model.Version).Select("stock", "version").Updates(map[string]interface{}{"stock": model.Stock, "version": model.Version + 1}).Error
			if err != nil {
				retryCount++
				log.Warnf("商品%d乐观锁重试，当前次数: %d", invInfo.GoodsId, retryCount)
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}
		// 重试次数耗尽仍未成功
		if retryCount >= maxOptimisticRetry {
			log.Errorf("商品%d乐观锁重试次数耗尽，更新失败", invInfo.GoodsId)
			// 这里错了直接return  因为这个订单中的一个商品归还失败 肯定要回滚所以后面没必要归还
			return do.OptionFail
		}

	}
	// 归还成功 继续操作
	return do.Continuing

}

var _ v1.InventoryStore = &inventorys{}
