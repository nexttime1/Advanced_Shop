package code

const (
	// ErrShopCartItemNotFound - 404: ShopCart item not found.
	ErrShopCartItemNotFound int = iota + 100701

	// ErrSubmitOrder - 400: Submit order error.
	ErrSubmitOrder

	// ErrNoGoodsSelect - 404: No Goods selected.
	ErrNoGoodsSelect

	// ErrOrderNotFound - 404: Order Not found
	ErrOrderNotFound

	// ErrOrderStatus - 404: Order Status fail
	ErrOrderStatus
)
