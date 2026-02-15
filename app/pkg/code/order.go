package code

const (
	// ErrShopCartItemNotFound - 404: ShopCart item not found.
	ErrShopCartItemNotFound int = iota + 100701 // 100701

	// ErrSubmitOrder - 400: Submit order error.
	ErrSubmitOrder // 100702

	// ErrNoGoodsSelect - 404: No Goods selected.
	ErrNoGoodsSelect // 100703

	// ErrOrderNotFound - 404: Order Not found
	ErrOrderNotFound // 100704

	// ErrOrderStatus - 400: Order Status update failed
	ErrOrderStatus // 100705

	// ErrInvalidParameter - 400: Invalid request parameter.
	ErrInvalidParameter // 100706
)
