package v1

import (
	proto "Advanced_Shop/api/goods/v1"
	proto2 "Advanced_Shop/api/inventory/v1"
	"gorm.io/gorm"
)

type DataFactory interface {
	Orders() OrderStore
	ShopCarts() ShopCartStore
	Goods() proto.GoodsClient
	Inventorys() proto2.InventoryClient

	Begin() *gorm.DB
}
