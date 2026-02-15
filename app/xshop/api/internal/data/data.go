package data

import (
	gpb "Advanced_Shop/api/goods/v1"
	opb "Advanced_Shop/api/order/v1"
)

type DataFactory interface {
	Goods() gpb.GoodsClient
	Users() UserData
	Order() opb.OrderClient
}
