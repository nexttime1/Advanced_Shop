package v1

import (
	v1 "Advanced_Shop/app/inventory/srv/internal/data/v1"
	code2 "Advanced_Shop/app/pkg/code"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"context"
	"database/sql"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/go-redsync/redsync/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"

	"Advanced_Shop/app/inventory/srv/internal/domain/do"
	"Advanced_Shop/app/inventory/srv/internal/domain/dto"
	"github.com/dtm-labs/client/dtmgrpc"

	"Advanced_Shop/pkg/log"
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
	// 使用屏障子事务
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		log.Errorf("订单%s创建屏障失败: %v", ordersn, err)
		return errors.WithCode(code.ErrUnknown, "创建屏障失败: %v", err)
	}
	// 新建实例
	rs := redsync.New(is.data.Pool())

	// 按照商品的id排序，然后从小大小逐个扣减库存，这样可以减少锁的竞争
	// 如果无序的话 那么就有可能订单a 扣减 1,3,4 订单B 扣减 3,2,1
	var detail = do.GoodsDetailList(details)
	sort.Sort(detail) // 实现了接口

	db := is.data.DB()

	// 直接用封装好的helper，闭包里的tx就是*gorm.DB，你的repo照常用
	return CallWithGorm(barrier, db, func(tx *gorm.DB) error {
		for _, goodsInfo := range detail {
			// 拿锁  一定是 先排序在拿锁
			mutex := rs.NewMutex(do.InventoryLockPrefix + strconv.FormatInt(int64(goodsInfo.GoodId), 10))
			defer mutex.Unlock()
			if err := mutex.Lock(); err != nil {
				return err
			}

			inv, err := is.data.Inventorys().GetWithTx(ctx, tx, uint64(goodsInfo.GoodId))
			if err != nil {
				return err
			}

			if inv.Stock < goodsInfo.Num {
				// TODO 错误码
				return status.Errorf(codes.Aborted, "库存不足: 商品%d", goodsInfo.GoodId)
			}

			if err := is.data.Inventorys().Reduce(ctx, tx, uint64(goodsInfo.GoodId), int(goodsInfo.Num)); err != nil {
				return err
			}

		}

		// 生成历史记录
		record := &do.StockSellDetailDO{
			OrderSn: ordersn,
			Status:  0,
			Detail:  detail,
			Version: 0,
		}
		err = is.data.Inventorys().CreateStockSellDetail(ctx, tx, record)
		if err != nil {
			log.Errorf("订单%s创建历史表失败", ordersn)
			return err
		}
		return nil
	})

}

func (is *inventoryService) Reback(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单%s归还库存", ordersn)

	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "创建屏障失败: %v", err)
	}

	// 新建 redsync 实例
	rs := redsync.New(is.data.Pool())

	gormDB := is.data.DB()
	var detail = do.GoodsDetailList(details)
	sort.Sort(detail) // 与 Sell 保持相同排序，防止死锁

	return CallWithGorm(barrier, gormDB, func(tx *gorm.DB) error {
		sellDetail, err := is.data.Inventorys().GetSellDetail(ctx, tx, ordersn)
		if err != nil {
			if errors.IsCode(err, code2.ErrInvSellDetailNotFound) {
				log.Infof("订单%s扣减记录不存在，空回滚", ordersn)
				return nil
			}
			return err
		}

		if sellDetail.Status == 2 || sellDetail.Status == 1 {
			log.Infof("订单%s已归还，跳过", ordersn)
			return nil
		}

		for _, goodsInfo := range detail {
			// 加 Redis 分布式锁
			mutex := rs.NewMutex(do.InventoryLockPrefix + strconv.FormatInt(int64(goodsInfo.GoodId), 10))
			defer mutex.Unlock()
			if err := mutex.Lock(); err != nil {
				log.Errorf("订单%s商品%d获取Redis锁失败: %v", ordersn, goodsInfo.GoodId, err)
				return errors.WithCode(code2.ErrRedisLock, "获取Redis锁失败: %v", err)
			}

			inv, err := is.data.Inventorys().GetWithTx(ctx, tx, uint64(goodsInfo.GoodId))
			if err != nil {
				log.Errorf("订单%s获取库存失败", ordersn)
				return err
			}
			inv.Stock += goodsInfo.Num

			if err := is.data.Inventorys().Increase(ctx, tx, inv); err != nil {
				log.Errorf("订单%s商品%d归还库存失败: %v", ordersn, goodsInfo.GoodId, err)
				return err
			}

			log.Infof("订单%s商品%d归还库存%d成功", ordersn, goodsInfo.GoodId, goodsInfo.Num)
		}

		if err = is.data.Inventorys().UpdateStockSellDetailStatus(ctx, tx, ordersn, 2); err != nil {
			log.Errorf("订单%s更新扣减库存记录失败", ordersn)
			return err
		}

		log.Infof("订单%s归还库存成功", ordersn)
		return nil
	})
}

func newInventoryService(s *service) *inventoryService {
	return &inventoryService{data: s.data, redisOptions: s.redisOptions}
}

var _ InventorySrv = &inventoryService{}

// CallWithGorm 让DTM屏障支持GORM事务
// 原理：拿到底层sql.DB开事务，把sql.Tx转给GORM，这样既能用屏障又能用GORM
func CallWithGorm(barrier *dtmcli.BranchBarrier, db *gorm.DB, busiCall func(tx *gorm.DB) error) error {
	// 1. 从gorm拿到底层 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 2. 调用DTM原生的CallWithDB，拿到 *sql.Tx
	return barrier.CallWithDB(sqlDB, func(sqlTx *sql.Tx) error {
		// 3. 把 *sql.Tx 包装成 *gorm.DB，这样你的repo方法全部可以继续用
		gormTx := db.WithContext(db.Statement.Context)

		// 关键：用DTM给的sqlTx替换gorm内部的连接
		// gorm提供了 ConnPool 接口，sql.Tx 实现了这个接口
		gormTx.Statement.ConnPool = sqlTx

		// 4. 执行业务逻辑，用的是包装好的gormTx
		return busiCall(gormTx)
	})
}
