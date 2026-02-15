package do

import "Advanced_Shop/app/pkg/gorm"

type ShoppingCartDOList struct {
	TotalCount int64             `json:"totalCount,omitempty"`
	Items      []*ShoppingCartDO `json:"items"`
}

type GetShoppingBatchResponse struct {
	GoodsId    []int32
	GoodNumMap map[int32]int32
}

// ShoppingCartModel
type ShoppingCartDO struct {
	gorm.Model
	User    int32 `gorm:"type:int;index;comment:用户ID"`
	Goods   int32 `gorm:"type:int;index;comment:商品ID"`
	Nums    int32 `gorm:"type:int;comment:商品数量"`
	Checked *bool `gorm:"comment:是否勾选（结算）"`
}

// TableName 重写购物车表名
func (ShoppingCartDO) TableName() string {
	return "shoppingcart"
}

type CartUpdateMap struct {
	Nums    int32 `json:"nums" structs:"nums"`
	Checked *bool `json:"checked" structs:"checked"`
}
