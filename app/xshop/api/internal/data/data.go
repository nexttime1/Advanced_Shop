package data

import (
	apb "Advanced_Shop/api/action/v1"
	gpb "Advanced_Shop/api/goods/v1"
	ipb "Advanced_Shop/api/inventory/v1"
	opb "Advanced_Shop/api/order/v1"
)

type DataFactory interface {
	Inventory() ipb.InventoryClient
	Goods() gpb.GoodsClient
	Users() UserData
	Order() opb.OrderClient
	Address() apb.AddressClient
	Collection() apb.UserFavClient
	Message() apb.MessageClient
}
