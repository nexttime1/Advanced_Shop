//go:generate codegen -type=int

package code

// Inventory: inventory service errors.
// Code must start with 1006xx.
const (
	// ErrInventoryNotFound - 404: Inventory not found.
	ErrInventoryNotFound int = iota + 100601

	// ErrInvSellDetailNotFound - 404: Inventory sell detail not found.
	ErrInvSellDetailNotFound

	// ErrInvNotEnough - 400: Inventory not enough.
	ErrInvNotEnough

	// ErrOptimisticRetry - 500: Optimistic lock retry limit exceeded.
	ErrOptimisticRetry
)
