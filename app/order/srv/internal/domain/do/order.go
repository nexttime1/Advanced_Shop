package do

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"Advanced_Shop/app/pkg/gorm"
)

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type OrderInfoDO struct {
	gorm.Model
	User         int32      `gorm:"type:int;index;comment:用户ID"`
	OrderSn      string     `gorm:"type:varchar(30);index;comment:订单编号（唯一）"`
	PayType      string     `gorm:"type:varchar(20);comment:支付方式（alipay/wechat）"`
	Status       string     `gorm:"type:varchar(20);comment:订单状态（PAYING/TRADE_SUCCESS/CLOSED）"`
	TradeNo      string     `gorm:"type:varchar(100);comment:第三方支付交易号"`
	OrderMount   float32    `gorm:"comment:订单总金额"`
	PayTime      *time.Time `gorm:"comment:支付时间"`
	Address      string     `gorm:"type:varchar(100);comment:收货地址"`
	SignerName   string     `gorm:"type:varchar(20);comment:签收人姓名"`
	SignerMobile string     `gorm:"type:varchar(11);comment:签收人手机号"`
	Post         string     `gorm:"type:varchar(20);comment:物流单号"`
}

// TableName 重写订单主表表名
func (OrderInfoDO) TableName() string {
	return "orderinfo"
}

type OrderGoodsModel struct {
	gorm.Model
	Order      int32   `gorm:"type:int;index;comment:订单ID"`
	Goods      int32   `gorm:"type:int;index;comment:商品ID"`
	GoodsName  string  `gorm:"type:varchar(100);comment:商品名称"`
	GoodsPrice float32 `gorm:"comment:商品单价"`
	GoodImages string  `gorm:"type:varchar(100);comment:商品图片"`
	Nums       int32   `gorm:"type:int;comment:商品数量"`
}

// TableName 重写订单商品明细表名
func (OrderGoodsModel) TableName() string {
	return "ordergoods"
}

type OrderInfoDOList struct {
	TotalCount int64          `json:"totalCount,omitempty"`
	Items      []*OrderInfoDO `json:"items"`
}

type OrderTransitionRequest struct {
	Id       int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`          // 订单ID（查询详情时必填，创建时不传）
	UserId   int32  `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`  // 用户ID
	Address  string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"` // 收货地址
	Name     string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`       // 签收人姓名
	Mobile   string `protobuf:"bytes,5,opt,name=mobile,proto3" json:"mobile,omitempty"`   // 签收人手机号
	Post     string `protobuf:"bytes,6,opt,name=post,proto3" json:"post,omitempty"`       // 物流单号（创建时可选，发货后填充）
	OrderSns string
}
