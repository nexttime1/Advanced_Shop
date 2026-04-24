package do

import "time"

type MQMessageType int

const (
	DirectPass MQMessageType = 0
	OptionFail MQMessageType = 1
	Continuing MQMessageType = 2
)

const (
	StockSellStatusPending    = 0 // 待处理
	StockSellStatusProcessing = 1 // 处理中（已抢占）
	StockSellStatusDone       = 2 // 已完成
)

const (
	MaxOptimisticRetry      = 10
	OptimisticRetryInterval = 100 * time.Millisecond
)

const (
	InventoryLockPrefix = "inventory_"
	OrderLockPrefix     = "order_"
)
