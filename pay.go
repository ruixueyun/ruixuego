package ruixuego

type OrderStatusRes struct {
	PlatformStatus int32  `json:"platform_status"`  // 订单状态 1 未支付 2  支付中 3 支付成功 4 支付关闭
	TradeStatus    string `json:"trade_status"`     // 三方支付状态
	TradeStateDesc string `json:"trade_state_desc"` // 三方支付描述
}
