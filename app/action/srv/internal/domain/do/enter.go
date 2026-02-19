package do

// 订单支付方式枚举（使用iota规范管理）
const (
	PayTypeAlipay = "alipay" // 支付宝支付
	PayTypeWechat = "wechat" // 微信支付
)

// 订单状态枚举（使用iota规范管理）
const (
	OrderStatusPaying        = "PAYING"         // 待支付
	OrderStatusTradeSuccess  = "TRADE_SUCCESS"  // 支付成功
	OrderStatusTradeClosed   = "TRADE_CLOSED"   // 交易关闭
	OrderStatusShipped       = "SHIPPED"        // 已发货
	OrderStatusReceived      = "RECEIVED"       // 已收货
	OrderStatusRefundSuccess = "REFUND_SUCCESS" // 退款成功
)
