package v1

import proto "Advanced_Shop/api/goods/v1"

type DataFactory interface {
	Address() AddressStore
	Collection() CollectionStore
	Goods() proto.GoodsClient
	Messages() MessageStore
}
