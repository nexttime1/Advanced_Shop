package do

import (
	bgorm "Advanced_Shop/app/pkg/gorm"
	"Advanced_Shop/pkg/errors"
	"database/sql/driver"
	"encoding/json"
)

type GoodsDetailList []GoodsDetail
type GoodsDetail struct {
	GoodId int32
	Num    int32
}

func (a GoodsDetailList) Len() int           { return len(a) }
func (a GoodsDetailList) Less(i, j int) bool { return a[i].GoodId < a[j].GoodId }
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

type OrderMQMessageRequest struct {
	Id       int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`          // 订单ID（查询详情时必填，创建时不传）
	UserId   int32  `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`  // 用户ID
	Address  string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"` // 收货地址
	Name     string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`       // 签收人姓名
	Mobile   string `protobuf:"bytes,5,opt,name=mobile,proto3" json:"mobile,omitempty"`   // 签收人手机号
	Post     string `protobuf:"bytes,6,opt,name=post,proto3" json:"post,omitempty"`       // 物流单号（创建时可选，发货后填充）
	OrderSns string
}

type GoodsInvInfo struct {
	GoodsId int32 `json:"goods_id"`
	Num     int32 `json:"num"`
}

type RebackInfo struct {
	GoodsInfo []*GoodsInvInfo `json:"goods_info"`
	OrderSn   string          `json:"order_sn"`
}
