//go:generate codegen -type=int

package code

// Order: order service errors.
// Code must start with 1007xx.
const (
	// ErrShopCartItemNotFound - 404: ShopCart item not found.
	ErrShopCartItemNotFound int = iota + 100701

	// ErrSubmitOrder - 500: Failed to submit order.
	ErrSubmitOrder

	// ErrNoGoodsSelect - 400: No goods selected.
	ErrNoGoodsSelect

	// ErrOrderNotFound - 404: Order not found.
	ErrOrderNotFound

	// ErrOrderStatus - 500: Failed to update order status.
	ErrOrderStatus

	// ErrInvalidParameter - 400: Invalid request parameter.
	ErrInvalidParameter
)
