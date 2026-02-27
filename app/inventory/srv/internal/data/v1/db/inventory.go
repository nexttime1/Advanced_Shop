package mysql

import (
	"Advanced_Shop/app/inventory/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/code"
	code2 "Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"context"

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

var _ v1.InventoryStore = &inventorys{}
