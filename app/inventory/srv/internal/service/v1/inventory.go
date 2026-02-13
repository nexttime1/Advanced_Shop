package v1

import (
	v1 "Advanced_Shop/app/inventory/srv/internal/data/v1"
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/pkg/errors"
	"context"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"sort"
	"time"

	"Advanced_Shop/app/inventory/srv/internal/domain/do"
	"Advanced_Shop/app/inventory/srv/internal/domain/dto"

	"Advanced_Shop/pkg/log"
)

const (
	inventoryLockPrefix = "inventory_"
	orderLockPrefix     = "order_"
)

const (
	maxOptimisticRetry      = 10 // 乐观锁重试
	optimisticRetryInterval = 100 * time.Millisecond
)

type InventorySrv interface {
	// Create 设置库存
	Create(ctx context.Context, inv *dto.InventoryDTO) error

	// Get 根据商品的id查询库存
	Get(ctx context.Context, goodsID uint64) (*dto.InventoryDTO, error)

	// Sell 扣减库存
	Sell(ctx context.Context, ordersn string, detail []do.GoodsDetail) error

	// Reback 归还库存
	Reback(ctx context.Context, ordersn string, detail []do.GoodsDetail) error
}

type inventoryService struct {
	data v1.DataFactory

	redisOptions *options.RedisOptions

	pool redsyncredis.Pool
}

func (is *inventoryService) Create(ctx context.Context, inv *dto.InventoryDTO) error {
	return is.data.Inventorys().Create(ctx, &inv.InventoryDO)
}

func (is *inventoryService) Get(ctx context.Context, goodsID uint64) (*dto.InventoryDTO, error) {
	inv, err := is.data.Inventorys().Get(ctx, goodsID)
	if err != nil {
		return nil, err
	}
	return &dto.InventoryDTO{InventoryDO: *inv}, nil
}

func (is *inventoryService) Sell(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单%s扣减库存", ordersn)

	// 新建实例
	rs := redsync.New(is.pool)

	// 按照商品的id排序，然后从小大小逐个扣减库存，这样可以减少锁的竞争
	// 如果无序的话 那么就有可能订单a 扣减 1,3,4 订单B 扣减 3,2,1
	var detail = do.GoodsDetailList(details)
	sort.Sort(detail) // 实现了接口

	txn := is.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			_ = txn.Rollback()
			log.Error("事务进行中出现异常，回滚")
			return
		}
	}()

	// 生成历史记录
	record := do.StockSellDetailDO{
		OrderSn: ordersn,
		Status:  1, // 1 代表 扣了没还
		Detail:  detail,
	}

	for _, goodsInfo := range detail {

		// 拿锁  一定是 先排序在拿锁
		mutex := rs.NewMutex(inventoryLockPrefix + ordersn)
		if err := mutex.Lock(); err != nil {
			log.Errorf("订单%s获取锁失败", ordersn)
			txn.Rollback()
			return err
		}

		inv, err := is.data.Inventorys().Get(ctx, uint64(goodsInfo.Goods))
		if err != nil {
			log.Errorf("订单%s获取库存失败", ordersn)
			txn.Rollback()
			return err
		}

		// 判断库存是否充足
		if inv.Stock < goodsInfo.Num {
			txn.Rollback() //回滚
			log.Errorf("商品%d库存%d不足, 现有库存: %d", goodsInfo.Goods, goodsInfo.Num, inv.Stock)
			return errors.WithCode(code.ErrInvNotEnough, "库存不足")
		}
		inv.Stock -= goodsInfo.Num

		err = is.data.Inventorys().Reduce(ctx, txn, uint64(goodsInfo.Goods), int(goodsInfo.Num))
		if err != nil {
			txn.Rollback() //回滚
			log.Errorf("订单%s扣减库存失败", ordersn)
			return err
		}

		//释放锁
		if _, err = mutex.Unlock(); err != nil {
			txn.Rollback() //回滚
			log.Errorf("订单%s释放锁出现异常", ordersn)
		}
	}

	err := is.data.Inventorys().CreateStockSellDetail(ctx, txn, &record)
	if err != nil {
		txn.Rollback() //回滚
		log.Errorf("订单%s创建扣减库存记录失败", ordersn)
		return err
	}

	txn.Commit()
	return nil
}

func (is *inventoryService) Reback(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单%s归还库存", ordersn)
	// 新建实例
	rs := redsync.New(is.pool)

	// 防止 抖动或者多次请求 同时间到达 直接开锁 最后放开
	mutex := rs.NewMutex(orderLockPrefix + ordersn)
	mutex.Lock()
	defer mutex.Unlock()
	txn := is.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			_ = txn.Rollback()
			log.Error("事务进行中出现异常，回滚")
			return
		}
	}()

	sellDetail, err := is.data.Inventorys().GetSellDetail(ctx, txn, ordersn)
	if err != nil {
		txn.Rollback()
		if errors.IsCode(err, code.ErrInvSellDetailNotFound) {
			//空回滚
			log.Errorf("订单%s扣减库存记录不存在, 忽略", ordersn)
			return nil
		}
		log.Errorf("订单%s获取扣减库存记录失败", ordersn)
		return err
	}

	if sellDetail.Status == 2 {
		log.Infof("订单%s扣减库存记录已经归还, 忽略", ordersn)
		return nil
	}

	var detail = do.GoodsDetailList(details)
	sort.Sort(detail)

	for _, goodsInfo := range detail {
		retryCount := 0
		var updateSuccess bool

		// 乐观锁重试逻辑
		for retryCount < maxOptimisticRetry {
			inv, err := is.data.Inventorys().Get(ctx, uint64(goodsInfo.Goods))
			if err != nil {
				txn.Rollback() //回滚
				log.Errorf("订单%s获取库存失败", ordersn)
				return err
			}
			inv.Stock += goodsInfo.Num

			row, err := is.data.Inventorys().Increase(ctx, txn, inv)
			if err != nil {
				txn.Rollback() // 回滚
				log.Errorf("订单%s归还库存失败", ordersn)
				return err
			}
			if row > 0 {
				updateSuccess = true
				log.Infof("订单%s商品%d乐观锁更新库存成功，库存变为%d", ordersn, inv.Stock)
				break
			}
			// 重试
			retryCount++
			log.Warnf("订单%s商品%d乐观锁冲突，重试次数: %d/%d",
				ordersn, inv.Goods, retryCount, maxOptimisticRetry)
			time.Sleep(optimisticRetryInterval)
		}
		if !updateSuccess {
			_ = txn.Rollback()
			errMsg := "订单" + ordersn + "商品" + string(goodsInfo.Goods) + "乐观锁重试次数耗尽，库存归还失败"
			log.Errorf(errMsg)
			return errors.WithCode(code.ErrOptimisticRetry, errMsg)
		}
	}

	err = is.data.Inventorys().UpdateStockSellDetailStatus(ctx, txn, ordersn, 2)
	if err != nil {
		txn.Rollback() //回滚
		log.Errorf("订单%s更新扣减库存记录失败", ordersn)
		return err
	}

	err = txn.Commit().Error
	if err != nil {
		log.Errorf("订单%s提交归还库存事务失败: %v", ordersn, err)
		return err
	}

	log.Infof("订单%s归还库存成功", ordersn)
	return nil
}

func newInventoryService(s *service) *inventoryService {
	return &inventoryService{data: s.data, redisOptions: s.redisOptions, pool: s.pool}
}

var _ InventorySrv = &inventoryService{}
