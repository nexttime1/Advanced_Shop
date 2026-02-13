package do

import (
	bgorm "Advanced_Shop/app/pkg/gorm"
	"Advanced_Shop/pkg/errors"
	"database/sql/driver"
	"encoding/json"
)

type GoodsDetail struct {
	Goods int32
	Num   int32
}

type GoodsDetailList []GoodsDetail

func (a GoodsDetailList) Len() int           { return len(a) }
func (a GoodsDetailList) Less(i, j int) bool { return a[i].Goods < a[j].Goods }
func (a GoodsDetailList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GoodsDetailList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, g)
}

type StockSellDetailDO struct {
	bgorm.Model `structs:"-"`
	OrderSn     string          `gorm:"type:varchar(200);index:unique"`
	Status      int32           //1 表示已扣减 2. 表示已归还
	Detail      GoodsDetailList `gorm:"type:json"`
	Version     int32           `gorm:"type:int"` //乐观锁
}

func (ssd *StockSellDetailDO) TableName() string {
	return "stock_sell_details"
}

type InventoryDO struct {
	bgorm.Model `structs:"-"`
	Goods       int32 `gorm:"type:int;index"`
	Stock       int32 `gorm:"type:int"`
	Version     int32 `gorm:"type:int"` //分布式锁
}

func (id *InventoryDO) TableName() string {
	return "inventory_models"
}
