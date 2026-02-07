package data

import (
	gpb "Advanced_Shop/api/goods/v1"
)

type DataFactory interface {
	Goods() gpb.GoodsClient
	Users() UserData
}
