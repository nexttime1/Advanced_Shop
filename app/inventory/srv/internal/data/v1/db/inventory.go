package mysql

import (
	"Advanced_Shop/app/inventory/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/code"
	code2 "Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"sort"
	"strconv"
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

func (i *inventorys) IncreaseSLock(ctx context.Context, txn *gorm.DB, inventory *do.InventoryDO) (int64, error) {
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

func (i *inventorys) Increase(ctx context.Context, txn *gorm.DB, inventory *do.InventoryDO) error {
	db := i.db
	if txn != nil {
		db = txn
	}
	err := db.Model(do.InventoryDO{}).Where("goods = ?", inventory.Goods).Update("stock", inventory.Stock).Error
	if err != nil {
		log.Errorf("increase inventory stock error: %v", err)
		return errors.WithCode(code2.ErrDatabase, "increase inventory stock error")
	}
	return err

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

func (i *inventorys) GetWithTx(ctx context.Context, txn *gorm.DB, goodsID uint64) (*do.InventoryDO, error) {
	inv := do.InventoryDO{}
	err := txn.Where("goods = ?", goodsID).First(&inv).Error
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

// AutoReback MQ 库存归还核心逻辑
func (i *inventorys) AutoReback(ctx context.Context, txn *gorm.DB, orderSn string, pool redsyncredis.Pool) (do.MQMessageType, error) {

	// 1. 查询归还记录（只查 status=2 已完成的，幂等判断）
	var history do.StockSellDetailDO

	// 先检查是否已经完成（幂等）
	err := txn.Where(&do.StockSellDetailDO{
		OrderSn: orderSn,
		Status:  do.StockSellStatusDone,
	}).Take(&history).Error

	if err == nil {
		// 已完成，幂等放行
		log.Infof("订单已完成归还，幂等跳过，OrderSn: %s", orderSn)
		return do.DirectPass, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return do.OptionFail, fmt.Errorf("查询归还记录失败: %w", err)
	}

	// 2. 查询待处理记录（status=0 待处理 或 status=1 处理中但事务回滚的）
	//    注意：status=1 可能是上次事务回滚后遗留，需要重新处理
	err = txn.Where("order_sn = ? AND status IN (?)",
		orderSn,
		[]int{do.StockSellStatusPending, do.StockSellStatusProcessing},
	).Take(&history).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有待处理记录，说明已经处理完成（被其他节点处理）
		log.Infof("未找到待处理归还记录，幂等跳过，OrderSn: %s", orderSn)
		return do.DirectPass, nil
	}

	if err != nil {
		return do.OptionFail, fmt.Errorf("查询待处理记录失败: %w", err)
	}

	//    抢占：status → 处理中(1)，使用乐观锁
	//    目的：多节点并发时，只有一个节点能抢到
	result := txn.Model(&history).
		Where("id = ? AND version = ? AND status IN (?)",
			history.ID,
			history.Version,
			[]int{do.StockSellStatusPending, do.StockSellStatusProcessing},
		).
		Updates(map[string]interface{}{
			"status":  do.StockSellStatusProcessing, // 标记为处理中
			"version": history.Version + 1,
		})

	if result.Error != nil {
		return do.OptionFail,
			fmt.Errorf("抢占归还记录失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		// 被其他节点抢占，本次跳过
		log.Infof("归还记录被其他节点抢占，跳过，OrderSn: %s", orderSn)
		return do.DirectPass, nil
	}

	// 更新本地version，后续操作使用新version
	history.Version = history.Version + 1
	history.Status = do.StockSellStatusProcessing

	// 4. 执行实际库存归还
	rebackInfo := &do.RebackInfo{
		GoodsInfo: history.Detail,
		OrderSn:   orderSn,
	}

	if err := Reback(ctx, txn, rebackInfo, pool); err != nil {
		// status=1 的记录也会随事务回滚到修改前的状态(0或1)
		// 下次重试时，查询 status IN (0,1) 仍然能找到，可以继续处理
		return do.OptionFail,
			fmt.Errorf("执行库存归还失败: %w", err)
	}

	// 5. 标记为已完成（status: 处理中→完成）
	err = txn.Model(&history).
		Where("id = ? AND version = ?", history.ID, history.Version).
		Updates(map[string]interface{}{
			"status":  do.StockSellStatusDone,
			"version": history.Version + 1,
		}).Error

	if err != nil {
		return do.OptionFail,
			fmt.Errorf("更新归还状态为完成失败: %w", err)
	}

	return do.Continuing, nil
}

// ==================== Redis 分布式锁执行归还 ====================
func Reback(ctx context.Context, tx *gorm.DB, info *do.RebackInfo, pool redsyncredis.Pool) error {
	rs := redsync.New(pool)
	//  GoodsDetailList 就是 []GoodsDetail，直接强转，复用已有排序接口
	sort.Sort(do.GoodsDetailList(info.GoodsInfo))

	// 按排好序的顺序逐个拿锁，拿锁失败则释放已持有的，避免泄漏
	mutexes := make([]*redsync.Mutex, 0, len(info.GoodsInfo))
	for _, invInfo := range info.GoodsInfo {
		mutex := rs.NewMutex(
			do.InventoryLockPrefix+strconv.FormatInt(int64(invInfo.GoodId), 10),
			redsync.WithExpiry(8*time.Second),
			redsync.WithTries(3),
			redsync.WithRetryDelay(100*time.Millisecond),
		)
		if err := mutex.LockContext(ctx); err != nil {
			for _, acquired := range mutexes {
				if _, unlockErr := acquired.Unlock(); unlockErr != nil {
					log.Errorf("归还释放锁失败，OrderSn: %s, err: %v", info.OrderSn, unlockErr)
				}
			}
			return fmt.Errorf("商品%d获取分布式锁失败: %w", invInfo.GoodId, err)
		}
		mutexes = append(mutexes, mutex)
	}

	// 函数退出统一释放所有锁
	defer func() {
		for _, mutex := range mutexes {
			if _, err := mutex.Unlock(); err != nil {
				log.Errorf("归还释放分布式锁失败，OrderSn: %s, err: %v", info.OrderSn, err)
			}
		}
	}()

	// 持锁后直接执行，无需乐观锁重试
	for _, invInfo := range info.GoodsInfo {
		if err := rebackSingleGoods(ctx, tx, invInfo); err != nil {
			return fmt.Errorf("商品%d库存归还失败: %w", invInfo.GoodId, err)
		}
	}
	return nil
}

// ==================== 单个商品归还（持锁后直接读写）====================
func rebackSingleGoods(ctx context.Context, tx *gorm.DB, invInfo do.GoodsDetail) error {

	var model do.InventoryDO
	if err := tx.Where("goods = ?", invInfo.GoodId).Take(&model).Error; err != nil {
		return fmt.Errorf("查询商品库存失败，goodsId: %d, err: %w", invInfo.GoodId, err)
	}

	newStock := model.Stock + invInfo.Num

	// ✅ 持有分布式锁，并发安全，不需要 version 作为更新条件
	result := tx.Model(&do.InventoryDO{}).
		Where("goods = ?", model.Goods).
		Updates(map[string]interface{}{
			"stock":   newStock,
			"version": model.Version + 1, // 依然自增，保留审计语义
		})

	if result.Error != nil {
		return fmt.Errorf("更新库存失败，goodsId: %d, err: %w", invInfo.GoodId, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("商品%d库存记录不存在", invInfo.GoodId)
	}

	log.Infof("商品%d库存归还成功，归还数量: %d，新库存: %d",
		invInfo.GoodId, invInfo.Num, newStock)
	return nil
}

var _ v1.InventoryStore = &inventorys{}
