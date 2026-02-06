package data

import (
	gpb "mxshop/api/goods/v1"
)

type DataFactory interface {
	Goods() gpb.GoodsClient
	Users() UserData
}
